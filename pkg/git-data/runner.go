package gitdata

import (
	"dingoeatingfuzz/git-data/ui"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Source int

const (
	GitSource Source = iota
	GitHubSource
)

// TODO: thread bubbletea in here somehow
type Script interface {
	Source() Source
	Name() string
	// TODO: If progress is going to channel bytes, maybe we need lifecycle hooks
	Run(git *Git, progress func(string, float64, bool))
}

type Runner struct {
	Scripts []Script
	Git     *Git
}

var p *tea.Program

func (r *Runner) Run() {
	groups := map[Source][]Script{}

	for _, s := range r.Scripts {
		groups[s.Source()] = append(groups[s.Source()], s)
	}

	m := ui.UiModel{}

	// First pass: initialize the model
	m.Groups = append(m.Groups, ui.Group{
		Name: "Git Sources",
	})

	for _, s := range groups[GitSource] {
		m.Groups[0].Bars = append(m.Groups[0].Bars, &ui.Bar{
			Bar:     progress.New(progress.WithDefaultGradient()),
			Name:    s.Name(),
			Message: "Pending...",
		})
	}

	m.Groups = append(m.Groups, ui.Group{
		Name: "GitHub Sources",
	})

	for _, s := range groups[GitHubSource] {
		m.Groups[1].Bars = append(m.Groups[1].Bars, &ui.Bar{
			Bar:     progress.New(progress.WithDefaultGradient()),
			Name:    s.Name(),
			Message: "Pending...",
		})
	}

	p = tea.NewProgram(m)

	// Second pass: run the scripts
	if len(groups[GitSource]) > 0 {
		// TODO: This should be concurrent, but we need to figure out file locking and such first
		go func() {
			for i, s := range groups[GitSource] {
				s.Run(r.Git, func(msg string, progress float64, done bool) {
					p.Send(ui.ProgressMsg{
						Group:   "Git Sources",
						BarIdx:  i,
						Message: msg,
						Percent: progress,
						Done:    done,
					})
				})
			}
		}()
	}

	if len(groups[GitHubSource]) > 0 {
		for i, s := range groups[GitHubSource] {
			go s.Run(r.Git, func(msg string, progress float64, done bool) {
				p.Send(ui.ProgressMsg{
					Group:   "GitHub Sources",
					BarIdx:  i,
					Message: msg,
					Percent: progress,
					Done:    done,
				})
			})
		}
	}

	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
