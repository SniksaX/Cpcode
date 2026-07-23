package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "strings"
    "time"

    "cpcode/src"
)

func main() {
    start := time.Now()

    cfg, rawPatterns, err := src.ParseFlags()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
        os.Exit(1)
    }

    stat, _ := os.Stdin.Stat()
    isPiped := (stat.Mode() & os.ModeCharDevice) == 0
    var pipedBytes []byte
    if isPiped {
        pipedBytes, _ = io.ReadAll(os.Stdin)
        if len(bytes.TrimSpace(pipedBytes)) == 0 {
            pipedBytes = nil
        }
    }

    if len(rawPatterns) == 0 && pipedBytes == nil {
        rawPatterns = []string{"."}
    }

    files, _, err := src.CollectFiles(cfg, rawPatterns)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error collecting files: %v\n", err)
        os.Exit(1)
    }

    cwd, _ := os.Getwd()
    treeStr := ""
    if len(files) > 0 {
        treeStr = src.GenerateTree(files, cwd)
    }
    if pipedBytes != nil {
        if treeStr != "" {
            treeStr += "\n"
        }
        treeStr += "└── [Terminal_Output (stdin)]"
    }

    if cfg.DryRun {
        duration := time.Since(start)
        src.PrintCLISummary(cfg, treeStr, int64(len(files)), 0, 0, 0, duration, false, 0)
        return
    }

    var results []src.FileResult
    var validCount, skippedCount int64
    if len(files) > 0 {
        results, validCount, skippedCount = src.ProcessFiles(files)
    }

    if pipedBytes != nil {
        results = append(results, src.FileResult{
            Path:    "Terminal_Output (stdin)",
            Content: pipedBytes,
        })
        validCount++
    }

    var payloadBuilder bytes.Buffer
    var existingTextSize int64

    if cfg.Append {
        existing := ""
        if cfg.Dest != "" {
            if b, err := os.ReadFile(cfg.Dest); err == nil {
                existing = string(b)
            }
        } else {
            existing = src.ReadFromClipboard()
        }
        if idx := strings.Index(existing, src.FooterMarker); idx != -1 {
            existing = strings.TrimSpace(existing[:idx]) + "\n\n"
        } else if strings.TrimSpace(existing) != "" {
            existing = strings.TrimSpace(existing) + "\n\n"
        }
        existingTextSize = int64(len(existing))
        payloadBuilder.WriteString(existing)
    }

    if cfg.Format == "xml" {
        payloadBuilder.WriteString("<?xml version=\"1.0\"?>\n<files>\n")
    }

    for _, res := range results {
        payloadBuilder.WriteString(src.FormatContent(res, cfg.Format))
    }

    if cfg.Format == "xml" {
        payloadBuilder.WriteString("</files>\n")
    }

    totalBytes := existingTextSize + int64(payloadBuilder.Len())
    newBytesSize := totalBytes - existingTextSize
    duration := time.Since(start)
    tokenEstimate := 0

    footer := src.BuildFooter(cfg, treeStr, validCount, skippedCount,
        totalBytes, newBytesSize, duration)
    payloadBuilder.WriteString(footer)

    fullPayload := payloadBuilder.String()

    if cfg.Tokens {
        tokenEstimate = len(fullPayload) / 4
    }

    if cfg.Dest != "" {
        err := os.WriteFile(cfg.Dest, []byte(fullPayload), 0644)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error writing to %s: %v\n", cfg.Dest, err)
            os.Exit(1)
        }
        action := "Written"
        if cfg.Append {
            action = "Appended"
        }
        fmt.Fprintf(os.Stderr, "%s to %s\n", action, cfg.Dest)
        src.PrintCLISummary(cfg, treeStr, validCount, skippedCount,
            totalBytes, newBytesSize, duration, cfg.Tokens, tokenEstimate)
    } else if cfg.Print {
        fmt.Println(fullPayload)
        fmt.Fprintf(os.Stderr, "\n Output printed to stdout.\n")
        src.PrintCLISummary(cfg, treeStr, validCount, skippedCount,
            totalBytes, newBytesSize, duration, cfg.Tokens, tokenEstimate)
    } else {
        if src.CopyToClipboard(fullPayload) {
            action := "Copied"
            if cfg.Append {
                action = "Appended"
            }
            fmt.Fprintf(os.Stderr, "%s to clipboard successfully!\n", action)
            src.PrintCLISummary(cfg, treeStr, validCount, skippedCount,
                totalBytes, newBytesSize, duration, cfg.Tokens, tokenEstimate)
        } else {
            fmt.Println(fullPayload)
            fmt.Fprintln(os.Stderr, "\n Failed to copy to clipboard. Printed to stdout instead.")
        }
    }
}