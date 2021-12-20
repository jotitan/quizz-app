package server

import (
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/logger"
	"github.com/quizz-app/model"
	"net/http"
	"strconv"
)

func createGameRoutes(server *gin.Engine) {
	api := server.Group("/api")
	api.POST("/game/create/:quizz_id", addCors(createGame))
	api.POST("/game/:game_id/start/:secure_id", addCors(startGame))
	api.GET("/game/:game_id/connect/:secure_id", addCors(connect))
	api.GET("/game/:game_id/quizz/:secure_id", addCors(getQuizzFromGame))
	api.POST("/game/:game_id/playNextQuestion/:secure_id", addCors(playNextQuestion))
	api.GET("/game/:game_id/music/:question/:secure_id", addCors(readMusic))
	api.POST("/game/:game_id/forceEndQuestion/:secure_id", addCors(forceEndQuestion))
	api.POST("/game/:game_id/computeScores/:secure_id", addCors(computeQuestionScore))
	api.GET("/game/:game_id/score/:secure_id", addCors(getScore))

	api.POST("/player/join/:id", addCors(joinGame))
	api.GET("/player/connect/:id/:id_player", addCors(connectGame))
	api.POST("/player/answer/:id/:id_player", addCors(answerQuestion))
}

func joinGame(c *gin.Context){
	game,err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
		return
	}
	id,err := gameService.Join(game,c.Query("name"))
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
		return
	}
	c.JSON(http.StatusOK,gin.H{"id":id})
}

// Create SSE connection from gameId,playerId and his name
func connectGame(c *gin.Context){
	playerId := c.Param("id_player")
	name := c.Query("name")
	game,err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
		return
	}
	if err := gameService.Connect(game,playerId,name,c) ; err != nil {
		c.String(http.StatusNotFound,err.Error())
	}
}

func answerQuestion(c *gin.Context){
	game,err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
		return
	}
	answer,_ := strconv.Atoi(c.Query("answer"))
	gameService.AnswerQuestion(game,c.Param("id_player"),answer)
}

func getScore(c *gin.Context){
	game,ok := getGame(c)
	if ok {
		scores := gameService.GetScore(game)
		c.JSON(http.StatusOK,scores)
	}
}

func computeQuestionScore(c *gin.Context){
	game,ok := getGame(c)
	if ok{
		c.JSON(http.StatusOK,gameService.ComputeScore(game))
	}
}

func forceEndQuestion(c *gin.Context){
	game,ok := getGame(c)
	if ok{
		gameService.ForceEndAnswer(game)
	}
}

func getGame(c *gin.Context)(*model.Game,bool){
	game,err := gameService.GetByIdAndCheck(c.Param("game_id"),c.Param("secure_id"))
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
		return nil,false
	}
	return game,true
}

func readMusic(c *gin.Context){
	game,ok := getGame(c)
	if !ok {
		c.String(http.StatusNotFound,"No game with id")
		return
	}
	err := quizzService.ReadMusic(game.Quizz,c.Param("question"),c.Writer)
	if err != nil {
		c.String(http.StatusNotFound,"No quizz with id")
	}
}

func playNextQuestion(c *gin.Context){
	game,ok := getGame(c)
	if ok{
		question,more :=gameService.PlayNextQuestion(game)
		if more {
			c.JSON(http.StatusOK,question)
		}else{
			c.JSON(http.StatusOK,gin.H{"end":true})
		}
	}
}

func startGame(c *gin.Context){
	game,ok := getGame(c)
	if ok {
		gameService.Start(game,c)
	}
}

func getQuizzFromGame(c *gin.Context){
	game,ok := getGame(c)
	if ok {
		c.JSON(http.StatusOK, game.Quizz)
	}
}

func connect(c *gin.Context){
	game,ok := getGame(c)
	if ok {
		gameService.ConnectMaster(game,c)
	}
}

func createGame(c *gin.Context){
	quizzId := c.Param("quizz_id")
	if game,err := gameService.Create(quizzId) ; err != nil {
		c.String(http.StatusBadRequest,err.Error())
	}else{
		logger.GetLogger2().Info("EditGame",quizzId,"created",game)
		c.JSON(http.StatusOK,gin.H{"id":game.Id,"secureId":game.SecureId})
	}
}
