package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MyMessage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Comeback string `json:"comeback"`
}

func (h *Handler) hello(c *gin.Context) {

	requestID := uuid.New().String()

	name := c.Param("name")
	message := MyMessage{
		ID:       "1",
		Name:     name,
		Value:    "hello" + name,
		Comeback: requestID,
	}
	// Преобразование сообщения в JSON что бы потом отправить через kafka
	bytes, err := json.Marshal(message)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal JSON"})
		return
	}
	// пишем сообщение в topic nikita
	msg := &sarama.ProducerMessage{
		Topic: "nikita",
		Key:   sarama.StringEncoder(message.Comeback),
		Value: sarama.ByteEncoder(bytes),
	}

	// отправка сообщения в Kafka
	pr := *h.Producer
	_, _, err = pr.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return
	}

	// когда api gateway получит сообщение обратно, нам надо понять какому типу ее отравлять
	// для этого в этой горутине мы будем смотреть мапу, ли в ней наша запись с нашим id или нет
	for {
		time.Sleep(10 * time.Second)
		fmt.Println(h.Answer[requestID])

		if val, ok := h.Answer[requestID]; ok {
			fmt.Println("Got this message from ApiGateway(main) -> ApiGateway(hello) whit map ", val)
			c.JSON(http.StatusOK, map[string]interface{}{
				"result": val,
			})
			delete(h.Answer, requestID)
			return
		}
	}

}
