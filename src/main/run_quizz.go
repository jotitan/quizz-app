package main

import (
	"github.com/quizz-app/config"
	"github.com/quizz-app/server"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Specify conf path as argument")
	}
	conf, err := config.ReadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	server.NewServer(*conf)
}
