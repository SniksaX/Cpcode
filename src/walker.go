package src

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

func CollectFiles(cfg Config, rawPatterns []string) (files []string, rootDir string, err error) {
    if len(rawPatterns) == 0 {
        rawPatterns = []string{"."}
    }

    resolved := resolvePatterns(rawPatterns)
    cwd, _ := os.Getwd()
    rootDir = cwd

    for _, p := range resolved {
        abs, _ := filepath.Abs(p)
        info, statErr := os.Stat(abs)
        if statErr == nil && info.IsDir() {
            rootDir = abs
            break
        }
    }

    ignoreList := NewIgnoreList(cfg, rootDir)
    fileSet := make(map[string]bool)

    for _, p := range resolved {
        matches, globErr := filepath.Glob(p)
        if globErr != nil || len(matches) == 0 {
            matches = []string{p}
        }
        for _, match := range matches {
            absMatch, _ := filepath.Abs(match)
            info, statErr := os.Stat(absMatch)
            if statErr != nil {
                continue
            }
            if !info.IsDir() {
                name := filepath.Base(absMatch)
                rel, _ := filepath.Rel(rootDir, absMatch)
                if ignoreList.ShouldSkip(name, rel, false) {
                    continue
                }
                fileSet[absMatch] = true
                continue
            }
            filepath.WalkDir(absMatch, func(path string, d os.DirEntry, walkErr error) error {
                if walkErr != nil {
                    fmt.Fprintf(os.Stderr, "⚠️ Walk error at %s: %v\n", path, walkErr)
                    return nil
                }
                name := d.Name()
                rel, relErr := filepath.Rel(rootDir, path)
                if relErr != nil {
                    rel = path
                }
                if d.IsDir() {
                    if ignoreList.ShouldSkip(name, rel, true) {
                        return filepath.SkipDir
                    }
                    return nil
                }
                if ignoreList.ShouldSkip(name, rel, false) {
                    return nil
                }
                fileSet[path] = true
                return nil
            })
        }
    }

    for f := range fileSet {
        files = append(files, f)
    }
    sort.Strings(files)
    return files, rootDir, nil
}

func resolvePatterns(rawPatterns []string) []string {
    var resolved []string
    var currentBase string
    for _, p := range rawPatterns {
        absCwdPath, _ := filepath.Abs(p)
        _, cwdErr := os.Stat(absCwdPath)
        isGlob := strings.ContainsAny(p, "*?[")

        if cwdErr != nil && !isGlob && currentBase != "" {
            absBasePath := filepath.Join(currentBase, p)
            if _, baseErr := os.Stat(absBasePath); baseErr == nil {
                p = absBasePath
            }
        }

        absMatch, _ := filepath.Abs(p)
        info, err := os.Stat(absMatch)
        if err == nil {
            if info.IsDir() {
                currentBase = absMatch
            } else {
                currentBase = filepath.Dir(absMatch)
            }
        } else {
            currentBase = filepath.Dir(absMatch)
        }
        resolved = append(resolved, p)
    }
    return resolved
}