package prompt

import (
	"context"
	"text/template"
	"bytes"
	"errors"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/josejulio/eve/internal/task"
)

var systemTemplate *template.Template

type taskTemplate struct {
	Id string
	Description string
}

func init() {
	prompt, err := template.ParseFiles("configs/task_prompt.tmpl")
	if err != nil {
		panic(err)
	}

	systemTemplate = prompt
}

func buildTaskTemplates(taskDefinition task.TaskDefinition) ([]taskTemplate) {

	var taskTemplates []taskTemplate

	for taskId, task := range taskDefinition.Tasks {
		taskTemplates = append(taskTemplates, taskTemplate{Id: taskId, Description: task.Description,})
	}

	return taskTemplates
}

func Task(ctx context.Context, llm llms.Model, taskDefinition task.TaskDefinition, query string) (string, error) {

	var taskTemplate = buildTaskTemplates(taskDefinition)

	var systemPrompt bytes.Buffer
	if err := systemTemplate.Execute(&systemPrompt, map[string]interface{} {"tasks": taskTemplate,}); err != nil {
		return "", err
	}

	var systemPromptString = systemPrompt.String()

	msg := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{systemPromptString}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{query}},
		},
	}

	log.Printf("Sending message:\n - user message:\n %s\n****************\n - system message: %s\n", query, systemPromptString)

	resp, err := llm.GenerateContent(ctx, msg)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	log.Printf("Response: %s", c1.Content)
	return c1.Content, nil
}