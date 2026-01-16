package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}
	vaultDir = fmt.Sprintf("%s/.totion", homeDir)
}

/* ---------- LIST ITEM ---------- */

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

/* ---------- MODEL ---------- */

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	list                   list.Model
	showingList            bool
	width                  int
	height                 int
}

func (m model) Init() tea.Cmd {
	return nil
}

/* ---------- UPDATE ---------- */

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h, v := docStyle.GetFrameSize()
		listHeight := msg.Height - v - fixedUIHeight()
		if listHeight < 0 {
			listHeight = 0
		}

		m.list.SetSize(msg.Width-h, listHeight)

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
			}

			if m.currentFile != nil {
				m.currentFile = nil
			}

			if m.showingList {
				if m.list.FilterState() == list.Filtering {
					break
				}
				m.showingList = false
			}

			m.noteTextArea.SetValue("")

			return m, nil

		case "ctrl+l":
			noteList := listFiles()
			m.list.SetItems(noteList)
			m.showingList = true
			return m, nil

		case "ctrl+n":
			m.createFileInputVisible = true
			m.showingList = false
			return m, nil

		case "ctrl+s":
			if m.currentFile == nil {
				break
			}

			_ = m.currentFile.Truncate(0)
			_, _ = m.currentFile.Seek(0, 0)
			_, _ = m.currentFile.WriteString(m.noteTextArea.Value())
			_ = m.currentFile.Close()

			m.currentFile = nil
			m.noteTextArea.SetValue("")
			return m, nil

		case "enter":
			if m.currentFile != nil {
				break
			}

			if m.showingList {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)

					content, err := os.ReadFile(filepath)
					
					if err != nil {
						log.Printf("Error reading file: %v", err)
						return m, nil
					}
					
					m.noteTextArea.SetValue(string(content))
					f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					
					if err != nil {
						log.Printf("Error reading file: %v", err)
						return m, nil
					}

					m.currentFile = f

					m.showingList = false
				}

				return m, nil
			}

			fileName := m.newFileInput.Value()
			if fileName != "" {
				path := fmt.Sprintf("%s/%s.md", vaultDir, fileName)
				if _, err := os.Stat(path); err == nil {
					break
				}

				f, err := os.Create(path)
				if err != nil {
					log.Fatal(err)
				}

				m.currentFile = f
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")
			}
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

/* ---------- VIEW ---------- */

func (m model) View() string {

	welcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 4)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2)

	welcome := welcomeStyle.Render("📒 Welcome to Totion")
	help := helpStyle.Render("Ctrl+N: new file | Ctrl+L: list | Esc: back | Ctrl+S: save | Ctrl+C: quit")

	var content string

	switch {
	case m.createFileInputVisible:
		content = m.newFileInput.View()
	case m.currentFile != nil:
		content = m.noteTextArea.View()
	case m.showingList:
		content = m.list.View()
	default:
		content = ""
	}

	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			welcome,
			"",
			content,
			"",
			help,
		),
	)
}

/* ---------- INIT MODEL ---------- */

func initializeModel() model {

	_ = os.MkdirAll(vaultDir, 0750)

	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.Width = 50
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.Focus()

	lst := list.New(listFiles(), list.NewDefaultDelegate(), 0, 0)
	lst.Title = "All Notes 📁"

	return model{
		newFileInput: ti,
		noteTextArea: ta,
		list:         lst,
	}
}

/* ---------- HELPERS ---------- */

func fixedUIHeight() int {
	// welcome + spacing + help box
	return 7
}

func listFiles() []list.Item {
	items := []list.Item{}

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		return items
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, err := e.Info()
		if err != nil {
			continue
		}

		items = append(items, item{
			title: e.Name(),
			desc:  "Modified: " + info.ModTime().Format("2006-01-02 15:04"),
		})
	}

	return items
}

/* ---------- MAIN ---------- */

func main() {
	p := tea.NewProgram(initializeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
