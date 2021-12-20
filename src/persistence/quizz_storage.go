package persistence

import (
	"github.com/quizz-app/model"
	"io"
)

// Store each quizz in filer for now

type QuizzStorage interface {
	Create(name string)(string,error)
	Get(id string)(model.Quizz,error)
	AddQuestion(id string, question model.Question)error
	GetAll() []model.LightQuizz
	DeleteQuestion(quizz model.Quizz,questionId string)error
	ReadMusic(quizz model.Quizz,questionId string, writer io.Writer)error
	StoreMusic(quizz model.Quizz, musicFile string, from,to int)(string,error)
	DeleteMusic(idQuizz,idQuestion string)error
	DeleteQuizz(quizz model.Quizz) error
}
