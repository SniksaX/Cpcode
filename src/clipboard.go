package src

import (
    "os"
    "os/exec"
    "runtime"
    "strings"
)

func CopyToClipboard(text string) bool {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin":
        cmd = exec.Command("pbcopy")
    case "windows":
        cmd = exec.Command("clip")
    case "linux":
        if os.Getenv("WAYLAND_DISPLAY") != "" {
            cmd = exec.Command("wl-copy")
        } else if _, err := exec.LookPath("xclip"); err == nil {
            cmd = exec.Command("xclip", "-selection", "clipboard")
        } else if _, err := exec.LookPath("xsel"); err == nil {
            cmd = exec.Command("xsel", "--clipboard", "--input")
        } else {
            return false
        }
    default:
        return false
    }
    cmd.Stdin = strings.NewReader(text)
    return cmd.Run() == nil
}

func ReadFromClipboard() string {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin":
        cmd = exec.Command("pbpaste")
    case "windows":
        cmd = exec.Command("powershell", "-NoProfile", "-Command", "Get-Clipboard -Raw")
    case "linux":
        if os.Getenv("WAYLAND_DISPLAY") != "" {
            cmd = exec.Command("wl-paste")
        } else if _, err := exec.LookPath("xclip"); err == nil {
            cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
        } else if _, err := exec.LookPath("xsel"); err == nil {
            cmd = exec.Command("xsel", "--clipboard", "--output")
        } else {
            return ""
        }
    default:
        return ""
    }
    out, err := cmd.Output()
    if err != nil {
        return ""
    }
    return strings.TrimRight(string(out), "\r\n")
}