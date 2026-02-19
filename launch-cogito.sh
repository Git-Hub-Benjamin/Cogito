#!/usr/bin/env bash
#
# launch-cogito.sh
#
# Intended to be bound to a KDE global keyboard shortcut.
# Behavior:
#   - If the focused window is a Konsole terminal, send "cogito\n" to the
#     currently active session via DBus.
#   - Otherwise, open a new Konsole running cogito.
#
# Dependencies: xdotool, xprop, qdbus (all standard on KDE/X11)
# For Wayland: replace xdotool/xprop with kdotool equivalents.

set -euo pipefail

# --- Detect the focused window's class and PID ---
ACTIVE_WIN=$(xdotool getactivewindow 2>/dev/null) || ACTIVE_WIN=""

if [[ -z "$ACTIVE_WIN" ]]; then
    # No active window detected; open a new Konsole
    exec konsole -e cogito
fi

# Get WM_CLASS (returns e.g. 'WM_CLASS(STRING) = "konsole", "konsole"')
WM_CLASS=$(xprop -id "$ACTIVE_WIN" WM_CLASS 2>/dev/null) || WM_CLASS=""

if [[ "$WM_CLASS" == *'"konsole"'* ]]; then
    # The focused window IS a Konsole instance.
    # Get its PID so we can find the correct DBus service.
    WIN_PID=$(xdotool getwindowpid "$ACTIVE_WIN" 2>/dev/null) || WIN_PID=""

    if [[ -z "$WIN_PID" ]]; then
        exec konsole -e cogito
    fi

    SERVICE="org.kde.konsole-${WIN_PID}"

    # Verify the DBus service exists
    if ! qdbus "$SERVICE" >/dev/null 2>&1; then
        # Fallback: try the generic service name (single-instance mode)
        SERVICE="org.kde.konsole"
        if ! qdbus "$SERVICE" >/dev/null 2>&1; then
            exec konsole -e cogito
        fi
    fi

    # Get the current (active) session ID from the first window
    SESSION_ID=$(qdbus "$SERVICE" /Windows/1 org.kde.konsole.Window.currentSession 2>/dev/null) || SESSION_ID=""

    if [[ -z "$SESSION_ID" ]]; then
        exec konsole -e cogito
    fi

    # Send the command to the active session
    qdbus "$SERVICE" "/Sessions/${SESSION_ID}" org.kde.konsole.Session.runCommand "cogito"
else
    # Focused window is not Konsole; open a new one
    exec konsole -e cogito
fi
