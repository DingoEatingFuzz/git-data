package gitjson

import (
	"dingoeatingfuzz/git.json/ui"

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
	// TODO: If progress is going to channel bytes, maybe we need lifecycle hooks
	Run(git *Git, progress func(string, float64))
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
	if len(groups[GitSource]) > 0 {
		m.Groups = append(m.Groups, ui.Group{
			Name: "Git Sources",
		})

		for range groups[GitSource] {
			m.Groups[0].Bars = append(m.Groups[0].Bars, &ui.Bar{
				Bar:     progress.New(progress.WithDefaultGradient()),
				Message: "Pending...",
			})
		}
	}

	if len(groups[GitHubSource]) > 0 {
		m.Groups = append(m.Groups, ui.Group{
			Name: "GitHub Sources",
		})

		for range groups[GitHubSource] {
			m.Groups[0].Bars = append(m.Groups[0].Bars, &ui.Bar{
				Bar:     progress.New(progress.WithDefaultGradient()),
				Message: "Pending...",
			})
		}
	}

	p = tea.NewProgram(m)

	// Yolo?
	go p.Run()

	/// if _, err := p.Run(); err != nil {
	/// 	fmt.Println("error running program:", err)
	/// 	os.Exit(1)
	/// }

	// Second pass: run the scripts
	if len(groups[GitSource]) > 0 {
		for i, s := range groups[GitSource] {
			// TODO: This should be concurrent, but we need to figure out file locking and such first
			s.Run(r.Git, func(msg string, progress float64) {
				// fmt.Println(fmt.Sprintf("Progress on %v (%d)", msg, int(progress*100)))
				p.Send(ui.ProgressMsg{
					Group:   "Git Sources",
					BarIdx:  i,
					Message: msg,
					Percent: progress,
				})
			})
		}
	}
}
