package service

import (
	"errors"
	"github.com/quizz-app/config"
	"github.com/quizz-app/model"
	"github.com/quizz-app/persistence"
	"github.com/quizz-app/persistence/injector"
	"io"
	"strings"
)

type QuizzService struct {
	storage persistence.QuizzStorage
}

type MusicDto struct {
	From int `json:"from"`
	To int `json:"to"`
	Path string `json:"path"`
	Filename string
	KeepExisting bool `json:"keepExisting"`
	Delete bool `json:"delete"`
}

type QuestionDto struct {
	Id string `json:"id,omitempty"`
	Title string `json:"title"`
	Answers []model.Answer `json:"answers"`
	MaxTime int `json:"time"`
	Position int `json:"position"`
	Music MusicDto `json:"music"`
}

func (qd QuestionDto)toQuestion()model.Question{
	return model.Question{Id:qd.Id,Title: qd.Title,Answers: qd.Answers,MaxTime: qd.MaxTime,Position: qd.Position}
}

func (gs QuizzService) Update(id string,quizz model.QuizzDto) (string,error) {
	if strings.EqualFold("",quizz.Name) {
		return "",errors.New("Missing name field")
	}
	if strings.EqualFold("",id) {
		return gs.storage.Create(quizz)
	}
	return gs.storage.Update(id,quizz)
}

func (gs QuizzService)GetAll()[]model.LightQuizz {
	return gs.storage.GetAll()
}

func (gs QuizzService)Get(id string)(model.Quizz,error){
	return gs.storage.Get(id)
}

func (gs QuizzService)DeleteQuizz(id string)error{
	quizz,err := gs.storage.Get(id)
	if err != nil {
		return err
	}
	return gs.storage.DeleteQuizz(quizz)
}

func (gs QuizzService)AddQuestion(id string,questionDto QuestionDto)error{
	quizz,err := gs.Get(id)
	if err != nil {
		return err
	}
	question := questionDto.toQuestion()
	if questionDto.Music.Delete {
		// Delete existing music
		gs.storage.DeleteMusic(id,questionDto.Id)
		question.MusicPath = ""
		question.Filename = ""
	}
	if !strings.EqualFold("",questionDto.Music.Path){
		musicPath, err := gs.storage.StoreMusic(quizz,questionDto.Music.Path,questionDto.Music.From,questionDto.Music.To)
		if err != nil {
			return err
		}
		question.MusicPath = musicPath
		question.Filename = questionDto.Music.Filename
	}
	if questionDto.Music.KeepExisting && !strings.EqualFold("",questionDto.Id){
		// Keep music info, search in quizz
		copyMusicFromQuestion(quizz,&question)
	}
	return gs.storage.AddQuestion(id,question)
}

func copyMusicFromQuestion(quizz model.Quizz,question *model.Question){
	q,err := quizz.GetQuestionById(question.Id)
	if err == nil {
		question.MusicPath = q.MusicPath
		question.Filename = q.Filename
	}
}

func (gs QuizzService) DeleteQuestion(quizz model.Quizz,questionId string) {
	gs.storage.DeleteQuestion(quizz,questionId)
}

func (gs QuizzService) ReadMusic(quizz model.Quizz,questionId string, writer io.Writer) error{
	return gs.storage.ReadMusic(quizz,questionId,writer)
}

func (gs QuizzService) GetCover(quizz model.Quizz,writer io.Writer)error{
	reader, err := gs.storage.GetCover(quizz)
	if err != nil {
		return err
	}
	if _,err := io.Copy(writer,reader) ; err != nil {
		return err
	}
	return reader.Close()
}

func NewQuizzService(conf config.Config) QuizzService {
	return QuizzService{storage:injector.GetQuizzStorage(conf)}
}
