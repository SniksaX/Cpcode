package src

import (
    "path/filepath"
    "sort"
    "strings"
)

type TreeNode struct {
    Children map[string]*TreeNode
}

func GenerateTree(files []string, rootDir string) string {
    rootNode := &TreeNode{Children: make(map[string]*TreeNode)}
    for _, f := range files {
        rel, err := filepath.Rel(rootDir, f)
        if err != nil {
            rel = f
        }
        parts := strings.Split(rel, string(filepath.Separator))
        current := rootNode
        for _, part := range parts {
            if current.Children[part] == nil {
                current.Children[part] = &TreeNode{Children: make(map[string]*TreeNode)}
            }
            current = current.Children[part]
        }
    }
    var sb strings.Builder
    sb.WriteString(filepath.Base(rootDir) + "/\n")
    walkTree(rootNode, "", &sb)
    return strings.TrimRight(sb.String(), "\n")
}

func walkTree(node *TreeNode, prefix string, sb *strings.Builder) {
    var keys []string
    for k := range node.Children {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for i, key := range keys {
        isLast := i == len(keys)-1
        connector := "├── "
        if isLast {
            connector = "└── "
        }
        sb.WriteString(prefix + connector + key + "\n")
        extension := "│   "
        if isLast {
            extension = "    "
        }
        if len(node.Children[key].Children) > 0 {
            walkTree(node.Children[key], prefix+extension, sb)
        }
    }
}