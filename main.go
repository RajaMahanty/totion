package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

type model struct {
	newFIleInput           textinput.Model
	createFileInputVisible bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFIleInput, cmd = m.newFIleInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	WelcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 4).
		Margin(1)

	HelpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginTop(1)

	welcome := WelcomeStyle.Render("📒 Welcome to Totion")
	view := ""
	if m.createFileInputVisible {
		view = m.newFIleInput.View()
	}
	help := HelpStyle.Render("Ctrl+N: new file | Ctrl+L: list files | Esc: back/save | Ctrl+S: save | Ctrl+C: quit")

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initializeModel() model {
	// initialize new files
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	return model{
		newFIleInput:           ti,
		createFileInputVisible: false,
	}
}

func main() {
	fmt.Print("Welcome to Totion!")
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there is an error: %v\n", err)
		os.Exit(1)
	}
}
