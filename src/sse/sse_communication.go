package sse

import (
	"fmt"
	"github.com/quizz-app/logger"
	"net/http"
)

type Message struct {
	Event   string
	Payload string
}

type SSECommunicate struct {
	Chanel chan Message
	writer http.ResponseWriter
	id     string
}

func NewSSECommunicate(writer http.ResponseWriter, id string) *SSECommunicate {
	return &SSECommunicate{Chanel: make(chan Message, 10), writer: writer, id: id}
}

func (sse *SSECommunicate) Run() {
	sse.writer.Header().Set("Content-Type", "text/event-stream")
	sse.writer.Header().Set("Cache-Control", "no-cache")
	sse.writer.Header().Set("Connection", "keep-alive")
	sse.writer.Header().Set("Access-Control-Allow-Origin", "*")

	sse.readMessages()
}

func (sse SSECommunicate) watchEndSSE(r *http.Request, watcher chan string) {
	go func() {
		<-r.Context().Done()
		logger.GetLogger2().Info("Stop connexion")
		//remove player
		close(sse.Chanel)
		watcher <- sse.id
	}()
}

// Blocking method which read message and send to user
func (sse *SSECommunicate) readMessages() {
	for {
		message, hasMore := <-sse.Chanel
		// End of message
		if !hasMore {
			break
		}
		sse.writer.Write([]byte(fmt.Sprintf("event: %s\n", message.Event)))
		sse.writer.Write([]byte("data: " + message.Payload + "\n\n"))
		sse.writer.(http.Flusher).Flush()
	}
}
