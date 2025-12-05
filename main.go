package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}
	vaultDir = fmt.Sprintf("%s/.totion", homeDir)
}

type model struct {
	newFIleInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
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
		case "ctrl+s":
			// save
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Cannot save the file :( ")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("Cannot save the file :( ")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("Cannot save the file :( ")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("Cannot save the file :( ")
				return m, nil
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil
		case "enter":
			if m.currentFile != nil {
				break
			}
			// todo: create file
			fileName := m.newFIleInput.Value()
			if fileName != "" {
				filepath := fmt.Sprintf("%s/%s.md", vaultDir, fileName)

				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}

				f, err := os.Create(filepath)
				if err != nil {
					log.Fatalf("%v", err)
				}

				m.currentFile = f
				m.createFileInputVisible = false
				m.newFIleInput.SetValue("")
			}
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFIleInput, cmd = m.newFIleInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
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

	help := HelpStyle.Render("Ctrl+N: new file | Ctrl+L: list files | Esc: back/save | Ctrl+S: save | Ctrl+C: quit")

	view := ""

	if m.createFileInputVisible {
		view = m.newFIleInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initializeModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// initialize new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	// textarea
	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.Focus()

	return model{
		newFIleInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
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
