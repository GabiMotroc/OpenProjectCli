package appCommands

import (
	"fmt"
	"onyxide/data"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var itCmd = &cobra.Command{
	Use:   "it",
	Short: "Manage apps",
	Long:  "Create, list, and manage applications.",
	RunE:  startInteractive,
}

func init() {
	AppCmd.AddCommand(itCmd)
}

func startInteractive(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}

func initialModel() model {
	a, _ := data.LoadApps()

	ti := textinput.New()
	ti.CharLimit = 156
	ti.SetWidth(20)
	ti.SetVirtualCursor(false)
	return model{
		apps:  a,
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		if m.inputting {
			switch msg.String() {
			case "enter":
				name := m.input.Value()
				if name != "" {
					if m.editingIndex >= 0 {
						m.apps[m.editingIndex].Name = name
					} else {
						m.apps = append(m.apps, data.App{Name: name})
						m.cursor = len(m.apps) - 1
					}
				}
				m.inputting = false
				m.editingIndex = -1
				m.input.SetValue("")
			case "esc":
				m.inputting = false
				m.input.SetValue("")
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.apps)-1 {
					m.cursor++
				}

			case "a":
				m.inputting = true
				m.input.SetValue("")
				cmd := m.input.Focus()
				cmds = append(cmds, cmd)
			case "e":
				if len(m.apps) > 0 {
					m.inputting = true
					m.editingIndex = m.cursor
					m.input.SetValue(m.apps[m.cursor].Name)
					cmd := m.input.Focus()
					cmds = append(cmds, cmd)
				}
			case "s":
				err := data.SaveApps(m.apps)
				if err != nil {
					return nil, nil
				}
				return m, tea.Quit

			case "d":
				m.apps = append(m.apps[:m.cursor], m.apps[m.cursor+1:]...)
				if m.cursor > 0 {
					m.cursor--
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.

	//return m, nil

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	s := "Manage apps\n\n"

	for i, app := range m.apps {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%d) %s %s\n", i, cursor, app.Name)
	}

	if m.inputting {
		s += "\n" + m.input.View() + "\nenter to confirm, esc to cancel"
	} else {
		s += "\nPress q to quit. Press a to add. Press s to save. Press d to delete\n"
	}
	return tea.NewView(s)
}

type model struct {
	apps         []data.App
	cursor       int
	input        textinput.Model
	inputting    bool
	editingIndex int
}
