# gbt ⚡️

A fast, keyboard-first Git branch picker with a modern TUI.

gbt lets you browse local branches, preview upstream commits, switch branches, and delete branches — all without leaving the terminal.

Built with Go and Bubble Tea.

---

## ✨ Features

- 🚀 Instant branch switching
- 🧹 Safe branch deletion (double Backspace to confirm)
- 👀 Upstream commit preview (last 5 commits)
- ⌨️ Fully keyboard-driven
- 🎨 Modern neon TUI
- ⚡ No git fetch, no remote clutter

Only local branches are listed.
If a branch has an upstream, you’ll see its recent commits in the side panel.

---

## 📦 Installation

### Option 1: Go install (recommended)

Run this command:

```bash
go install https://github.com/moa-os/git-branch-tui@v0.1.0
```

Make sure ~/go/bin is on your PATH.

---

### Option 2: Install script (macOS / Linux)

Run:

git clone https://github.com/moa-os/git-branch-tui
cd git-branch-tui
./scripts/install.sh

The install script:
- builds gbt
- installs it to ~/.local/bin
- adds that directory to your PATH for zsh (safely)

---

## 🚀 Usage

Run gbt from inside any git repository:

gbt

### Key bindings

- Up / Down — navigate branches
- Enter — checkout branch and exit
- Backspace — delete branch (press twice to confirm)
- Esc — cancel delete
- q or Ctrl+C — quit

---

## 📄 License

MIT
