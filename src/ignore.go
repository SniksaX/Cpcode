package src

import (
    "os"
    "path/filepath"
    "strings"
)

var defaultIgnoredDirs = map[string]struct{}{
    ".git": {}, ".svn": {}, ".hg": {}, ".idea": {}, ".vscode": {}, ".settings": {},
    "node_modules": {}, "bower_components": {}, "jspm_packages": {},
    "dist": {}, "build": {}, "out": {}, "target": {}, "bin": {}, "obj": {},
    "__pycache__": {}, "venv": {}, ".venv": {}, "env": {}, ".env": {}, "eggs": {}, ".eggs": {},
    ".mypy_cache": {}, ".pytest_cache": {}, ".tox": {}, ".ruff_cache": {},
    "vendor": {}, "bundle": {}, "gems": {}, "_build": {},
    "cov": {}, "coverage": {}, ".nyc_output": {},
    "tmp": {}, "temp": {}, "logs": {}, "log": {}, "private": {},
    "site-packages": {}, "static": {}, "media": {}, "assets": {}, "qdrant_storage": {},
}

var defaultIgnoredFiles = map[string]struct{}{
    ".DS_Store": {}, "Thumbs.db": {}, "desktop.ini": {},
    "package-lock.json": {}, "yarn.lock": {}, "pnpm-lock.yaml": {}, "bun.lockb": {},
    "Gemfile.lock": {}, "composer.lock": {}, "poetry.lock": {}, "Pipfile.lock": {}, "Cargo.lock": {},
    "LICENSE": {}, "LICENSE.txt": {}, "LICENSE.md": {}, "CHANGELOG.md": {},
    ".gitignore": {}, ".dockerignore": {}, ".npmignore": {}, ".env": {},
}

var defaultIgnoredExts = map[string]struct{}{
    ".exe": {}, ".dll": {}, ".so": {}, ".dylib": {}, ".class": {}, ".jar": {}, ".war": {}, ".ear": {},
    ".zip": {}, ".tar": {}, ".gz": {}, ".rar": {}, ".7z": {}, ".iso": {}, ".img": {}, ".dmg": {},
    ".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".bmp": {}, ".ico": {}, ".svg": {}, ".webp": {}, ".tiff": {},
    ".mp3": {}, ".wav": {}, ".mp4": {}, ".avi": {}, ".mov": {}, ".mkv": {}, ".flv": {}, ".webm": {},
    ".pdf": {}, ".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {}, ".ppt": {}, ".pptx": {},
    ".db": {}, ".sqlite": {}, ".sqlite3": {}, ".parquet": {}, ".hdf5": {},
    ".pyc": {}, ".pyo": {}, ".pyd": {}, ".o": {}, ".a": {}, ".lib": {},
    ".eot": {}, ".otf": {}, ".ttf": {}, ".woff": {}, ".woff2": {},
}

type IgnoreList struct {
    ignoreDirs       map[string]struct{}
    ignoreFiles      map[string]struct{}
    ignoreExts       map[string]struct{}
    omitPatterns     []string
    includePatterns  []string
    gitignorePatterns []string
    useGitignore     bool
    rootDir          string
}

func NewIgnoreList(cfg Config, gitignoreRoot string) *IgnoreList {
    il := &IgnoreList{
        ignoreDirs:      copyMap(defaultIgnoredDirs),
        ignoreFiles:     copyMap(defaultIgnoredFiles),
        ignoreExts:      copyMap(defaultIgnoredExts),
        omitPatterns:    cfg.OmitPatterns,
        includePatterns: cfg.IncludePatterns,
        useGitignore:    cfg.UseGitignore,
        rootDir:         gitignoreRoot,
    }
    if cfg.NoIgnore {
        il.ignoreDirs = make(map[string]struct{})
        il.ignoreFiles = make(map[string]struct{})
        il.ignoreExts = make(map[string]struct{})
    }
    if cfg.UseGitignore && gitignoreRoot != "" {
        pats, err := parseGitignore(gitignoreRoot)
        if err == nil {
            il.gitignorePatterns = pats
        }
    }
    return il
}

func (il *IgnoreList) ShouldSkip(name, rel string, isDir bool) bool {
    if isDir {
        if _, ok := il.ignoreDirs[name]; ok {
            return true
        }
    } else {
        if _, ok := il.ignoreFiles[name]; ok {
            return true
        }
        if _, ok := il.ignoreExts[filepath.Ext(name)]; ok {
            return true
        }
    }

    if il.useGitignore {
        for _, pat := range il.gitignorePatterns {
            if matchGitignorePattern(pat, name, rel, isDir) {
                return true
            }
        }
    }

    if il.matchesPatterns(name, rel, isDir, il.omitPatterns) {
        return true
    }

    if !isDir && len(il.includePatterns) > 0 {
        if !il.matchesPatterns(name, rel, false, il.includePatterns) {
            return true
        }
    }

    return false
}

func (il *IgnoreList) matchesPatterns(name, rel string, isDir bool, patterns []string) bool {
    for _, p := range patterns {
        if matchPattern(p, name, rel, isDir) {
            return true
        }
    }
    return false
}

func matchPattern(pattern, name, rel string, isDir bool) bool {
    if strings.Contains(pattern, string(filepath.Separator)) {
        matched, _ := filepath.Match(pattern, rel)
        return matched
    }
    if strings.HasSuffix(pattern, "/") {
        if isDir {
            dirPat := strings.TrimSuffix(pattern, "/")
            return matchBasePattern(dirPat, name)
        }
        return false
    }
    return matchBasePattern(pattern, name)
}

func matchBasePattern(p, name string) bool {
    if matched, _ := filepath.Match(p, name); matched {
        return true
    }
    if strings.Contains(p, "*") {
        mod := p
        if strings.HasPrefix(mod, "*") && !strings.HasSuffix(mod, "*") {
            mod += "*"
        }
        if strings.HasSuffix(mod, "*") && !strings.HasPrefix(mod, "*") {
            mod = "*" + mod
        }
        if matched, _ := filepath.Match(mod, name); matched {
            return true
        }
    } else if p != "" && strings.Contains(name, p) {
        return true
    }
    return false
}

func matchGitignorePattern(pat, name, rel string, isDir bool) bool {
    if strings.HasPrefix(pat, "!") {
        return false
    }
    if strings.HasPrefix(pat, "/") {
        pat = pat[1:]
    }
    dirOnly := false
    if strings.HasSuffix(pat, "/") {
        dirOnly = true
        pat = strings.TrimSuffix(pat, "/")
    }
    if strings.Contains(pat, "/") {
        if dirOnly && !isDir {
            return false
        }
        matched, _ := filepath.Match(pat, rel)
        return matched
    }
    if dirOnly && !isDir {
        return false
    }
    matched, _ := filepath.Match(pat, name)
    return matched
}

func copyMap(src map[string]struct{}) map[string]struct{} {
    dst := make(map[string]struct{}, len(src))
    for k := range src {
        dst[k] = struct{}{}
    }
    return dst
}

func parseGitignore(root string) ([]string, error) {
    data, err := os.ReadFile(filepath.Join(root, ".gitignore"))
    if err != nil {
        return nil, err
    }
    lines := strings.Split(string(data), "\n")
    var patterns []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        patterns = append(patterns, line)
    }
    return patterns, nil
}