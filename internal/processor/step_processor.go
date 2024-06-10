package processor

import(
	"context"
	"strconv"
	"errors"

	"github.com/tmc/langchaingo/llms"

	"github.com/josejulio/eve/internal/prompt"
	"github.com/josejulio/eve/internal/task"
	"github.com/josejulio/eve/internal/session"
	"github.com/josejulio/eve/internal/actions"

)

func StepProcessor(ctx context.Context, input string, session session.Session, taskDefinition task.TaskDefinition, llm llms.Model) (*ProcessorResponse, error) {
	processedInput := false
	currentStepPath := session.GetStepPath()
	currentTask := session.GetTask()

	response := &ProcessorResponse{Messages: []string{}}

	for {		
		if currentTask != "" {
			// We have a task - check what's the next step and stop when requiring input from the user
			// If we haven't processed it already
			steps := taskDefinition.Tasks[currentTask].Steps

			// We should support nested paths - for now we do not.
			step := steps[currentStepPath[0]]

			// Process step
			if step.TaskStepCollect.Collect != "" {
				// Processing a collect step
				if processedInput {
					// Input required - break loop
					break
				}

				response, err := prompt.Collect(ctx, llm, []task.TaskStepCollect{
					step.TaskStepCollect,
				}, input)

				if err != nil {
					return nil, err
				}

				isValid := false

				value, ok := response[step.TaskStepCollect.Collect]
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
						isValid = true
					}
				}

				if isValid {
					session.SetSlot(step.TaskStepCollect.Collect, value)
				}

				processedInput = true
			} else if step.TaskStepUtter.Utter != "" {
				// Processing a utter step
				response.Messages = append(response.Messages, step.TaskStepUtter.Utter)
			} else if step.TaskStepAction.Action != "" {
				// Processing action step
				// ToDo: Allow to write messages
				newMessages, err := actions.ExecuteAction(step.TaskStepAction.Action, session)

				if err != nil {
					return nil, err
				}

				response.Messages = append(response.Messages, newMessages...)

			} else if step.TaskStepIf.If != "" {
				// Processing if step
				return nil, errors.New("Not implemented yet")
			}

			// Step processed - increment and check we are still in a valid step or exit
			currentStepPath = []int{currentStepPath[0]+1}
			if currentStepPath[0] >= len(steps) {
				// End of flow - break flow if there is no other flow
				currentTask = ""
				currentStepPath = []int{0}
				break
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
	Messages []string `json:" messages"`
}
