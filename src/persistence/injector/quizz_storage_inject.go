package injector

import (
	"github.com/quizz-app/config"
	"github.com/quizz-app/persistence"
	"github.com/quizz-app/persistence/filer"
	logpersist "github.com/quizz-app/persistence/log"
)

func GetQuizzStorage(conf config.Config)persistence.QuizzStorage {
	switch conf.Storage {
	case "filer":return filer.NewQuizzStorage(conf.FilerPath,conf.FfmpegPath)
	default:return logpersist.NewQuizzStorage()
	}
}
