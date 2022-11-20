package persistence

import (
	"github.com/quizz-app/model"
	"io"
)

// Store each quizz in filer for now

type QuizzStorage interface {
	Create(dto model.QuizzDto) (string, error)
	Update(id string, dto model.QuizzDto) (string, error)
	Get(id string) (model.Quizz, error)
	GetCover(quizz model.Quizz) (io.ReadCloser, error)
	AddQuestion(id string, question model.Question) error
	GetAll() []model.LightQuizz
	DeleteQuestion(quizz model.Quizz, questionId string) error
	ReadMusic(quizz model.Quizz, questionId string, writer io.Writer) error
	StoreMusic(quizz model.Quizz, musicFile string, from, to int) (string, string, error)
	DeleteMusic(idQuizz, idQuestion string) error
	DeleteQuizz(quizz model.Quizz) error
}
