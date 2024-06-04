package prompt

import (
	"context"
	"text/template"
	"bytes"
	"errors"

	"github.com/tmc/langchaingo/llms"
)

var systemTemplate *template.Template

func init() {
	prompt, err := template.ParseFiles("configs/task_prompt.tmpl")
	if err != nil {
		panic(err)
	}

	systemTemplate = prompt
}

func Task(ctx context.Context, llm llms.Model, query string) (string, error) {

	var systemPrompt bytes.Buffer
	if err := systemTemplate.Execute(&systemPrompt, nil); err != nil {
		return "", err
	}

	msg := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{systemPrompt.String()}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{query}},
		},
	}

	resp, err := llm.GenerateContent(ctx, msg)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}