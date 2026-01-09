package main

import (
	"github.com/charmbracelet/bubbles/list"
)

type branchItem struct {
	name      string
	isCurrent bool
}

func (b branchItem) Title() string       { return b.name }
func (b branchItem) Description() string { return "" }
func (b branchItem) FilterValue() string { return b.name }

type mode int

const (
	modeBrowse mode = iota
	modeConfirmDelete
)

type model struct {
	l             list.Model
	branches      []branchItem
	current       string
	status        string
	errMsg        string
	mode          mode
	pendingDel    string
	width, height int
	repoName      string

	// Right panel state (shows only when selected branch has upstream)
	selectedBranchName string
	upstream           string
	logText            string
	logErr             string
}

type loadedMsg struct {
	branches []branchItem
	current  string
	err      error
}

type errMsg string
type statusMsg string

type upstreamMsg struct {
	branch   string
	upstream string
	err      error
}

type logMsg struct {
	ref  string
	text string
	err  error
}

type checkoutDoneMsg struct {
	branch string
}

type deleteDoneMsg struct {
	branch string
}
