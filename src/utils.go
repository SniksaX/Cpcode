package src

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const FooterSeparator = "=================================================="
const FooterMarker = "<!-- cpcode:footer -->"

var (
    homeDir    string
    homePrefix string
)

func init() {
    d, err := os.UserHomeDir()
    if err == nil {
        homeDir = d
        homePrefix = d + string(filepath.Separator)
    }
}

func prettyPath(path string) string {
    if homePrefix != "" && strings.HasPrefix(path, homePrefix) {
        return "~" + strings.TrimPrefix(path, homeDir)
    }
    return path
}

func FormatBytes(b int64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func BuildFooter(cfg Config, treeStr string, validCount, skippedCount,
    totalBytes, newBytesSize int64, duration time.Duration) string {

    var sb strings.Builder
    title := "Project Structure"
    if cfg.Append {
        title = "Appended Files Structure"
    }

    switch cfg.Format {
    case "xml":
        sb.WriteString(fmt.Sprintf("\n<!-- %s -->\n", FooterMarker))
        sb.WriteString(fmt.Sprintf("<!-- %s -->\n", FooterSeparator))
        sb.WriteString(fmt.Sprintf("<!-- %s:\n%s\n-->\n", title, treeStr))
        sb.WriteString("<!-- --- Stats ---\n")
        sb.WriteString(fmt.Sprintf("Files/Streams Added : %d\n", validCount))
        if cfg.Append {
            sb.WriteString(fmt.Sprintf("Added Size          : %s\n", FormatBytes(newBytesSize)))
        }
        sb.WriteString(fmt.Sprintf("Total Size          : %s\n", FormatBytes(totalBytes)))
        sb.WriteString(fmt.Sprintf("Time Taken          : %v -->\n", duration))
    case "plain":
        sb.WriteString(fmt.Sprintf("\n%s\n%s\n%s:\n%s\n\n--- Stats ---\n",
            FooterMarker, FooterSeparator, title, treeStr))
        sb.WriteString(fmt.Sprintf("Files/Streams Added : %d\n", validCount))
        if cfg.Append {
            sb.WriteString(fmt.Sprintf("Added Size          : %s\n", FormatBytes(newBytesSize)))
        }
        sb.WriteString(fmt.Sprintf("Total Size          : %s\n", FormatBytes(totalBytes)))
        sb.WriteString(fmt.Sprintf("Time Taken          : %v\n", duration))
    default:
        sb.WriteString(fmt.Sprintf("\n%s\n%s\n%s:\n%s\n\n--- Stats ---\n", FooterMarker, FooterSeparator, title, treeStr))
        sb.WriteString(fmt.Sprintf("Files/Streams Added : %d\n", validCount))
        if cfg.Append {
            sb.WriteString(fmt.Sprintf("Added Size          : %s\n", FormatBytes(newBytesSize)))
        }
        sb.WriteString(fmt.Sprintf("Total Size          : %s\n", FormatBytes(totalBytes)))
        sb.WriteString(fmt.Sprintf("Time Taken          : %v\n", duration))
    }
    return sb.String()
}

func PrintCLISummary(cfg Config, treeStr string, validCount, skippedCount,
    totalBytes, newBytesSize int64, duration time.Duration, showTokens bool, tokenEstimate int) {
    if cfg.Quiet {
        return
    }
    var cli strings.Builder
    title := "Project Structure"
    if cfg.Append {
        title = "Appended Files Structure"
    }
    cli.WriteString(fmt.Sprintf("\n🌳 %s:\n%s\n\n", title, treeStr))
    cli.WriteString("Execution Stats\n")
    cli.WriteString(fmt.Sprintf("├─ Streams Added : %d\n", validCount))
    if skippedCount > 0 {
        cli.WriteString(fmt.Sprintf("├─ Skipped       : %d (Binary/Large/Empty)\n", skippedCount))
    }
    if cfg.DryRun {
        cli.WriteString("├─ Mode          : DRY RUN\n")
    } else if cfg.Append {
        cli.WriteString(fmt.Sprintf("├─ Appended Size : %s\n", FormatBytes(newBytesSize)))
        cli.WriteString(fmt.Sprintf("├─ Total Payload : %s\n", FormatBytes(totalBytes)))
    } else {
        cli.WriteString(fmt.Sprintf("├─ Payload Size  : %s\n", FormatBytes(totalBytes)))
    }
    if showTokens {
        cli.WriteString(fmt.Sprintf("├─ Est. Tokens   : %d\n", tokenEstimate))
    }
    cli.WriteString(fmt.Sprintf("└─ Time Taken    : %v\n", duration))
    fmt.Fprint(os.Stderr, cli.String())
}