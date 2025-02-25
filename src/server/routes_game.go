package server

import (
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/model"
	"net/http"
	"strconv"
)

func createGameRoutes(server *gin.Engine) {
	api := server.Group("/api/game")
	api.Use(IsAdmin())
	api.Use(addCorsUse())
	api.POST("/create/:quizz_id", createGame)
	api.POST("/:game_id/start/:secure_id", startGame)
	api.GET("/:game_id/connect/:secure_id", connect)
	api.GET("/:game_id/quizz/:secure_id", getQuizzFromGame)
	api.POST("/:game_id/playNextQuestion/:secure_id", playNextQuestion)
	api.GET("/:game_id/music/:question/:secure_id", readMusic)
	api.POST("/:game_id/forceEndQuestion/:secure_id", forceEndQuestion)
	api.POST("/:game_id/computeScores/:secure_id", computeQuestionScore)
	api.GET("/:game_id/score/:secure_id", getScore)

	apiClient := server.Group("/api/player")
	apiClient.Use(isPlayer())
	apiClient.Use(addCorsUse())
	apiClient.POST("/join/:id", joinGame)
	apiClient.GET("/connect/:id/:id_player", connectGame)
	apiClient.GET("/:id/:id_player", detailPlayer)
	apiClient.POST("/answer/:id/:id_player", answerQuestion)
}

func joinGame(c *gin.Context) {
	game, err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id, idPosition, err := gameService.Join(game, c.Query("name"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "idPosition": idPosition})
}

// Create SSE connection from gameId,playerId and his name
func connectGame(c *gin.Context) {
	playerId := c.Param("id_player")
	name := c.Query("name")
	game, err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := gameService.Connect(game, playerId, name, c); err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

// Create SSE connection from gameId,playerId and his name
func detailPlayer(c *gin.Context) {
	playerId := c.Param("id_player")
	game, err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if idPosition, err := gameService.DetailPlayer(game, playerId, c); err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"position": idPosition})
	}
}

func answerQuestion(c *gin.Context) {
	game, err := gameService.GetById(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	answer, _ := strconv.Atoi(c.Query("answer"))
	gameService.AnswerQuestion(game, c.Param("id_player"), answer)
}

func getScore(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		scores := gameService.GetScore(game)
		c.JSON(http.StatusOK, scores)
	}
}

func computeQuestionScore(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		c.JSON(http.StatusOK, gameService.ComputeScore(game))
	}
}

func forceEndQuestion(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		gameService.ForceEndAnswer(game)
	}
}

func getGame(c *gin.Context) (*model.Game, bool) {
	game, err := gameService.GetByIdAndCheck(c.Param("game_id"), c.Param("secure_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return nil, false
	}
	return game, true
}

func readMusic(c *gin.Context) {
	game, ok := getGame(c)
	if !ok {
		c.String(http.StatusNotFound, "No game with id")
		return
	}
	err := quizzService.ReadMusic(game.Quizz, c.Param("question"), c.Writer)
	if err != nil {
		c.String(http.StatusNotFound, "No quizz with id")
	}
}

func playNextQuestion(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		question, more := gameService.PlayNextQuestion(game)
		if more {
			c.JSON(http.StatusOK, question)
		} else {
			c.JSON(http.StatusOK, gin.H{"end": true})
		}
	}
}

func startGame(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		gameService.Start(game, c)
	}
}

func getQuizzFromGame(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		c.JSON(http.StatusOK, game.Quizz)
	}
}

func connect(c *gin.Context) {
	game, ok := getGame(c)
	if ok {
		gameService.ConnectMaster(game, c)
	}
}

func createGame(c *gin.Context) {
	quizzId := c.Param("quizz_id")
	scoreWithTime := false
	if c.Request.FormValue("time_score") == "true" {
		scoreWithTime = true
	}
	if game, err := gameService.Create(quizzId, scoreWithTime); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"id": game.Id, "secureId": game.SecureId})
	}
}
