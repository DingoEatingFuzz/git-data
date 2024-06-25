package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type UiModel struct {
	Name   string
	Groups []Group
}

type ProgressMsg struct {
	Group   string
	BarIdx  int
	Message string
	Percent float64
	Done    bool
}

type DoneMsg struct{}

func (m UiModel) Init() tea.Cmd {
	return nil
}

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
	case DoneMsg:
		// Check if all groups and bars are done, if so, quit
		for _, g := range m.Groups {
			for _, b := range g.Bars {
				if !b.Done {
					return m, nil
				}
			}
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m UiModel) View() string {
	var groups []string

	for _, g := range m.Groups {
		groups = append(groups, g.View())
	}

	return m.Name + "\n\n" + strings.Join(groups, "\n") + "\n"
}
