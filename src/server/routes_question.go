package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/quizz-app/logger"
	"github.com/quizz-app/service"
	"io"
	"net/http"
	"os"
)

func createQuestionRoutes(server *gin.Engine) {
	server.GET("/api/quizz/:id/question/:question", addCorsUse(), getQuestionDetail)
	server.POST("/api/quizz/:id/question", addCorsUse(), IsAdmin(), createOrUpdateQuestion)
	server.DELETE("/api/quizz/:id/question/:question", addCorsUse(), IsAdmin(), deleteQuestion)
	server.OPTIONS("/api/quizz/:id/question/:question", addCorsUse(), empty)
}

func getQuestionDetail(c *gin.Context) {
	logger.GetLogger2().Info("Get question detail")
}

func empty(c *gin.Context) {}

func deleteQuestion(c *gin.Context) {
	logger.GetLogger2().Info("delete question")
	quizz, err := quizzService.Get(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "No quizz with id")
		return
	}
	quizzService.DeleteQuestion(quizz, c.Param("question"))
}

// return temp path
func copyTemp(input io.Reader) (string, error) {
	f, err := os.CreateTemp("", "temp_music_mp3")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, input)
	return f.Name(), err
}

func createOrUpdateQuestion(c *gin.Context) {
	logger.GetLogger2().Info("Create or update question")

	q := service.QuestionDto{}
	data := []byte(c.Request.FormValue("question"))
	if err := json.Unmarshal(data, &q); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	file, header, err := c.Request.FormFile("music")
	// Extract music
	if err == nil {
		q.Music.Path, err = copyTemp(file)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		q.Music.Filename = header.Filename
	}
	if err := quizzService.AddQuestion(c.Param("id"), q); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, "")
	}
}
