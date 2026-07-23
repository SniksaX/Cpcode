package src

import (
    "flag"
    "os"
    "strings"
)


type Config struct {
    OmitPatterns    []string
    IncludePatterns []string
    Dest            string
    Print           bool
    NoIgnore        bool
    Append          bool
    Tokens          bool
    DryRun          bool
    Format          string
    UseGitignore    bool
    Quiet           bool
}

func ParseFlags() (Config, []string, error) {
    args := os.Args[1:]

    omitPatterns, args := extractPatterns(args, "o", "omit")

    includePatterns, args := extractPatterns(args, "i", "include")

    args = reorderArgs(args)

    var cfg Config
    fs := flag.NewFlagSet("cpcode", flag.ContinueOnError)
    fs.StringVar(&cfg.Dest, "d", "", "Save output to file")
    fs.StringVar(&cfg.Dest, "destination", "", "Save output to file")
    fs.BoolVar(&cfg.Print, "p", false, "Print to stdout instead of copying")
    fs.BoolVar(&cfg.Print, "print", false, "Print to stdout instead of copying")
    fs.BoolVar(&cfg.NoIgnore, "no-ignore", false, "Disable default ignore lists")
    fs.BoolVar(&cfg.Append, "a", false, "Append to existing clipboard/file content")
    fs.BoolVar(&cfg.Append, "append", false, "Append to existing clipboard/file content")
    fs.BoolVar(&cfg.Tokens, "t", false, "Estimate token count")
    fs.BoolVar(&cfg.Tokens, "tokens", false, "Estimate token count")
    fs.BoolVar(&cfg.DryRun, "dry-run", false, "Dry run (show summary without copying)")
    fs.StringVar(&cfg.Format, "format", "markdown", "Output format: markdown, xml, plain")
    fs.BoolVar(&cfg.UseGitignore, "use-gitignore", false, "Respect .gitignore patterns from root directory")
    fs.BoolVar(&cfg.Quiet, "q", false, "Suppress CLI summary")
    fs.BoolVar(&cfg.Quiet, "quiet", false, "Suppress CLI summary")

    if err := fs.Parse(args); err != nil {
        if err == flag.ErrHelp {
            os.Exit(0)
        }
        return cfg, nil, err
    }

    cfg.OmitPatterns = cleanPatterns(omitPatterns)
    cfg.IncludePatterns = cleanPatterns(includePatterns)

    cfg.Format = strings.ToLower(strings.TrimSpace(cfg.Format))
    if cfg.Format != "markdown" && cfg.Format != "xml" && cfg.Format != "plain" {
        cfg.Format = "markdown"
    }

    return cfg, fs.Args(), nil
}

func reorderArgs(args []string) []string {
    var flags []string
    var positional []string

    i := 0
    for i < len(args) {
        arg := args[i]
        if strings.HasPrefix(arg, "-") && arg != "-" {
            flags = append(flags, arg)
            if !strings.Contains(arg, "=") {
                switch arg {
                case "-d", "-destination", "-format":
                    if i+1 < len(args) {
                        flags = append(flags, args[i+1])
                        i++
                    }
                }
            }
        } else {
            positional = append(positional, arg)
        }
        i++
    }
    return append(flags, positional...)
}

func extractPatterns(args []string, shortFlag, longFlag string) (patterns, remaining []string) {
    remaining = args
    i := 0
    for i < len(remaining) {
        arg := remaining[i]
        if arg == "-"+shortFlag || arg == "-"+longFlag {
            if i+1 >= len(remaining) {
                i++
                continue
            }
            var collected []string
            j := i + 1
            for j < len(remaining) {
                p := strings.TrimSpace(remaining[j])
                if strings.HasSuffix(p, ",") {
                    collected = append(collected, strings.TrimSuffix(p, ","))
                    j++
                } else {
                    collected = append(collected, p)
                    j++
                    break
                }
            }
            patterns = append(patterns, collected...)
            remaining = append(remaining[:i], remaining[j:]...)
        } else if strings.HasPrefix(arg, "-"+shortFlag+"=") || strings.HasPrefix(arg, "-"+longFlag+"=") {
            eqIdx := strings.Index(arg, "=")
            firstPart := strings.TrimSpace(arg[eqIdx+1:])
            var collected []string
            if strings.HasSuffix(firstPart, ",") {
                collected = append(collected, strings.TrimSuffix(firstPart, ","))
                j := i + 1
                for j < len(remaining) {
                    p := strings.TrimSpace(remaining[j])
                    if strings.HasSuffix(p, ",") {
                        collected = append(collected, strings.TrimSuffix(p, ","))
                        j++
                    } else {
                        collected = append(collected, p)
                        j++
                        break
                    }
                }
                patterns = append(patterns, collected...)
                remaining = append(remaining[:i], remaining[j:]...)
            } else {
                collected = append(collected, firstPart)
                patterns = append(patterns, collected...)
                remaining = append(remaining[:i], remaining[i+1:]...)
            }
        } else {
            i++
        }
    }
    return
}

func cleanPatterns(raw []string) []string {
    var cleaned []string
    for _, p := range raw {
        p = strings.TrimSpace(p)
        if p != "" {
            cleaned = append(cleaned, p)
        }
    }
    return cleaned
}