#!/usr/bin/env bash
set -euo pipefail

# Render the matching scheduler template with the real repo path and activate it.
REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PLACEHOLDER="/ABSOLUTE/PATH/TO/llama-watch"

install_launchd() {
  local src="${REPO_DIR}/scripts/launchd/com.gultekinmakif.llama-watch.refresh.plist"
  local target="${HOME}/Library/LaunchAgents/com.gultekinmakif.llama-watch.refresh.plist"
  mkdir -p "${HOME}/Library/LaunchAgents"
  # Unload any prior copy before re-rendering; ignore "not loaded" on first run.
  if [ -f "${target}" ]; then
    launchctl unload "${target}" 2>/dev/null || true
  fi
  sed "s|${PLACEHOLDER}|${REPO_DIR}|g" "${src}" > "${target}"
  launchctl load -w "${target}"
  echo "launchd: installed at ${target}"
}

install_systemd() {
  if [ "$(id -u)" != "0" ]; then
    echo "systemd install requires sudo; re-run as: sudo \"$0\" $*"
    exit 1
  fi
  local svc_src="${REPO_DIR}/scripts/systemd/llama-watch-refresh.service"
  local tmr_src="${REPO_DIR}/scripts/systemd/llama-watch-refresh.timer"
  local svc_target="/etc/systemd/system/llama-watch-refresh.service"
  local tmr_target="/etc/systemd/system/llama-watch-refresh.timer"
  sed "s|${PLACEHOLDER}|${REPO_DIR}|g" "${svc_src}" > "${svc_target}"
  sed "s|${PLACEHOLDER}|${REPO_DIR}|g" "${tmr_src}" > "${tmr_target}"
  systemctl daemon-reload
  systemctl enable --now llama-watch-refresh.timer
  echo "systemd: installed ${svc_target} and ${tmr_target}, timer enabled"
}

install_crontab() {
  local line="0 * * * * cd ${REPO_DIR} && ./scripts/refresh.sh >> var/log/refresh.log 2>&1"
  # Idempotent: match the full assembled line so a stray comment or stale entry from
  # another checkout cannot trick the gate into a silent no-op.
  if crontab -l 2>/dev/null | grep -Fqx "${line}"; then
    echo "crontab: already installed for ${REPO_DIR}"
    return 0
  fi
  (crontab -l 2>/dev/null; echo "${line}") | crontab -
  echo "crontab: installed hourly entry for ${REPO_DIR}"
}

case "$(uname -s)" in
  Darwin) install_launchd ;;
  Linux)  install_systemd ;;
  *)      install_crontab ;;
esac
