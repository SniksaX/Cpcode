<img width="224" height="44" alt="cpcode-lockup-white" src="https://github.com/user-attachments/assets/171294f6-e732-487d-a8e0-8dce034956d7" />


**`cpcode`** is a fast CLI tool written in Go that packages your source code, directory structures, and terminal outputs into single, well-formatted prompt contexts optimized for Large Language Models (LLMs) like ChatGPT, Claude, and Gemini.

By default, `cpcode` formats your codebase as Markdown, appends an ASCII file tree with execution stats, and copies the result directly to your system clipboard in milliseconds.

---

## Key Features

- **Blazing Fast Concurrent Processing**: Leverages multi-core parallel worker pools to scan and process large repositories in milliseconds.
- **Clipboard-First Workflow**: Copies output directly to your system clipboard (`pbcopy`, `xclip`, `wl-copy`, Windows clipboard) with automatic fallback to `stdout`.
- **Automatic Tree Generation**: Generates clean ASCII directory structures at the bottom of your context payload.
- **Smart Default Filtering**: Automatically skips binaries, images, archives, lock files (`package-lock.json`, `Cargo.lock`), dependency directories (`node_modules`, `vendor`, `.venv`), and files over 1 MB.
- **Smart Appending**: Easily chain multiple directories, files, or stdin terminal streams together with `--append` without breaking context syntax.
- **LLM Token Estimator**: Built-in `--tokens` flag estimates context size before sending prompts to your model.
- **Piping Support**: Pipe terminal logs, diffs, or standard output directly into your context (`cat log.txt | cpcode`).
- **Multiple Formats**: Export output as `markdown` (default), `xml`, or `plain` text.

---

## Installation

### From Source

```bash
git clone git@github.com:SniksaX/Cpcode.git
cd cpcode
go build -o cpcode main.go
mv cpcode /usr/local/bin/
```

---

## Quick Start

```bash
# 1. Copy current directory to clipboard (Markdown format)
cpcode

# 2. Copy specific files or subdirectories
cpcode main.go src/

# 3. Print output to stdout instead of clipboard
cpcode -p

# 4. Save output to a file
cpcode -d llm_context.md

# 5. Append terminal logs to your clipboard context
git status | cpcode -a
```

---

## Usage & Options

```text
cpcode [flags] [files/directories...]
```

### Flags Reference

| Short | Long | Description | Default |
| :--- | :--- | :--- | :--- |
| `-d` | `--destination` | Save output to a specified file | `""` |
| `-p` | `--print` | Print payload to `stdout` instead of clipboard | `false` |
| `-a` | `--append` | Append output to existing file or clipboard | `false` |
| `-t` | `--tokens` | Include an estimated LLM token count | `false` |
| `-o` | `--omit` | Patterns/files to ignore (comma-separated) | `""` |
| `-i` | `--include` | Patterns/files to explicitly include | `""` |
| | `--dry-run` | Show execution summary without copying/saving | `false` |
| | `--format` | Output format (`markdown`, `xml`, `plain`) | `markdown` |
| | `--use-gitignore`| Respect root `.gitignore` rules | `false` |
| | `--no-ignore` | Disable default built-in ignore lists | `false` |
| `-q` | `--quiet` | Suppress CLI summary stderr output | `false` |

---

## Examples

### 1. Filtering Specific Files

**Omit tests or log files:**
```bash
cpcode -o "*.test.go, *.log, docs/"
```

**Include only specific extensions:**
```bash
cpcode -i "*.go, *.ts"
```

### 2. Respecting `.gitignore`

By default, `cpcode` uses its smart built-in ignore list. To additionally enforce your repository's `.gitignore`:

```bash
cpcode --use-gitignore
```

### 3. Appending Stdin / Logs

If you're debugging an issue, copy your source code first, then append the error stack trace from your terminal:

```bash
# Step 1: Copy main source files
cpcode src/

# Step 2: Pipe terminal output to append to clipboard context
go test ./... | cpcode -a
```

### 4. XML Formatting & Token Estimates

For models like Claude that perform well with XML tag formatting:

```bash
cpcode --format xml -t
```

### 5. Dry Run (Previewing Context)

Preview which files will be collected, the ASCII tree structure, and statistics without affecting your clipboard:

```bash
cpcode --dry-run
```

---

## Output Format Preview

By default, `cpcode` outputs Markdown syntax blocks with pretty file paths and appends execution statistics:

````markdown
```go
// ~/Desktop/other/cpcode/package/main.go
package main

import (
    "fmt"
    ...
)
...
```

<!-- cpcode:footer -->
==================================================
Project Structure:
package/
├── main.go
└── src
    ├── clipboard.go
    ├── flags.go
    ├── ignore.go
    ├── processor.go
    ├── tree.go
    ├── utils.go
    └── walker.go

--- Stats ---
Files/Streams Added : 8
Total Size          : 29.7 KB
Time Taken          : 591.129µs
````

---

## Default Ignore Lists

`cpcode` automatically skips common non-essential context items to prevent excessive token usage:

- **Directories**: `.git`, `node_modules`, `dist`, `build`, `bin`, `vendor`, `__pycache__`, `.venv`, `.vscode`, `.idea`, etc.
- **Files**: `.DS_Store`, `package-lock.json`, `yarn.lock`, `Cargo.lock`, `*.lock`, `.env`, binary files, etc.
- **Extensions**: `.exe`, `.dll`, `.png`, `.jpg`, `.pdf`, `.zip`, `.mp4`, `.db`, `.pyc`, `.woff2`, etc.
- **File Size**: Any file exceeding 1 MB or containing binary bytes (`\0`).

*(Note: Use `--no-ignore` to bypass built-in ignore rules).*

---

## Contributing

Contributions are welcome. Feel free to open issues or submit pull requests to enhance performance, add support for more platforms, or introduce useful options.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---
