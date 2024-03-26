package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	maxWidth = 200
)

type ProgressMsg struct {
	Group   string
	BarIdx  int
	Message string
	Percent float64
}

type Bar struct {
	Bar     progress.Model
	Message string
}

type Group struct {
	Name string
	Bars []*Bar
}

func (m Group) Init() tea.Cmd {
	return nil
}

func (m Group) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for _, b := range m.Bars {
			b.Bar.Width = msg.Width - 4
			if b.Bar.Width > maxWidth {
				b.Bar.Width = maxWidth
			}
		}
		return m, nil

	case progress.FrameMsg:
		var cmds []tea.Cmd
		for _, b := range m.Bars {
			progressModel, c := b.Bar.Update(msg)
			cmds = append(cmds, c)
			b.Bar = progressModel.(progress.Model)
		}
		return m, tea.Batch(cmds...)

	case ProgressMsg:
		// fmt.Println(fmt.Sprintf("From group: %v %d, Msg: %v", m.Name, msg.BarIdx, msg))
		if msg.Group == m.Name && len(m.Bars) > msg.BarIdx {
			b := m.Bars[msg.BarIdx]
			b.Message = msg.Message
			c := b.Bar.SetPercent(msg.Percent)
			return m, c
		}
	}

	return m, nil
}

func (m Group) View() string {
	var bars []string

	for _, b := range m.Bars {
		bars = append(bars, b.Bar.View()+"\n"+b.Message)
	}

	return m.Name + "\n" + strings.Join(bars, "\n")
}
