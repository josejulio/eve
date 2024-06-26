package processor

import(
	"context"
	"strconv"
	"text/template"
	"bytes"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/expr-lang/expr"

	"github.com/josejulio/eve/internal/prompt"
	"github.com/josejulio/eve/internal/task"
	"github.com/josejulio/eve/internal/session"
	"github.com/josejulio/eve/internal/actions"
)

func templateUtterance(utterance string, session session.Session) (string, error) {
	templ, err := template.New("utterance").Parse(utterance)

	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = templ.Execute(&output, map[string]interface{}{
		"slots": session.GetSlots(),
	})

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

const thenStepIndex = 0
const elseStepIndex = 1

func getStep(t task.Task, path []int) (*task.TaskStep, error) {
	stepPath, err := getStepsForPath(t, path)
	if err != nil {
		return nil, err
	}

	return &stepPath[len(stepPath) - 1], nil
}

func getStepsForPath(t task.Task, path []int) ([]task.TaskStep, error) {
	var step task.TaskStep = t.Steps[path[0]]
	stepPath := []task.TaskStep{step}

	path = path[1:]

	for {
		// We will consume 2 at once
		if len(path) < 1 {
			break
		}

		var elseOrThen, index int
		elseOrThen, index, path = path[0], path[1], path[2:]
		
		if step.TaskStepIf.If == "" {
			return nil, errors.New("Invalid path: Does not point to an TakeStepIf")
		}

		var subSteps []task.TaskStep
		if elseOrThen == thenStepIndex {
			subSteps = step.TaskStepIf.Then
		} else if elseOrThen == elseStepIndex {
			subSteps = step.TaskStepIf.Else
		} else {
			return nil, errors.New("Invalid path: Does not have an else/then index")
		}

		if index >= len(subSteps) {
			return nil, errors.New("Invalid path: Does not have a valid else subpath")
		}

		step = subSteps[index]
		stepPath = append(stepPath, step)
	}

	return stepPath, nil
}

func incrementStepPath(t task.Task, path []int) ([]int, error) {
	stepPath, err := getStepsForPath(t, path)
	if err != nil {
		return nil, err
	}

	lastStep := stepPath[len(stepPath) - 1]

	if lastStep.Next != "" { // Find the step with the selected id
		return stepPathForStepId(t, lastStep.Next)
	}

	for {
		if len(path) == 0 {
			return nil, errors.New("No next step found for path")
		} else if (len(path) == 1) {
			// Root level, special case
			nextIndex := path[0] + 1
			if nextIndex < len(t.Steps) {
				return []int{nextIndex}, nil
			} else {
				return nil, errors.New("No next step found for path")
			}
		} else {
			var prevStep task.TaskStep
			stepPath, prevStep = stepPath[:len(stepPath) - 1], stepPath[len(stepPath) - 1]

			if prevStep.TaskStepIf.If != "" { // If is not the If step, skip it, nothing to do at this level
				var elseOrThen, prevIndex int
				path, elseOrThen, prevIndex = path[:len(path) - 2], path[len(path) - 2], path[len(path) - 1]
		
				var subSteps []task.TaskStep
				if elseOrThen == thenStepIndex {
					subSteps = prevStep.TaskStepIf.Then
				} else if elseOrThen == elseStepIndex {
					subSteps = prevStep.TaskStepIf.Else
				} else {
					return nil, errors.New("Invalid path: Does not have an else/then index")
				}
		
				nextIndex := prevIndex + 1
		
				if nextIndex < len(subSteps) {
					newPath := make([]int, 0)
					newPath = append(newPath, path...)
					return append(newPath, elseOrThen, nextIndex), nil
				}
			}
		}		
	}

}

func stepPathForStepId(t task.Task, stepId string) ([]int, error) {
	return nil, errors.New("Not implemented")
}


func StepProcessor(ctx context.Context, input string, session session.Session, taskDefinition task.TaskDefinition, llm llms.Model) (*ProcessorResponse, error) {
	processedInput := false
	currentStepPath := session.GetStepPath()
	currentTask := session.GetTask()

	response := &ProcessorResponse{Messages: []string{}}

	for {
		incrementPath := true
		if currentTask != "" {
			// We have a task - check what's the next step and stop when requiring input from the user
			// If we haven't processed it already
			stepTask := taskDefinition.Tasks[currentTask]

			// We should support nested paths - for now we do not.
			step, err := getStep(stepTask, currentStepPath)

			if err != nil {
				return nil, err
			}

			// Process step
			if step.TaskStepCollect.Collect != "" {
				// Processing a collect step
				if processedInput {
					// Input required - break loop
					break
				}

				collectResponse, err := prompt.Collect(ctx, llm, []task.TaskStepCollect{
					step.TaskStepCollect,
				}, input)

				if err != nil {
					return nil, err
				}

				isValid := false

				value, ok := collectResponse[step.TaskStepCollect.Collect]
				if ok {
					// Do validations - our LLM could be using a wrong type
					// For now only checking if a number is an integer
					if step.TaskStepCollect.Type == "number" {
						if _, ok := value.(int); ok {
							isValid = true
						} else if strValue, ok := value.(string); ok {
							value, err = strconv.Atoi(strValue)
							if err == nil {
								isValid = true
							}
						}
					} else {
						if value != "" {
							isValid = true
						}
					}
				}

				if isValid {
					session.SetSlot(step.TaskStepCollect.Collect, value)
				} else {
					incrementPath = false
					response.Messages = append(response.Messages, "I couldn't get that. Please try again!")
				}

				processedInput = true
			} else if step.TaskStepUtter.Utter != "" {
				// Processing a utter step
				msg, err := templateUtterance(step.TaskStepUtter.Utter, session)
				if err != nil {
					return nil, err
				}

				response.Messages = append(response.Messages, msg)
			} else if step.TaskStepAction.Action != "" {
				// Processing action step
				newMessages, err := actions.ExecuteAction(step.TaskStepAction.Action, session)

				if err != nil {
					return nil, err
				}

				response.Messages = append(response.Messages, newMessages...)

			} else if step.TaskStepIf.If != "" {
				// Processing if step
				conditionResult, err := evaluateCondition(step.TaskStepIf.If, session)
				if err != nil {
					return nil, err
				}

				if conditionResult {
					if len(step.TaskStepIf.Then) > 0 {
						// Condition is true, proceed to the then step
						currentStepPath = append(currentStepPath, thenStepIndex, 0)
						incrementPath = false
					}
				} else {
					// Condition is false, move to the else step if present
					if len(step.TaskStepIf.Else) > 0 {
						currentStepPath = append(currentStepPath, elseStepIndex, 0)
						incrementPath = false
					}
				}
			}

			// Step processed - increment and check we are still in a valid step or exit
			if incrementPath {
				currentStepPath, err = incrementStepPath(stepTask, currentStepPath)
				if err != nil {
					currentTask = ""
					currentStepPath = []int{0}
					break
				}
			}
		} else {
			// Process the input to find what the user wants to do
			response, err := prompt.Task(ctx, llm, taskDefinition, input)
			
			if err != nil {
				return nil, err
			}

			processedInput = true

			// Found a task - check if it is something that we have available and start it
			if _, ok := taskDefinition.Tasks[response.Tasks[0]]; ok {
				currentTask = response.Tasks[0]
				currentStepPath = []int{0}
			}
		}
	}

	// Update session
	session.SetTask(currentTask)
	session.SetStepPath(currentStepPath)

	return response, nil
}

type ProcessorResponse struct {
	Messages []string `json:"messages"`
}

func evaluateCondition(condition string, session session.Session) (bool, error) {
    env := map[string]interface{}{
        "slots": session.GetSlots(),
    }

    program, err := expr.Compile(condition, expr.Env(env))
    if err != nil {
        return false, err
    }

    output, err := expr.Run(program, env)
    if err != nil {
        return false, err
    }

    result, ok := output.(bool)
    if !ok {
        return false, errors.New("Condition did not evaluate to a boolean")
    }

    return result, nil
}
