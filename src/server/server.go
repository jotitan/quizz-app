package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/config"
	"github.com/quizz-app/music"
	"github.com/quizz-app/service"
	"net/http"
	"path/filepath"
	"strings"
)

var quizzService service.QuizzService
var gameService service.GameService
var cutter music.Cutter


func NewServer(conf config.Config){
	initServices(conf)

	server := gin.Default()
	s := http.Server{
		Handler: server,
		Addr: fmt.Sprintf(":%s",conf.Port),
	}

	server.GET("/health",health)
	createQuizzRoutes(server)
	createQuestionRoutes(server)
	createGameRoutes(server)
	createStaticRoute(server,conf.WebResources)
	//server.Run(fmt.Sprintf(":%s",conf.Port))
	s.ListenAndServe()
}

func createStaticRoute(server *gin.Engine, resources string ){
	server.NoRoute(staticFiles(resources).serve)
}

type staticFiles string

func (s staticFiles)serve(c *gin.Context){
	url := c.Request.URL.Path[1:]
	if strings.HasPrefix(url,"static/") || filepath.Ext(url) != ""{
		http.ServeFile(c.Writer,c.Request,filepath.Join(string(s),url))
		return
	}
	http.ServeFile(c.Writer,c.Request,filepath.Join(string(s),"index.html"))
}

func initServices(conf config.Config){
	quizzService = service.NewQuizzService(conf)
	gameService = service.NewGameService(quizzService)
}

func health(c *gin.Context){
	c.JSON(200,gin.H{"status":true})
}

func addCors(handler func(*gin.Context))func(*gin.Context){
	return func(c *gin.Context) {
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods","GET, PUT, POST, DELETE, HEAD, OPTIONS")

		handler(c)
	}
}
