package model

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
)

type QuizzDto struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	RemoveImage            bool   `json:"removeImage"`
	ImageDescription       multipart.File
	ImageDescriptionHeader *multipart.FileHeader
}

type Quizz struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       bool       `json:"image"`
	Id          string     `json:"id"`
	Questions   []Question `json:"questions"`
}

func (q Quizz) GetQuestionById(id string) (Question, error) {
	pos, err := q.GetPositionQuestionById(id)
	if err == nil {
		return q.Questions[pos], nil
	}
	return Question{}, err
}

func (q Quizz) GetPositionQuestionById(id string) (int, error) {
	for i, question := range q.Questions {
		if strings.EqualFold(question.Id, id) {
			return i, nil
		}
	}
	return 0, errors.New("impossible to find question")
}

func (q Quizz) GetNextId() string {
	max := -1
	for _, question := range q.Questions {
		value, _ := strconv.Atoi(question.Id)
		if value > max {
			max = value
		}
	}
	return fmt.Sprintf("%d", max+1)
}

type LightQuizz struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       bool   `json:"img"`
	Id          string `json:"id"`
	Nb          int    `json:"nb"`
}

func NewLightQuizz(quizz Quizz) *LightQuizz {
	return &LightQuizz{
		Id:          quizz.Id,
		Description: quizz.Description,
		Name:        quizz.Name,
		Image:       quizz.Image,
		Nb:          len(quizz.Questions),
	}
}

type Question struct {
	Id        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	Answers   []Answer `json:"answers"`
	MaxTime   int      `json:"time"`
	Position  int      `json:"position"`
	MusicPath string   `json:"path"`
	Filename  string   `json:"filename"`
}

func (q Question) GetGoodAnswers() map[int]string {
	good := make(map[int]string)
	for pos, answer := range q.Answers {
		if answer.Good {
			good[pos] = answer.Text
		}
	}
	return good
}

type Answer struct {
	Text string `json:"text"`
	Good bool   `json:"good"`
}

type LightQuestion struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

func ToLightQuizz(quizzes []Quizz) []LightQuizz {
	lg := make([]LightQuizz, len(quizzes))
	for pos, g := range quizzes {
		lg[pos] = LightQuizz{Id: g.Id, Name: g.Name}
	}
	return lg
}

func ToLightQuestions(questions []Question) []LightQuestion {
	lq := make([]LightQuestion, len(questions))
	for pos, q := range questions {
		lq[pos] = LightQuestion{Id: q.Id, Title: q.Title}
	}
	return lq
}
