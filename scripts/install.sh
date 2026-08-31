#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
install_root="${AUTO_COMPANY_HOME:-$HOME/.auto-company}"
bin_dir="$install_root/bin"
install_codex_skills="${INSTALL_CODEX_SKILLS:-0}"

cd "$repo_root"
go test ./...
node scripts/validate-plugin.mjs
node scripts/test-guard.mjs

mkdir -p "$bin_dir"
go build -trimpath -o "$bin_dir/autoco" ./cmd/autoco
chmod +x "$bin_dir/autoco"

echo "Installed CLI: $bin_dir/autoco"
echo "PATH was not modified. Add $bin_dir to PATH when ready."

if [[ "$install_codex_skills" == "1" ]]; then
  mkdir -p "$HOME/.agents/skills"
  for skill in "$repo_root"/skills/*; do
    [[ -d "$skill" ]] || continue
    name="$(basename "$skill")"
    destination="$HOME/.agents/skills/auto-company-$name"
    rm -rf "$destination"
    cp -R "$skill" "$destination"
  done
  echo "Installed Auto Company skills for Codex under $HOME/.agents/skills"
fi

echo
echo "Claude Code local plugin:"
echo "  claude --plugin-dir \"$repo_root\""
