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
		return checkoutDoneMsg{branch: branch}
	}
}

func deleteCmd(branch string) tea.Cmd {
	return func() tea.Msg {
		if err := deleteBranch(branch); err != nil {
			return errMsg(err.Error())
		}
		return deleteDoneMsg{branch: branch}
	}
}

func loadUpstreamCmd(branch string) tea.Cmd {
	return func() tea.Msg {
		up, err := getUpstream(branch)
		return upstreamMsg{branch: branch, upstream: up, err: err}
	}
}

func loadLogCmd(ref string) tea.Cmd {
	return func() tea.Msg {
		txt, err := getCommitLog(ref)
		return logMsg{ref: ref, text: txt, err: err}
	}
}

func initialModel() model {
	l := list.New([]list.Item{}, delegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return model{
		l:        l,
		status:   "Loading…",
		mode:     modeBrowse,
		repoName: getRepoName(),
	}
}

func (m model) Init() tea.Cmd { return loadBranchesCmd() }

func (m *model) setItems(branches []branchItem) {
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
		// header + divider + footer = 3 lines (+ confirm appended in body)
		m.l.SetSize(msg.Width, msg.Height-3)
		return m, nil

	case checkoutDoneMsg:
		// Successful checkout -> exit the app
		return m, tea.Quit

	case deleteDoneMsg:
		// Deletion finished; refresh the list after the branch is actually gone.
		m.errMsg = ""
		m.status = "Deleted \u201c" + msg.branch + "\u201d"
		return m, loadBranchesCmd()

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

		// Start upstream/log loading for current selection (if any)
		if b, ok := m.selectedBranch(); ok {
			m.selectedBranchName = b.name
			m.upstream = ""
			m.logText = ""
			m.logErr = ""
			return m, loadUpstreamCmd(b.name)
		}

		return m, nil

	case upstreamMsg:
		// Ignore stale results
		if msg.branch != m.selectedBranchName {
			return m, nil
		}

		// If no upstream, hide panel (common case). We don't treat it as an error.
		if msg.err != nil || strings.TrimSpace(msg.upstream) == "" {
			m.upstream = ""
			m.logText = ""
			m.logErr = ""
			return m, nil
		}

		m.upstream = msg.upstream
		m.logText = "Loading…"
		m.logErr = ""
		return m, loadLogCmd(m.upstream)

	case logMsg:
		// Ignore stale results
		if msg.ref != m.upstream {
			return m, nil
		}
		if msg.err != nil {
			m.logErr = msg.err.Error()
			m.logText = ""
			return m, nil
		}
		m.logErr = ""
		m.logText = msg.text
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
				if len(m.branches) == 0 {
					return m, nil
				}
				b, ok := m.selectedBranch()
				if !ok {
					return m, nil
				}
				// Checkout then reload
				return m, checkoutCmd(b.name)

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

			// Let the list handle navigation/filter keys
			prevSel := m.l.Index()

			var cmd tea.Cmd
			m.l, cmd = m.l.Update(msg)

			// If selection changed, update side panel state
			if m.l.Index() != prevSel {
				if b, ok := m.selectedBranch(); ok {
					m.selectedBranchName = b.name
					m.upstream = ""
					m.logText = ""
					m.logErr = ""
					return m, tea.Batch(cmd, loadUpstreamCmd(b.name))
				}
			}

			return m, cmd

		case modeConfirmDelete:
			switch msg.String() {
			case "backspace", "delete":
				branch := m.pendingDel
				m.mode = modeBrowse
				m.pendingDel = ""
				return m, deleteCmd(branch)

			case "esc":
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
	// Header: brand + repo on the left, counts on the right
	left := sBrand.Render("gbt") + sChip.Render(m.repoName)
	right := sChipAccent.Render(fmt.Sprintf("%d", len(m.branches))) + sChip.Render("branches")

	space := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if space < 1 {
		space = 1
	}

	headerLine := sHeader.Width(m.width).Render(left + strings.Repeat(" ", space) + right)
	divider := sDivider.Render(strings.Repeat("─", max(0, m.width)))

	// Footer
	var footer string
	if m.errMsg != "" {
		footer = sFooterErr.Width(m.width).Render("⚠ " + m.errMsg)
	} else {
		footer = sFooter.Width(m.width).Render("↑/↓ navigate   enter checkout   ⌫ delete   q quit   •   " + m.status)
	}

	confirm := ""
	if m.mode == modeConfirmDelete {
		confirm = sConfirm.Render(fmt.Sprintf("Press ⌫ again to delete “%s” (Esc to cancel)", m.pendingDel))
	}

	// Body: list + optional right panel
	body := m.renderBody()
	if confirm != "" {
		body = body + "\n" + confirm
	}

	return strings.Join([]string{
		headerLine,
		divider,
		body,
		footer,
	}, "\n")
}

func (m model) renderBody() string {
	// Always show a right panel (unless the terminal is extremely narrow).
	leftW := int(float64(m.width) * 0.70)
	if leftW < 50 {
		leftW = 50
	}
	rightW := m.width - leftW
	if rightW < 24 {
		// Too narrow to sensibly render a split view
		return m.l.View()
	}

	// Render list in left pane
	m2 := m
	m2.l.SetSize(leftW, m.l.Height())
	left := m2.l.View()

	// Right panel content
	header := sChipAccent.Render("Commits") + " " + sChip.Render(m.selectedBranchName)
	if strings.TrimSpace(m.upstream) != "" {
		header += " " + sChip.Render(m.upstream)
	} else {
		header += " " + sChip.Render("no upstream")
	}

	content := ""
	switch {
	case strings.TrimSpace(m.upstream) == "":
		// Empty state when branch isn't “remote” (no upstream)
		content = "Select a branch with an upstream to see commits."
	case m.logErr != "":
		content = "⚠ " + m.logErr
	case strings.TrimSpace(m.logText) == "":
		content = "No commits to show."
	default:
		// Truncate each commit line to fit the panel.
		maxLine := rightW - 4 // accounts for border + padding
		lines := strings.Split(m.logText, "\n")
		for i := range lines {
			lines[i] = truncateLine(lines[i], maxLine)
		}
		content = strings.Join(lines, "\n")
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cFaint).
		Padding(0, 1).
		Width(rightW).
		Height(m.l.Height()).
		Render(
			header + "\n" +
				sDivider.Render(strings.Repeat("─", max(0, rightW-2))) + "\n" +
				content,
		)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, panel)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func truncateLine(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	if maxLen == 1 {
		return "…"
	}
	return string(r[:maxLen-1]) + "…"
}
