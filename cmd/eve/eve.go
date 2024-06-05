package main

import (
	"fmt"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"

	"github.com/josejulio/eve/internal/api"
	"github.com/josejulio/eve/internal/task"
)

func main() {
	viper.AddConfigPath(".")  
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Loading task files
	taskDefinition, err := task.LoadTaskFile()
	if err != nil {
		panic(fmt.Errorf("fatal error loading task file: %w", err))
	}

	llm, err := openai.New(openai.WithBaseURL(viper.Get("llm.base_url").(string)))

	if err != nil {
		panic(fmt.Errorf("fatal error loading llm: %w", err))
	}

	r := gin.Default()
	r.GET("/health", api.HealthGetAPI)
	r.GET("/talk", func(c *gin.Context) {
		api.TalkPostAPI(c, llm, *taskDefinition)
	})
	r.Run()
}
