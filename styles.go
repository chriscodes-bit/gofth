package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	FthWhite  = lipgloss.Color("#FFFFFF")
	FthOrange = lipgloss.Color("#FF5500")
	FthBlack  = lipgloss.Color("#000000")
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(FthOrange).
			Background(FthBlack).
			Padding(0, 1)

	outputStyle = lipgloss.NewStyle().
			Foreground(FthOrange).
			Background(FthBlack).
			Padding(0, 2).
			MarginTop(1)
)

func GetDictionaryTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(FthWhite).
		Background(FthOrange).
		Bold(true)
	return s
}
