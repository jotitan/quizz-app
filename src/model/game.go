package model

import (
	"errors"
	"fmt"
	"github.com/quizz-app/sse"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Game struct {
	Quizz Quizz
	// Join id on 4 letters
	Id string
	// SecureId to connect master
	SecureId string
	// Players
	players            []*Player
	counterPlayer      int32
	byNames            map[string]*Player
	byIds              map[string]*Player
	currentAnswers     playerAnswers
	currentQuestion    int
	isStart            bool
	canAnswer          bool
	masterCommunicate  *sse.SSECommunicate
	computeScoreByTime bool
}

type playerAnswers struct {
	answers     map[string]playerAnswer
	computeToDo bool
	locker      sync.Mutex
	time        time.Time
}

func (pa *playerAnswers) reset() {
	pa.locker.Lock()
	defer pa.locker.Unlock()
	pa.answers = make(map[string]playerAnswer)
	pa.computeToDo = false
	pa.time = time.Now()
}

// return total of answers
func (pa *playerAnswers) addAnswer(idPlayer string, answer int) (int, error) {
	pa.locker.Lock()
	defer pa.locker.Unlock()
	if _, exist := pa.answers[idPlayer]; exist {
		return 0, errors.New("player already answered")
	}
	pa.answers[idPlayer] = playerAnswer{idPlayer: idPlayer, answer: answer, date: time.Now()}
	return len(pa.answers), nil
}

func (pa *playerAnswers) HasAnswered(login string) bool {
	_, exist := pa.answers[login]
	return exist
}

func (pa *playerAnswers) GetUsers(players map[string]*Player) []string {
	users := make([]string, 0, len(pa.answers))
	for id := range pa.answers {
		users = append(users, players[id].Login)
	}
	return users
}

type playerAnswer struct {
	idPlayer string
	answer   int
	date     time.Time
}

type Player struct {
	Login      string
	Id         string
	IdPosition int32
	score      int
	fineScore  int
	rank       int
	connected  bool
	messages   chan sse.Message
}

func (g Game) IsPlayerExists(name string) bool {
	_, exist := g.byNames[name]
	return exist
}

func (g *Game) IsStarted() bool {
	return g.isStart
}

func (g *Game) Start() {
	g.isStart = true
}

func (g *Game) SSEMaster(communicate *sse.SSECommunicate) {
	g.masterCommunicate = communicate
}

func (g *Game) ShouldCompute() bool {
	return g.currentAnswers.computeToDo
}

func (g *Game) IsQuestionRunning() bool {
	return g.canAnswer
}

func (g *Game) GetAnswered() []string {
	return g.currentAnswers.GetUsers(g.byIds)
}

func (g *Game) HasPlayerAnswered(login string) bool {
	return g.currentAnswers.HasAnswered(login)
}

func (g Game) CheckPlayer(name, id string) (bool, int32) {
	player, exist := g.byNames[name]
	if !exist {
		return false, 0
	}
	return strings.EqualFold(player.Id, id), player.IdPosition
}

func (g Game) CheckPlayerById(id string) (bool, int32) {
	player, exist := g.byIds[id]
	if !exist {
		return false, 0
	}
	return true, player.IdPosition
}

func (g Game) GetUsersNames() []string {
	names := make([]string, len(g.players))
	for i, p := range g.players {
		names[i] = p.Login
	}
	return names
}

func (g *Game) ConnectPlayer(name string, communicate *sse.SSECommunicate) {
	g.byNames[name].connected = true
	g.byNames[name].messages = communicate.Chanel
}

func (g *Game) SendScoreToPlayer(isFinish bool) {
	for _, player := range g.players {
		if player.connected {
			score, rank := g.GetPlayerScore(player.Login)
			payload := fmt.Sprintf("{\"rank\":%d,\"score\":%d,\"end\":%t}", rank, score, isFinish)
			player.messages <- sse.Message{Event: "score", Payload: payload}
		}
	}
}

func (g *Game) SendMessages(message sse.Message) {
	for _, player := range g.players {
		if player.connected {
			player.messages <- message
		}
	}
}

func (g *Game) SendMasterMessage(message sse.Message) {
	g.masterCommunicate.Chanel <- message
}

func (g *Game) AddPlayer(name, id string) *Player {
	player := &Player{Login: name, Id: fmt.Sprintf("%s%d", id, len(g.players)+1), score: 0, fineScore: 0, connected: false}
	player.IdPosition = atomic.AddInt32(&g.counterPlayer, 1)
	g.players = append(g.players, player)
	g.byNames[name] = player
	g.byIds[player.Id] = player
	return player
}

func (g *Game) HasMoreQuestion() bool {
	return g.currentQuestion < len(g.Quizz.Questions)
}

func (g *Game) PreviousQuestion() (Question, error) {
	if g.currentQuestion == 0 {
		return Question{}, errors.New("no previous question")
	}
	return g.Quizz.Questions[g.currentQuestion-1], nil
}

func (g *Game) GetCurrentQuestion() int {
	return g.currentQuestion
}

func (g *Game) NextQuestion() (Question, error) {
	if !g.HasMoreQuestion() {
		return Question{}, errors.New("no more question")
	}
	g.canAnswer = true
	g.currentAnswers.reset()
	q := g.Quizz.Questions[g.currentQuestion]
	q.Position = g.currentQuestion
	g.currentQuestion++
	return q, nil
}

func (g *Game) EndQuestion() {
	g.canAnswer = false
	g.currentAnswers.computeToDo = true
	message := sse.Message{Event: "end-answers", Payload: "{}"}
	g.SendMasterMessage(message)
	g.SendMessages(message)
}

// return login of player
func (g *Game) increaseScore(playerId string, restTime int) string {
	player, _ := g.GetPlayerById(playerId)
	player.score++
	player.fineScore += restTime
	return player.Login
}

func (g *Game) GetPlayerById(idPlayer string) (*Player, error) {
	player, exist := g.byIds[idPlayer]
	if !exist {
		return nil, errors.New("unknown player")
	}
	return player, nil
}

func (g *Game) SaveAnswerPlayer(idPlayer string, answer int) (bool, error) {
	if !g.canAnswer {
		return false, errors.New("question end")
	}
	nbAnswers, err := g.currentAnswers.addAnswer(idPlayer, answer)
	if err != nil {
		return false, err
	}
	return nbAnswers >= len(g.players), nil
}

type Repartition struct {
	Answer string `json:"answer"`
	Nb     int    `json:"nb"`
}

func (g *Game) GetRepartitionWithAnswers() []Repartition {
	answers := g.Quizz.Questions[g.currentQuestion-1].Answers
	repartitions := make([]Repartition, len(answers))
	for pos, answer := range answers {
		repartitions[pos] = Repartition{Answer: answer.Text, Nb: 0}
	}
	for _, a := range g.currentAnswers.answers {
		repartitions[a.answer].Nb++
	}
	return repartitions
}

func (g *Game) ForceEndAnswer() {
	g.EndQuestion()
}

// ComputeScore return list of winners
func (g *Game) ComputeScore() ScoreQuestion {
	question := g.Quizz.Questions[g.currentQuestion-1]
	good := question.GetGoodAnswers()
	results := ScoreQuestion{Winners: make([]string, 0), GoodAnswers: make([]string, 0, len(good))}
	for _, g := range good {
		results.GoodAnswers = append(results.GoodAnswers, g)
	}
	for playerId, answer := range g.currentAnswers.answers {
		if _, haveGood := good[answer.answer]; haveGood {
			restTime := question.MaxTime - int(answer.date.Sub(g.currentAnswers.time).Seconds())
			results.Winners = append(results.Winners, g.increaseScore(playerId, restTime))
		}
	}
	results.Repartition = g.GetRepartitionWithAnswers()
	g.currentAnswers.reset()
	g.computePlayersRank()
	g.SendScoreToPlayer(len(g.Quizz.Questions) == g.currentQuestion)
	// If last question, notify user
	return results
}

func (g *Game) getRankCompute() func(a, b int) bool {
	if g.computeScoreByTime {
		return func(a, b int) bool { return g.players[a].fineScore > g.players[b].fineScore }
	}
	return func(a, b int) bool { return g.players[a].score > g.players[b].score }
}

func (g *Game) computePlayersRank() {
	sort.Slice(g.players, g.getRankCompute())
	for i, player := range g.players {
		player.rank = i
	}
}

func (g *Game) GetPlayerScore(login string) (int, int) {
	if player, exist := g.byNames[login]; exist {
		if g.computeScoreByTime {
			return player.fineScore, player.rank
		}
		return player.score, player.rank
	}
	return 0, 0
}

func (g *Game) GetScore() map[string]RecapScore {
	scores := make(map[string]RecapScore)
	for _, player := range g.players {
		scores[player.Login] = g.getRecapScore(*player)
	}
	return scores
}

func (g *Game) getRecapScore(p Player) RecapScore {
	if g.computeScoreByTime {
		return RecapScore{p.fineScore, p.score}
	}
	return RecapScore{p.score, p.fineScore}
}

func NewGame(quizz Quizz, uniqueId, secureId string, scoreWithTime bool) *Game {
	return &Game{
		Quizz:              quizz,
		players:            make([]*Player, 0),
		byNames:            make(map[string]*Player),
		byIds:              make(map[string]*Player),
		currentQuestion:    0,
		isStart:            false,
		Id:                 uniqueId,
		SecureId:           secureId,
		canAnswer:          false,
		computeScoreByTime: scoreWithTime,
	}
}

type ScoreQuestion struct {
	Winners          []string       `json:"winners"`
	GoodAnswers      []string       `json:"good"`
	Repartition      []Repartition  `json:"repartition"`
	RepartitionAsMap map[string]int `json:"repartition_map"`
}

type RecapScore struct {
	Score     int
	TimeScore int
}
