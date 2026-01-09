#!/usr/bin/env bash
set -euo pipefail

BIN_NAME="gbt"
PREFIX="${PREFIX:-$HOME/.local}"
BIN_DIR="$PREFIX/bin"

echo "🔨 Building $BIN_NAME..."
go build -o "$BIN_NAME"

echo "📦 Installing to $BIN_DIR/$BIN_NAME"
mkdir -p "$BIN_DIR"
install -m 755 "$BIN_NAME" "$BIN_DIR/$BIN_NAME"

# ---- PATH setup for zsh ----
ensure_zsh_path() {
  local export_line='export PATH="$HOME/.local/bin:$PATH"'
  local marker_start="# >>> gbt installer >>>"
  local marker_end="# <<< gbt installer <<<"

  # If already on PATH in this shell, skip
  case ":$PATH:" in
    *":$HOME/.local/bin:"*) return 0 ;;
  esac

  # Only touch zsh configs if zsh is in use or config exists
  if [[ "${SHELL:-}" != */zsh && ! -f "$HOME/.zshrc" ]]; then
    return 0
  fi

  add_block() {
    local file="$1"
    [[ -f "$file" ]] || touch "$file"

    if grep -Fq "$export_line" "$file"; then
      return 0
    fi

    {
      echo ""
      echo "$marker_start"
      echo "$export_line"
      echo "$marker_end"
    } >> "$file"
  }

  add_block "$HOME/.zshrc"
  add_block "$HOME/.zprofile"

  echo ""
  echo "🧠 Added ~/.local/bin to PATH in:"
  echo "   • ~/.zshrc"
  echo "   • ~/.zprofile"
  echo ""
  echo "Restart your terminal or run:"
  echo "  source ~/.zshrc"
}

ensure_zsh_path

echo ""
if command -v "$BIN_NAME" >/dev/null 2>&1; then
  echo "✅ Installed! Run: $BIN_NAME"
else
  echo "⚠️  Installed, but $BIN_NAME is not on PATH yet."
  echo "   Restart your terminal or run: source ~/.zshrc"
fi