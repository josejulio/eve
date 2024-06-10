package prompt

import (
	"text/template"
	"bytes"
	"log"

	"github.com/tmc/langchaingo/llms"
)

func buildPrompt(promptTemplate *template.Template, arguments map[string]interface{}, query string) ([]llms.MessageContent, error) {

	var prompt bytes.Buffer
	if err := promptTemplate.Execute(&prompt, arguments); err != nil {
		return nil, err
	}

	var promptString = prompt.String()

	log.Printf("Built prompt:\n - user message:\n %s\n****************\n - system message: %s\n", query, promptString)

	return []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{promptString}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{query}},
		},
	}, nil
}
