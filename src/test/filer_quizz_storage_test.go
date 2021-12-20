package test

import (
	"fmt"
	"github.com/quizz-app/model"
	"github.com/quizz-app/persistence/filer"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tempDir = os.TempDir()

func TestRand(t *testing.T){
	fmt.Println(fmt.Sprintf("%c",65))
}

func TestCreateQuizz(t *testing.T){
	storage := filer.NewQuizzStorage(tempDir)

	id,err := storage.Create("test")
	if err != nil {
		t.Error("No error possible",err)
		return
	}
	pathToTest := filepath.Join(tempDir,"quizz." + id + ".json")
	f,err := os.Open(pathToTest)
	if err != nil {
		t.Error("Impossible to open file",err)
		return
	}
	f.Close()
}

func TestGetQuizz(t *testing.T){
	storage := filer.NewQuizzStorage(tempDir)

	id,_ := storage.Create("test2")
	quizz,err := storage.Get(id)
	if err != nil {
		t.Error("Impossible to get quizz",err)
		return
	}
	if !strings.EqualFold(quizz.Id,id) {
		t.Error("Content must be the same but find", quizz.Id)
	}
	if !strings.EqualFold(quizz.Name,"test2") {
		t.Error("Content must be the same but find", quizz.Name)
	}
	if _,err := storage.Get("toto") ; err == nil {
		t.Error("Unknown quizz must fail")
	}
}

func TestAddQuestion(t *testing.T){
	storage := filer.NewQuizzStorage(tempDir)
	id,_ := storage.Create("test-with-question")
	err := storage.AddQuestion(id,model.Question{Title:"New question",Answers:[]model.Answer{{Text: "First response",Good:true},{Text: "Bad response"},{Text: "Second bad"}}})
	if err != nil {
		t.Error("No error can appear")
	}
	quizz,_ := storage.Get(id)
	if size := len(quizz.Questions) ; size != 1 {
		t.Error("Must have 1 question but found",size)
	}

}

func TestRemoveQuestion(t *testing.T) {
	storage := filer.NewQuizzStorage(tempDir)
	id, _ := storage.Create("test-with-question2")
	storage.AddQuestion(id, model.Question{Title: "New question", Answers: []model.Answer{}})
	storage.AddQuestion(id, model.Question{Title: "New question2", Answers: []model.Answer{}})
	storage.AddQuestion(id, model.Question{Title: "New question3", Answers: []model.Answer{}})
	storage.AddQuestion(id, model.Question{Title: "New question4", Answers: []model.Answer{}})
	quizz,_ := storage.Get(id)
	if size := len(quizz.Questions); size != 4 {
		t.Error("Must have 4 question but found", size)
	}
	storage.DeleteQuestion(quizz,quizz.Questions[2].Id)
	quizz,_ = storage.Get(id)
	if size := len(quizz.Questions); size != 3 {
		t.Error("Must have 3 question but found", size)
	}
	storage.DeleteQuestion(quizz,quizz.Questions[0].Id)
	quizz,_ = storage.Get(id)
	if size := len(quizz.Questions); size != 2 {
		t.Error("Must have 2 question but found", size)
	}
	storage.DeleteQuestion(quizz,quizz.Questions[1].Id)
	quizz,_ = storage.Get(id)
	if size := len(quizz.Questions); size != 1 {
		t.Error("Must have 1 question but found", size)
	}
}
