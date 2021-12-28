package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/logger"
	"github.com/quizz-app/model"
	"net/http"
)

func createQuizzRoutes(server *gin.Engine){
	api := server.Group("/api")
	api.GET("/quizz/:id", addCors(getQuizz))
	api.GET("/quizz/:id/cover", addCors(getCover))
	api.GET("/quizzes", addCors(getQuizzes))
	api.GET("/quizz/:id/questions", addCors(getQuestionsOfQuizz))
	api.POST("/quizz", addCors(createQuizz))
	api.POST("/quizz/:id", addCors(updateQuizz))
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

func getCover(c *gin.Context){
	logger.GetLogger2().Info("Get quizz")
	quizz,err := quizzService.Get(c.Param("id"))
	if err != nil || !quizz.Image{
		c.String(http.StatusNotFound,"No quizz with id")
		return
	}
	err = quizzService.GetCover(quizz,c.Writer)
	if err != nil{
		c.String(http.StatusNotFound,err.Error())
	}
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
	createOrUpdateQuizz(c,"","Create")
}

func updateQuizz(c *gin.Context){
	createOrUpdateQuizz(c,c.Param("id"),"Update")
}

func createOrUpdateQuizz(c *gin.Context,id,method string){
	logger.GetLogger2().Info(method + " quizz")
	quizz := model.QuizzDto{}
	json.Unmarshal([]byte(c.Request.FormValue("quizz")),&quizz)
	if file,header,err := c.Request.FormFile("image") ; err == nil {
		quizz.ImageDescription = file
		quizz.ImageDescriptionHeader = header
	}

	if id,err := quizzService.Update(id,quizz) ; err !=nil {
		c.String(http.StatusBadRequest,err.Error())
	}else{
		c.JSON(http.StatusCreated,gin.H{"id":id})
	}
}
