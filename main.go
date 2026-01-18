package main

import (
	"log"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textarea        textarea.Model
	dictionaryTable table.Model
	filePicker      filepicker.Model
	showTable       bool
	showFilePicker  bool
	output          string
	width           int
	height          int
	err             error
	promptMode      string
	promptInput     string
	currentFile     string
}
