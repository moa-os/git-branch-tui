package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// palette (works well on dark terminals)
var (
	cBg      = lipgloss.Color("#0b1020") // deep navy
	cPanel   = lipgloss.Color("#111a33")
	cText    = lipgloss.Color("#e5e7eb")
	cMuted   = lipgloss.Color("#93a4c7")
	cAccent  = lipgloss.Color("#7dd3fc") // sky
	cAccent2 = lipgloss.Color("#a78bfa") // violet
	cDanger  = lipgloss.Color("#fb7185") // rose
	cSuccess = lipgloss.Color("#34d399") // emerald
)

var (
	// header
	sAppTitle = lipgloss.NewStyle().Bold(true).Foreground(cText)
	sPill     = lipgloss.NewStyle().
			Foreground(cBg).
			Background(cAccent).
			Bold(true).
			Padding(0, 1).
			MarginLeft(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cAccent)

	sPillAlt = lipgloss.NewStyle().
			Foreground(cBg).
			Background(cAccent2).
			Bold(true).
			Padding(0, 1).
			MarginLeft(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cAccent2)

	sHeader = lipgloss.NewStyle().
		Background(cPanel).
		Foreground(cMuted).
		Padding(0, 1)

	// list rows
	sRow = lipgloss.NewStyle().
		Foreground(cText).
		Padding(0, 1)

	sRowDim = lipgloss.NewStyle().
		Foreground(cMuted).
		Padding(0, 1)

	sSelected = lipgloss.NewStyle().
			Foreground(cBg).
			Background(cAccent).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cAccent)

	sCurrentBadge = lipgloss.NewStyle().
			Foreground(cSuccess).
			Bold(true)

	// footer “toast”
	sToast = lipgloss.NewStyle().
		Background(cPanel).
		Foreground(cMuted).
		Padding(0, 1)

	sToastErr = lipgloss.NewStyle().
			Background(cPanel).
			Foreground(cDanger).
			Bold(true).
			Padding(0, 1)

	sConfirm = lipgloss.NewStyle().
			Foreground(cDanger).
			Bold(true).
			Padding(0, 1)
)

// Custom delegate so we can render badges and nicer selection
type delegate struct{}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	b, ok := item.(branchItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	prefix := "  "
	if b.isCurrent {
		prefix = sCurrentBadge.Render("●") + " "
	}

	line := prefix + b.name

	if isSelected {
		fmt.Fprint(w, sSelected.Render(line))
		return
	}

	if !b.isCurrent {
		fmt.Fprint(w, sRowDim.Render(line))
		return
	}

	fmt.Fprint(w, sRow.Render(line))
}
