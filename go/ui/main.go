package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type UiModel struct {
	Groups []Group
}

func (m UiModel) Init() tea.Cmd {
	return nil
}

// This should take a specific message that is just "group" "name" "message" "percent"
// Then it should find the appropriate field in the Model and update it along with the
// corresponding progress bar

func (m UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case progress.FrameMsg:
		var cmds []tea.Cmd
		for i, g := range m.Groups {
			gModel, c := g.Update(msg)
			m.Groups[i] = gModel.(Group)
			cmds = append(cmds, c)
		}
		return m, tea.Batch(cmds...)
	case ProgressMsg:
		for i, g := range m.Groups {
			if g.Name == msg.Group {
				gModel, c := g.Update(msg)
				m.Groups[i] = gModel.(Group)
				return m, c
			}
		}
	}
	return m, nil
}

func (m UiModel) View() string {
	var groups []string

	for _, g := range m.Groups {
		groups = append(groups, g.View())
	}

	return strings.Join(groups, "\n\n")
}
