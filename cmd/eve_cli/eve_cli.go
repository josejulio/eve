package main

import (
    "fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    messages  []string
    textInput textinput.Model

}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "What do you want to do?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 44


	return model{
		// Our to-do list is a grocery list
		messages:  []string{},
		textInput: ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter":
			m.messages = append(m.messages, "User: " + m.textInput.Value())
			m.textInput.SetValue("")
        }
    }

	m.textInput, cmd = m.textInput.Update(msg)
    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, cmd
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
    // The header
    s := "EVE command line?\n\n"

    // Iterate over our choices
    for _, message := range m.messages {
        // Render the row
        s += fmt.Sprintf("%s\n", message)
    }

	s += "\n\n" + m.textInput.View()

    // The footer
    s += "\n\nPress ctrl+c to quit.\n"

    // Send the UI for rendering
    return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}