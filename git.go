package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s", msg)
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

func getRepoName() string {
	// Best-effort: show folder name. Works even if git commands fail.
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Base(wd)
}

func getBranches() ([]branchItem, string, error) {
	out, err := runGit("branch", "--no-color")
	if err != nil {
		return nil, "", err
	}

	lines := strings.Split(out, "\n")
	items := make([]branchItem, 0, len(lines))
	current := ""

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}

		isCurrent := strings.HasPrefix(line, "* ")
		name := strings.TrimSpace(strings.TrimPrefix(line, "* "))
		name = strings.TrimSpace(strings.TrimPrefix(name, "  "))

		if isCurrent {
			current = name
		}

		items = append(items, branchItem{name: name, isCurrent: isCurrent})
	}

	return items, current, nil
}

func checkoutBranch(name string) error {
	_, err := runGit("checkout", name)
	return err
}

func deleteBranch(name string) error {
	_, err := runGit("branch", "-d", name)
	return err
}
