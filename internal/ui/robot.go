package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	robotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Pinkish robot
	textStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Grey text
	doneStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green for success
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red for error
)

// Robot Frames
var frames = []string{
	`
      \ /
     [o_o]
     /| |\
     /   \
    `,
	`
      | |
     [o_o]
     /| |\
     /   \
    `,
	`
      / \
     [o_o]
     /| |\
     /   \
    `,
	`
      | |
     [-_-]
     /| |\
     /   \
    `,
}

type model struct {
	frameIndex int
	task       func() (string, error)
	output     string
	err        error
	done       bool
	width      int
	height     int
}

type tickMsg time.Time
type taskFinishedMsg struct {
	output string
	err    error
}

func initialModel(task func() (string, error)) model {
	return model{
		frameIndex: 0,
		task:       task,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		runTask(m.task),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if m.done {
			return m, nil
		}
		m.frameIndex = (m.frameIndex + 1) % len(frames)
		return m, tickCmd()

	case taskFinishedMsg:
		m.done = true
		m.output = msg.output
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		// We don't print the result here to keep the TUI clean.
		// The result is returned by RunProgram.
		return ""
	}

	// Robot Animation
	robot := robotStyle.Render(frames[m.frameIndex])

	// Status Text
	status := textStyle.Render("Processing...")

	return fmt.Sprintf("\n%s\n\n%s\n\n", robot, status)
}

// Commands
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func runTask(task func() (string, error)) tea.Cmd {
	return func() tea.Msg {
		out, err := task()
		return taskFinishedMsg{output: out, err: err}
	}
}

// RunProgram runs the TUI with the given task and returns the result.
func RunProgram(task func() (string, error)) (string, error) {
	p := tea.NewProgram(initialModel(task))
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	finalModel := m.(model)
	return finalModel.output, finalModel.err
}
