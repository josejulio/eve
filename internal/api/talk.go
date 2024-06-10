package api

import (
	"log"
//	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"

<<<<<<< Updated upstream
//	"github.com/josejulio/eve/internal/prompt"
	"github.com/josejulio/eve/internal/task"
	"github.com/josejulio/eve/internal/session"
//	"github.com/josejulio/eve/internal/actions"
	"github.com/josejulio/eve/internal/processor"
=======
	"github.com/sriroopar/eve/internal/prompt"
	"github.com/sriroopar/eve/internal/task"
>>>>>>> Stashed changes
)

func TalkPostAPI(c *gin.Context, llm llms.Model, taskDefinition task.TaskDefinition, session session.Session) {
	input := c.Query("input")
	
	response, err := processor.StepProcessor(c.Request.Context(), input, session, taskDefinition, llm)
	if err != nil {
		log.Print(err)
		c.JSON(400, err)
	} else {
		c.JSON(200, response)
	}
}
