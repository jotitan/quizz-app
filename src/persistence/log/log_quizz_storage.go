package logpersist

import (
	"github.com/quizz-app/logger"
	"github.com/quizz-app/model"
	"github.com/quizz-app/persistence"
	"io"
)

// Store each quizz in filer for now

type logQuizzStorage struct {
}

func (fgs logQuizzStorage) Get(id string) (model.Quizz, error) {
	logger.GetLogger2().Info("get")
	return model.Quizz{},nil
}

func (fgs logQuizzStorage) Create(name string) (string, error) {
	logger.GetLogger2().Info("create")
	return "",nil
}

func (fgs logQuizzStorage) AddQuestion(id string, question model.Question) error{
	logger.GetLogger2().Info("add question")
	return nil
}

func (fgs logQuizzStorage) GetAll() []model.LightQuizz {
	return []model.LightQuizz{}
}

func (fgs logQuizzStorage) DeleteQuestion(quizz model.Quizz,questionId string) error{
	return nil
}

func (fgs logQuizzStorage) StoreMusic(quizz model.Quizz, musicFile string, from, to int) (string, error) {
	return "",nil
}

func (fgs logQuizzStorage) ReadMusic(quizz model.Quizz, questionId string, writer io.Writer) error {
	return nil
}

func (fgs logQuizzStorage) DeleteMusic(idQuizz, idQuestion string) error {
	panic("implement me")
}

func (fgs logQuizzStorage) DeleteQuizz(quizz model.Quizz) error {
	panic("implement me")
}

func NewQuizzStorage()persistence.QuizzStorage {
	return logQuizzStorage{}
}
