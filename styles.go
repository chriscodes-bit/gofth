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
	FthGrey   = lipgloss.Color("#343434")
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(FthWhite).
			Padding(0, 1).
			Bold(true).
			BorderStyle(lipgloss.ThickBorder())

	outputStyle = lipgloss.NewStyle().
			Foreground(FthWhite).
			Background(FthGrey).
			Padding(1, 1).
			BorderStyle(lipgloss.ThickBorder())
)

func GetDictionaryTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(FthOrange).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(FthWhite).
		Background(FthOrange).
		Bold(true)
	return s
}
