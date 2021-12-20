package test

import (
	"github.com/quizz-app/config"
	"github.com/quizz-app/service"
	"os"
	"strings"
	"testing"
)

var tempDirOther = os.TempDir()

var conf = config.Config{Storage: "filer",FilerPath: tempDirOther}
var quizzService = service.NewQuizzService(conf)
var gameService = service.NewGameService(quizzService)


func TestCreateGame(t *testing.T) {
	id,_ := quizzService.Create("test game")
	game,err := gameService.Create(id)
	if err != nil {
		t.Error("Impossible to create game",err)
		return
	}
	if len(game.Id) != 4 {
		t.Error("Must have a for letter game")
	}
	if _,err := gameService.GetById(game.Id) ; err != nil {
		t.Error("Must found game with given id")
	}
}

func TestConnectGame(t *testing.T) {
	id,_ := quizzService.Create("test game2")
	game,_ := gameService.Create(id)
	id,err := gameService.Join(game,"super")
	if err != nil {
		t.Error("Must connect")
	}
	if strings.EqualFold("",id) {
		t.Error("Id must exist")
	}

}




