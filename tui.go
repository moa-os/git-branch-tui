package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func loadBranchesCmd() tea.Cmd {
	return func() tea.Msg {
		branches, current, err := getBranches()
		return loadedMsg{branches: branches, current: current, err: err}
	}
}

func checkoutCmd(branch string) tea.Cmd {
	return func() tea.Msg {
		if err := checkoutBranch(branch); err != nil {
			return errMsg(err.Error())
		}
		return statusMsg("Checked out “" + branch + "”")
	}
}

func deleteCmd(branch string) tea.Cmd {
	return func() tea.Msg {
		if err := deleteBranch(branch); err != nil {
			return errMsg(err.Error())
		}
		return statusMsg("Deleted “" + branch + "”")
	}
}

func initialModel() model {
	l := list.New([]list.Item{}, delegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)

	return model{
		l:        l,
		status:   "Loading…",
		mode:     modeBrowse,
		repoName: getRepoName(),
	}
}

func (m model) Init() tea.Cmd { return loadBranchesCmd() }

func (m *model) setItems(branches []branchItem) {
	m.l.ResetFilter()

	items := make([]list.Item, 0, len(branches))
	for _, b := range branches {
		items = append(items, b)
	}
	m.l.SetItems(items)

	// select current branch by default
	if m.current != "" {
		for i, b := range branches {
			if b.name == m.current {
				m.l.Select(i)
				break
			}
		}
	}
}

func (m model) selectedBranch() (branchItem, bool) {
	it := m.l.SelectedItem()
	if it == nil {
		return branchItem{}, false
	}
	b, ok := it.(branchItem)
	return b, ok
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.l.SetSize(msg.Width, msg.Height-3) // header + footer + (maybe confirm)
		return m, nil

	case loadedMsg:
		if msg.err != nil {
			m.errMsg = msg.err.Error()
			m.status = "Not a git repo"
			return m, nil
		}
		m.errMsg = ""
		m.branches = msg.branches
		m.current = msg.current
		m.status = fmt.Sprintf("On %s", m.current)
		m.setItems(m.branches)
		return m, nil

	case errMsg:
		m.errMsg = string(msg)
		return m, nil

	case statusMsg:
		m.errMsg = ""
		m.status = string(msg)
		return m, nil

	case tea.KeyMsg:
		switch m.mode {

		case modeBrowse:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "enter":
				b, ok := m.selectedBranch()
				if !ok {
					return m, nil
				}
				return m, tea.Batch(checkoutCmd(b.name), loadBranchesCmd())

			case "backspace", "delete":
				b, ok := m.selectedBranch()
				if !ok {
					return m, nil
				}
				if b.isCurrent || b.name == m.current {
					m.errMsg = "You can’t delete the currently checked-out branch."
					return m, nil
				}
				m.mode = modeConfirmDelete
				m.pendingDel = b.name
				m.errMsg = ""
				return m, nil
			}

			var cmd tea.Cmd
			m.l, cmd = m.l.Update(msg)
			return m, cmd

		case modeConfirmDelete:
			switch msg.String() {
			case "y", "Y":
				branch := m.pendingDel
				m.mode = modeBrowse
				m.pendingDel = ""
				return m, tea.Batch(deleteCmd(branch), loadBranchesCmd())

			case "n", "N", "esc":
				m.mode = modeBrowse
				m.pendingDel = ""
				return m, nil

			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	left := sAppTitle.Render("gbt") + sPill.Render(m.repoName)
	right := sPillAlt.Render(fmt.Sprintf("%d branches", len(m.branches)))

	// header line
	header := sHeader.Width(m.width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left,
			lipgloss.PlaceHorizontal(m.width-lipgloss.Width(left), lipgloss.Right, right),
		),
	)

	// footer “toast”
	var footer string
	if m.errMsg != "" {
		footer = sToastErr.Width(m.width).Render("⚠ " + m.errMsg)
	} else {
		footer = sToast.Width(m.width).Render("↑/↓ navigate   enter checkout   ⌫ delete   / filter   q quit   •   " + m.status)
	}

	confirm := ""
	if m.mode == modeConfirmDelete {
		confirm = sConfirm.Render(fmt.Sprintf("Delete “%s”? (y/n)", m.pendingDel))
	}

	body := m.l.View()
	if confirm != "" {
		body = body + "\n" + confirm
	}

	return strings.Join([]string{
		header,
		body,
		footer,
	}, "\n")
}
