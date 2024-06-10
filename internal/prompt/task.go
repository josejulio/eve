package prompt

import (
	"context"
	"text/template"
	"errors"
	"log"
	"gopkg.in/yaml.v3"

	"github.com/tmc/langchaingo/llms"
	"github.com/josejulio/eve/internal/task"
)

var systemTaskTemplate *template.Template

type taskTemplate struct {
	Id string
	Description string
}

type TaskResponse struct {
	Tasks []string `yaml:"tasks"`
}

func init() {
	prompt, err := template.ParseFiles("configs/task_prompt.tmpl")
	if err != nil {
		panic(err)
	}

	systemTaskTemplate = prompt
}

func buildTaskTemplates(taskDefinition task.TaskDefinition) ([]taskTemplate) {

	var taskTemplates []taskTemplate

	for taskId, task := range taskDefinition.Tasks {
		taskTemplates = append(taskTemplates, taskTemplate{Id: taskId, Description: task.Description,})
	}

	return taskTemplates
}

func Task(ctx context.Context, llm llms.Model, taskDefinition task.TaskDefinition, query string) (*TaskResponse, error) {

	var taskTemplate = buildTaskTemplates(taskDefinition)

	msg, err := buildPrompt(
		systemTaskTemplate, 
		map[string]interface{} {"tasks": taskTemplate,},
		query,
	)

	if err != nil {
		return nil, err
	}

	resp, err := llm.GenerateContent(ctx, msg)
	if err != nil {
		return nil, err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return nil, errors.New("empty response from model")
	}
	c1 := choices[0]
	log.Printf("Response: %s", c1.Content)


	var taskResponse TaskResponse
	err = yaml.Unmarshal([]byte(c1.Content), &taskResponse)

	if err != nil {
		return nil, err
	}

	return &taskResponse, nil
}