package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Enter Forth code here..."
	ti.Focus()

	return model{
		textarea:       ti,
		filePicker:     makeFilePicker(),
		output:         "",
		width:          80,
		height:         24,
		err:            nil,
		showFilePicker: false,
	}
}

func makeFilePicker() filepicker.Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".fth"} // Only .fth files (but dirs still show)

	fp.SetHeight(20)

	cwd, _ := os.Getwd()
	fp.CurrentDirectory = cwd

	return fp
}

func makeWordsTable() table.Model {
	columns := []table.Column{
		{Title: "Word", Width: 15},
		{Title: "Type", Width: 12},
		{Title: "Description", Width: 50},
	}

	rows := []table.Row{}

	// Add builtins
	for name, word := range Builtins {
		rows = append(rows, table.Row{name, word.Category, word.Description})
	}

	// Add user-defined words
	for name, word := range UserWords {
		rows = append(rows, table.Row{name, word.Category, word.Description})
	}

	if len(rows) == 0 {
		rows = append(rows, table.Row{"No words defined", "", ""})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(20),
		table.WithFocused(true),
	)

	t.SetStyles(GetDictionaryTableStyle())

	return t
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.showFilePicker {
			m.filePicker.SetHeight(m.height - 5)
		}

	case tea.KeyMsg:
		// Handle file picker
		// Handle file picker - ALWAYS update it first
		if m.showFilePicker {
			m.filePicker, cmd = m.filePicker.Update(msg)
			cmds = append(cmds, cmd)

			switch msg.Type {
			case tea.KeyEnter:
				// Check if a file was selected
				selected, path := m.filePicker.DidSelectFile(msg)
				if selected {
					data, err := os.ReadFile(path)
					if err != nil {
						output.Write(fmt.Sprintf("Error loading: %s", err))
					} else {
						m.textarea.SetValue(string(data))
						m.currentFile = path
						output.Write(fmt.Sprintf("Loaded from %s", path))
					}
					m.output = output.String()
					m.showFilePicker = false
					return m, tea.Batch(cmds...)
				}
			case tea.KeyEsc:
				m.showFilePicker = false
			}
			return m, tea.Batch(cmds...)
		}

		if m.promptMode != "" {
			switch msg.Type {
			case tea.KeyEnter:
				switch m.promptMode {
				case "save":
					// Create directory if it doesn't exist
					dir := filepath.Dir(m.promptInput)
					if dir != "." && dir != "" {
						err := os.MkdirAll(dir, 0755)
						if err != nil {
							output.Write(fmt.Sprintf("Error creating directory: %s", err))
							m.output = output.String()
							m.promptMode = ""
							m.promptInput = ""
							return m, nil
						}
					}
					err := os.WriteFile(m.promptInput, []byte(m.textarea.Value()), 0644)
					if err != nil {
						output.Write(fmt.Sprintf("Error saving: %s", err))
					} else {
						output.Write(fmt.Sprintf("Saved to %s", m.promptInput))
					}
				}
				m.output = output.String()
				m.promptMode = ""
				m.promptInput = ""
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				m.promptMode = ""
				m.promptInput = ""
				return m, nil
			case tea.KeyBackspace:
				if len(m.promptInput) > 0 {
					m.promptInput = m.promptInput[:len(m.promptInput)-1]
				}
				return m, nil
			default:
				if msg.Type == tea.KeyRunes {
					m.promptInput += string(msg.Runes)
				}
				return m, nil
			}
		}

		// If table is showing, let it handle arrow keys
		if m.showTable {
			switch msg.Type {
			case tea.KeyUp, tea.KeyDown:
				m.dictionaryTable, cmd = m.dictionaryTable.Update(msg)
				return m, cmd
			}
		}

		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyTab:
			if m.textarea.Focused() {
				m.textarea.InsertString("    ")
				return m, nil
			}
		case tea.KeyCtrlN:
			m.textarea.SetValue("")
			m.currentFile = ""
			output.Write("New buffer created")
			m.output = output.String()
			return m, nil
		case tea.KeyCtrlS:
			if m.currentFile != "" {
				// Save to existing file
				dir := filepath.Dir(m.currentFile)
				if dir != "." && dir != "" {
					os.MkdirAll(dir, 0755)
				}
				err := os.WriteFile(m.currentFile, []byte(m.textarea.Value()), 0644)
				if err != nil {
					output.Write(fmt.Sprintf("Error saving: %s", err))
				} else {
					output.Write(fmt.Sprintf("Saved to %s", m.currentFile))
				}
				m.output = output.String()
			} else {
				// Prompt for new filename
				m.promptMode = "save"
				m.promptInput = ""
			}
			return m, nil
		case tea.KeyCtrlO:
			m.showFilePicker = true
			m.filePicker = makeFilePicker()
			return m, m.filePicker.Init()
		case tea.KeyF5:
			output.Clear()
			parseForthCode(m.textarea.Value())
			m.output = output.String()
			if m.showTable {
				m.dictionaryTable = makeWordsTable()
			}
			cmds = append(cmds, tea.ClearScreen)
		case tea.KeyF2:
			m.showTable = !m.showTable
			if m.showTable {
				m.dictionaryTable = makeWordsTable()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		if m.showFilePicker {
			m.filePicker, cmd = m.filePicker.Update(msg)
			return m, cmd
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Status bar
	statusLeft := statusBarStyle.Render("GoFth Editor - Forth Interpreter")

	var statusRightText string
	if m.promptMode != "" {
		statusRightText = fmt.Sprintf("%s: %s", strings.ToUpper(m.promptMode), m.promptInput+"_")
	} else if m.currentFile != "" {
		statusRightText = fmt.Sprintf("File: %s | F5: Run | F2: Dict | Ctrl+N: New | Ctrl+S: Save", filepath.Base(m.currentFile))
	} else {
		statusRightText = "Unsaved | F5: Run | F2: Dict | Ctrl+N: New | Ctrl+S: Save As"
	}

	statusRight := statusBarStyle.Render(statusRightText)
	gap := m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight)
	if gap < 0 {
		gap = 0
	}
	statusBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusLeft,
		strings.Repeat(" ", gap),
		statusRight,
	)

	if m.showFilePicker {
		return fmt.Sprintf("%s\n%s", statusBar, m.filePicker.View())
	}

	if m.showTable {
		return fmt.Sprintf("%s\n%s", statusBar, m.dictionaryTable.View())
	}

	outputSection := outputStyle.
		Width(m.width).
		Render(m.output)

	textareaHeight := max(m.height-lipgloss.Height(statusBar)-lipgloss.Height(outputSection)-2, 3)
	m.textarea.SetHeight(textareaHeight)

	return fmt.Sprintf(
		"%s\n%s\n\n--- Output ---\n%s",
		statusBar,
		m.textarea.View(),
		outputSection,
	)
}

type errMsg error
