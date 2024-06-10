package main

import (
    "encoding/json"
    "fmt"
    "net/http"
	"net/url"
    "io/ioutil"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
    messages  []string
    textInput textinput.Model
	loading bool
	err string
	spinner  spinner.Model
}

type response struct {
	Messages []string `json:"messages"`
}

type errMsg struct {
    err error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "What do you want to do?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 44

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		
		messages:  []string{},
		textInput: ti,
		loading: false,
		err: "",
		spinner: s,
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
			if m.loading == false {
				m.loading = true
				m.err = ""
				return m, tea.Batch(
					sendMessage(m.textInput.Value()),
					m.spinner.Tick, 
				)
			}
        }

		case response:
			m.loading = false
			m.messages = append(m.messages, "User: " + m.textInput.Value())
			m.textInput.SetValue("")
			for _, message := range msg.Messages {
				m.messages = append(m.messages, "Eve: " + message)
			}
		
		case errMsg:
			m.loading = false
			m.err = fmt.Sprintf("Failed to send the message - please try again. (%v)", msg.err)
    }
	

	if m.loading == false {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.spinner, cmd = m.spinner.Update(msg)
	}

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, cmd
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
    // The header
    s := "Eve command line utility\n\n"

    // Iterate over our choices
    for _, message := range m.messages {
        // Render the row
        s += fmt.Sprintf("%s\n", message)
    }

	s += "\n\n" + m.textInput.View()

	if m.err != "" {
		s += "\n" + m.err
	}

	if m.loading {
		s += "\n" + m.spinner.View() + " Loading..."
	}

    // The footer
    s += "\n\nPress ctrl+c to quit.\n"

    // Send the UI for rendering
    return s
}

func sendMessage(message string) tea.Cmd {
	return func() tea.Msg {

		baseURL, err := url.Parse("http://localhost:8080/talk") // Replace with your URL
        if err != nil {
            return errMsg{err}
        }

		params := url.Values{}
        params.Add("input", message)
		baseURL.RawQuery = params.Encode()


		resp, err := http.Get(baseURL.String())
		if err != nil {
			return errMsg{err}
		}


		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errMsg{err}
		}

		if resp.StatusCode != http.StatusOK {
            return errMsg{fmt.Errorf("received non-200 status code: %d", resp.StatusCode)}
        }

		var eveResp response
		if err = json.Unmarshal(body, &eveResp); err != nil {
			return errMsg{err}
		}

		return eveResp
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}