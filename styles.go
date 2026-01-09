package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// -------- Palette (modern, calm, one accent) --------

var (
	cSurface = lipgloss.Color("#0B1020") // deep midnight
	cText    = lipgloss.Color("#F8FAFF") // bright white
	cMuted   = lipgloss.Color("#8B93B5") // cool grey-violet
	cFaint   = lipgloss.Color("#2A2F55") // neon-ish divider
	cAccent  = lipgloss.Color("#00F5FF") // neon cyan
	cDanger  = lipgloss.Color("#FF2FD6") // neon magenta
)

// -------- Layout styles --------

var (
	sHeader = lipgloss.NewStyle().
		Background(cSurface).
		Foreground(cMuted).
		Padding(0, 1)

	sDivider = lipgloss.NewStyle().
			Foreground(cFaint)

	sBrand = lipgloss.NewStyle().
		Bold(true).
		Foreground(cText)

	sChip = lipgloss.NewStyle().
		Foreground(cText).
		Background(lipgloss.Color("#182443")).
		Padding(0, 1).
		MarginLeft(1)

	sChipAccent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06121A")).
			Background(cAccent).
			Bold(true).
			Padding(0, 1).
			MarginLeft(1)

	sFooter = lipgloss.NewStyle().
		Background(cSurface).
		Foreground(cMuted).
		Padding(0, 1)

	sFooterErr = lipgloss.NewStyle().
			Background(cSurface).
			Foreground(cDanger).
			Bold(true).
			Padding(0, 1)

	sConfirm = lipgloss.NewStyle().
			Foreground(cDanger).
			Bold(true).
			Padding(0, 1)
)

// -------- Row rendering (modern selection) --------

var (
	sRow = lipgloss.NewStyle().
		Foreground(cText).
		Padding(0, 1)

	sRowDim = lipgloss.NewStyle().
		Foreground(cMuted).
		Padding(0, 1)

	// selection = left bar + soft highlight (no rounded border)
	sSelected = lipgloss.NewStyle().
			Foreground(cText).
			Background(lipgloss.Color("#152849")). // soft highlight
			Padding(0, 1)

	sSelBar = lipgloss.NewStyle().
		Foreground(cAccent)

	sRowActive = lipgloss.NewStyle().
			Foreground(cAccent).
			Padding(0, 1)

	sCurrentArrow = lipgloss.NewStyle().
			Foreground(cAccent).
			Bold(true)
)

// Custom delegate to render selection + current branch arrow
type delegate struct{}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	b, ok := item.(branchItem)
	if !ok {
		return
	}

	selected := index == m.Index()

	// Build a plain-text line so styles can apply uniformly (incl. background on hover)
	line := "  " + b.name
	if b.isCurrent {
		line = "▶ " + b.name
	}

	// 1) Hover / selected row ALWAYS gets background
	if selected {
		bar := sSelBar.Render("│")

		// Current + selected: keep the hover background but use accent text color
		if b.isCurrent {
			fmt.Fprint(w, bar+sSelected.Foreground(cAccent).Render(line))
			return
		}

		fmt.Fprint(w, bar+sSelected.Render(line))
		return
	}

	// 2) Current branch (not selected): accent text
	if b.isCurrent {
		fmt.Fprint(w, " "+sRowActive.Render(line))
		return
	}

	// 3) Normal row
	fmt.Fprint(w, " "+sRowDim.Render(line))
}
