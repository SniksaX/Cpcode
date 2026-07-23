package src

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "runtime"
    "sort"
    "strings"
    "sync"
    "sync/atomic"
)

const maxFileSize = 1_000_000

type FileResult struct {
    Path    string
    Content []byte
}

func ProcessFiles(files []string) ([]FileResult, int64, int64) {
    numWorkers := runtime.NumCPU() * 4
    jobs := make(chan int, len(files))
    resultsChan := make(chan FileResult, len(files))
    var wg sync.WaitGroup
    var validCount, skippedCount int64

    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := range jobs {
                path := files[i]
                content, skip := readTextFile(path)
                if skip {
                    atomic.AddInt64(&skippedCount, 1)
                    continue
                }
                resultsChan <- FileResult{Path: prettyPath(path), Content: content}
                atomic.AddInt64(&validCount, 1)
            }
        }()
    }

    for i := range files {
        jobs <- i
    }
    close(jobs)
    wg.Wait()
    close(resultsChan)

    var results []FileResult
    for res := range resultsChan {
        results = append(results, res)
    }
    sort.Slice(results, func(i, j int) bool { return results[i].Path < results[j].Path })
    return results, validCount, skippedCount
}

func readTextFile(path string) ([]byte, bool) {
    f, err := os.Open(path)
    if err != nil {
        return nil, true
    }
    defer f.Close()

    fi, err := f.Stat()
    if err != nil || fi.Size() == 0 || fi.Size() > maxFileSize {
        return nil, true
    }

    header := make([]byte, 8192)
    n, err := f.Read(header)
    if err != nil && err != io.EOF {
        return nil, true
    }
    if bytes.IndexByte(header[:n], 0) != -1 {
        return nil, true
    }

    var content []byte
    if fi.Size() <= int64(n) {
        content = make([]byte, n)
        copy(content, header[:n])
    } else {
        content = make([]byte, fi.Size())
        copy(content, header[:n])
        _, err = io.ReadFull(f, content[n:])
        if err != nil && err != io.ErrUnexpectedEOF {
            return nil, true
        }
    }

    if len(bytes.TrimSpace(content)) == 0 {
        return nil, true
    }
    return content, false
}

func FormatContent(fr FileResult, format string) string {
    var b strings.Builder
    switch format {
    case "xml":
        b.WriteString(fmt.Sprintf("<file path=\"%s\">\n<content><![CDATA[", fr.Path))
        b.Write(fr.Content)
        b.WriteString("]]></content>\n</file>\n")
    case "plain":
        b.Write(fr.Content)
        b.WriteString("\n---\n")
    default:
        ext := strings.TrimPrefix(filepath.Ext(fr.Path), ".")
        if ext == "" {
            ext = "txt"
        }
        b.WriteString("```")
        b.WriteString(ext)
        b.WriteString("\n// ")
        b.WriteString(fr.Path)
        b.WriteString("\n")
        b.Write(fr.Content)
        b.WriteString("\n```\n\n")
    }
    return b.String()
}