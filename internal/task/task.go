package task

import (
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type TaskDefinition struct {
	Name string `yaml:"name"`
	Tasks map[string]Task `yaml:"tasks"`
}

type TaskStep struct {
	// Collect step
	Id string `yaml:"id"`
	TaskStepCollect `yaml:",inline"`
	TaskStepIf `yaml:",inline"`
	TaskStepAction `yaml:",inline"`
	TaskStepUtter `yaml:",inline"`
	Next string `yaml:"next"`
}

type TaskStepCollect struct {
	Name string `yaml:"name"`
	Collect string `yaml:"collect"`
	Choices interface{} `yaml:"choices"`
	Type string `yaml:"type"`
}

type TaskStepIf struct {
	If string `yaml:"if,omitempty"`
	Then []TaskStep `yaml:"then"`
	Else []TaskStep `yaml:"else"`
}

type TaskStepAction struct {
	Action string `yaml:"action,omitempty"`
}

type TaskStepUtter struct {
	Utter string `yaml:"utter,omitempty"`
}


type Task struct {
	Name string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Steps []TaskStep `yaml:"steps,omitempty"`
}


func LoadTaskFile() (*TaskDefinition, error) {
	filename, _ := filepath.Abs("./configs/tasks.yml")
    yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
        return nil, err
    }

	var taskDefinition TaskDefinition

	err = yaml.Unmarshal(yamlFile, &taskDefinition)
	if err != nil {
        return nil, err
    }

	return &taskDefinition, nil
}