package handler

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Handler struct{
	Producer *sarama.SyncProducer
	Answer   map[string] string
}

func NewHandler(producer *sarama.SyncProducer)*Handler{
	return &Handler{
		Producer: producer,
		Answer: make(map[string]string),
	}
}
func (h *Handler)InitRoutes() *gin.Engine {
	router := gin.Default()
	api:= router.Group("api")
	{
		api.GET("/hello/:name",h.hello)
	}
	return router
}
