package api

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/josejulio/eve/internal/prompt"
)

func TalkPostAPI(c *gin.Context, llm llms.Model) {
	input := c.Query("input")

	response, err := prompt.Task(c.Request.Context(), llm, input)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, err)
	}

	c.JSON(200, response)
}
