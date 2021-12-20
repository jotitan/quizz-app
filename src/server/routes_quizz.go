package server

import (
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/logger"
	"github.com/quizz-app/model"
	"net/http"
)

func createQuizzRoutes(server *gin.Engine){
	api := server.Group("/api")
	api.GET("/quizz/:id", addCors(getQuizz))
	api.GET("/quizzes", addCors(getQuizzes))
	api.GET("/quizz/:id/questions", addCors(getQuestionsOfQuizz))
	api.POST("/quizz", addCors(createQuizz))
	api.DELETE("/quizz/:id", addCors(deleteQuizz))
	api.OPTIONS("/quizz/:id", addCors(empty))
}

func getQuizzes(c *gin.Context) {
	c.JSON(http.StatusOK, quizzService.GetAll())
}

func getQuizz(c *gin.Context){
	logger.GetLogger2().Info("Get quizz")
	quizz,err := quizzService.Get(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound,"No quizz with id")
		return
	}
	c.JSON(http.StatusOK, quizz)
}

func getQuestionsOfQuizz(c *gin.Context){
	logger.GetLogger2().Info("Get questions")
	quizz,err := quizzService.Get(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound,"No quizz with id")
		return
	}
	c.JSON(http.StatusOK,model.ToLightQuestions(quizz.Questions))
}

func deleteQuizz(c *gin.Context){
	if err := quizzService.DeleteQuizz(c.Param("id")) ; err != nil {
		c.String(http.StatusNotFound,"Impossible to delete quizz")
		return
	}
	c.JSON(http.StatusOK,"")
}

func createQuizz(c *gin.Context){
	logger.GetLogger2().Info("Create quizz")
	name := c.Query("name")
	if id,err := quizzService.Create(name) ; err !=nil {
		c.String(http.StatusBadRequest,err.Error())
	}else{
		c.JSON(http.StatusCreated,gin.H{"id":id})
	}
}
