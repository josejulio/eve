package prompt

import (
	"context"
	"text/template"
	"errors"
	"log"
	//"strings"
	"gopkg.in/yaml.v3"

	"github.com/tmc/langchaingo/llms"
	"github.com/josejulio/eve/internal/task"
)


var systemCollectTemplate *template.Template

type collectTemplate struct {
	Name string `yaml:"name`
	Type string `yaml:"type, omitempty`
	Description string `yaml:"description`
	Choices interface{} `yaml:"allowed,omitempty"`
}


func init() {
	prompt, err := template.ParseFiles("configs/collect_prompt.tmpl")
	if err != nil {
		panic(err)
	}

	systemCollectTemplate = prompt
}

func getCollectType(collect task.TaskStepCollect) (string) {
	var collectType string
	if collect.Type == "" {
		collectType = "string"
	} else {
		collectType = collect.Type
	}

	return collectType
}


func buildCollectTemplates(collects []task.TaskStepCollect) (string) {

	var collectTemplates []collectTemplate

	for _, collect := range collects {
		collectTemplates = append(collectTemplates, collectTemplate{
			Name: collect.Collect, 
			Description: collect.Name, 
			Type: collect.Type, 
			Choices: collect.Choices,
		})
	}

	data, _ := yaml.Marshal(collectTemplates)

	return string(data)
}

func Collect(ctx context.Context, llm llms.Model, collects []task.TaskStepCollect, query string) (map[string]interface{}, error) {


	msg, err := buildPrompt(
		systemCollectTemplate, 
		map[string]interface{} {"CollectData": buildCollectTemplates(collects),},
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

	var responseMap map[string]interface{}

	err = yaml.Unmarshal([]byte(c1.Content), &responseMap)

	if err != nil {
		return nil, err
	}

	return responseMap, nil
}