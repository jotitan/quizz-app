package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/model"
	"github.com/quizz-app/sse"
	"math/rand"
	"strings"
)

// Manage an execution of a quizz

type GameService struct {
	quizzService QuizzService
	runningGames map[string]*model.Game
}

func NewGameService(quizzService QuizzService) GameService {
	return GameService{quizzService, make(map[string]*model.Game)}
}

func (ps GameService) GetByIdAndCheck(gameId, secureID string) (*model.Game, error) {
	game, err := ps.GetById(gameId)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(game.SecureId, secureID) {
		return nil, errors.New("no access to this game")
	}
	return game, nil
}

func (ps GameService) GetById(gameId string) (*model.Game, error) {
	if game, exist := ps.runningGames[gameId]; exist {
		return game, nil
	}
	return nil, errors.New("game not found")
}

func (ps GameService) Join(game *model.Game, name string) (string, error) {
	if game.IsStarted() {
		return "", errors.New("game already started")
	}
	if strings.EqualFold("", name) {
		return "", errors.New("must specify a real name")
	}
	// Check if name already exist
	if game.IsPlayerExists(name) {
		return "", errors.New("player with name already exists")
	}
	player := game.AddPlayer(name, ps.generateId(6))
	return player.Id, nil
}

func (ps GameService) Create(quizzId string) (*model.Game, error) {
	quizz, err := ps.quizzService.Get(quizzId)
	if err != nil {
		return nil, err
	}
	game := model.NewGame(quizz, ps.generateUniqueId(), ps.generateId(12))
	ps.runningGames[game.Id] = game
	return game, nil
}

func (ps GameService) Connect(game *model.Game, playerId, playerName string, c *gin.Context) error {
	if !game.CheckPlayer(playerName, playerId) {
		return errors.New("impossible to connect user")
	}
	sseCommunication := sse.NewSSECommunicate(c.Writer, playerId)
	game.ConnectPlayer(playerName, sseCommunication)
	if !game.IsStarted() {
		game.SendMasterMessage(sse.Message{Event: "join", Payload: fmt.Sprintf("{\"player\":\"%s\"}", playerName)})
	}
	// Let connexion open, blocking
	sseCommunication.Chanel <- sse.Message{Event: "welcome", Payload: ps.definePlayerStatusGame(game, playerName)}
	sseCommunication.Run()

	return errors.New("end of connection")
}

func (ps GameService) generateUniqueId() string {
	// generate 4 random letter and check if already exist, round again if yes
	code := ps.generateId(4)
	if _, exist := ps.runningGames[code]; exist {
		return ps.generateUniqueId()
	}
	return code
}

func (ps GameService) generateId(size int) string {
	// generate 4 random letter and check if already exist, round again if yes
	code := ""
	for i := 0; i < size; i++ {
		code += randomLetter()
	}
	return code
}

func (ps GameService) Start(game *model.Game, c *gin.Context) {
	game.Start()
}

func (ps GameService) ComputeScore(game *model.Game) model.ScoreQuestion {
	return game.ComputeScore()
}

func (ps GameService) ConnectMaster(game *model.Game, c *gin.Context) {
	// Create SSE for master communication
	sseCommunication := sse.NewSSECommunicate(c.Writer, game.Id)
	game.SSEMaster(sseCommunication)
	// Message sended depends on game status

	game.SendMasterMessage(sse.Message{Event: "welcome", Payload: ps.defineStatusGame(game)})
	sseCommunication.Run()
}

func (ps GameService) definePlayerStatusGame(game *model.Game, login string) string {
	if !game.IsStarted() {
		return "{\"status\":\"waiting\"}"
	}
	if game.IsQuestionRunning() {
		// Check if player already answered
		if game.HasPlayerAnswered(login) {
			return "{\"status\":\"waiting\"}"
		}
		q, _ := game.PreviousQuestion()
		return fmt.Sprintf("{\"status\":\"question\",%s,\"question\":{\"title\":\"%s\",\"nb\":%d,\"position\":%d}}",
			getScoreAsString(game, login), q.Title, len(q.Answers), q.Position)
	}
	return fmt.Sprintf("{\"status\":\"score\",%s,\"end\":%t}", getScoreAsString(game, login), len(game.Quizz.Questions) == game.GetCurrentQuestion())
}

func getScoreAsString(game *model.Game, login string) string {
	score, rank := game.GetPlayerScore(login)
	return fmt.Sprintf("\"score\":%d,\"rank\":%d", score, rank)
}

func (ps GameService) defineStatusGame(game *model.Game) string {
	// different game status when connect :
	// 1 - not started (waiting users, sended connected)
	// 2 - running question (send users who have answered and current question)
	// 3 - others, show total
	if !game.IsStarted() {
		data, _ := json.Marshal(game.GetUsersNames())
		return fmt.Sprintf("{\"status\":\"waiting\",\"users\":%s}", string(data))
	}
	if game.IsQuestionRunning() {
		pq, err := game.PreviousQuestion()
		if err == nil {
			data, _ := json.Marshal(pq)
			users, _ := json.Marshal(game.GetAnswered())
			totalPlayers := len(game.GetUsersNames())
			return fmt.Sprintf("{\"status\":\"answer\",\"total_players\":%d,\"users\":%s,\"question\":%s}", totalPlayers, users, string(data))
		}
	}
	if game.ShouldCompute() {
		return "{\"status\":\"compute_score\"}"
	}

	return fmt.Sprintf("{\"status\":\"score\",\"current\":%d}", game.GetCurrentQuestion())
}

func (ps GameService) ForceEndAnswer(game *model.Game) {
	game.ForceEndAnswer()
}

func (ps GameService) PlayNextQuestion(game *model.Game) (model.Question, bool) {
	// Check if still have question
	if !game.HasMoreQuestion() {
		return model.Question{}, false
	}
	q, _ := game.NextQuestion()

	// Send message to players
	game.SendMessages(sse.Message{Event: "question", Payload: fmt.Sprintf("{\"title\":\"%s\",\"nb\":%d,\"position\":%d}",
		q.Title, len(q.Answers), q.Position)})
	return q, true
}

func (ps GameService) AnswerQuestion(game *model.Game, idPlayer string, answer int) error {
	player, err := game.GetPlayerById(idPlayer)
	if err != nil {
		return err
	}
	allAnswers, err := game.SaveAnswerPlayer(idPlayer, answer)
	if err != nil {
		return err
	}
	game.SendMasterMessage(sse.Message{Event: "answer", Payload: fmt.Sprintf("{\"player\":\"%s\"}", player.Login)})
	if allAnswers {
		game.EndQuestion()
	}
	return nil
}

func (ps GameService) GetScore(game *model.Game) map[string]model.RecapScore {
	return game.GetScore()
}

func randomLetter() string {
	return fmt.Sprintf("%c", 65+rand.Intn(26))
}
