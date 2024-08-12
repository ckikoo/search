package parse

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Info struct {
	Title   string
	Content string
	Url     string
}

func parseTitle(content string) string {
	contents := strings.Split(content, "<title>")
	contents = strings.Split(contents[1], "</title>")
	return contents[0]
}

func parseContent(file string) string {
	const (
		LABEL = iota
		CONTENT
	)
	var state int = LABEL
	var content strings.Builder

	for _, c := range file {
		switch state {
		case LABEL:
			if c == '>' {
				state = CONTENT
			}
		case CONTENT:
			if c == '<' {
				state = LABEL
			} else {
				if c == '\n' {
					c = 0
				}
				if c != 0 { // 忽略被替换为0的字符
					content.WriteRune(c)
				}
			}
		}
	}
	return content.String()
}

func parseUrl(path string) string {
	return path
}

// parseFile 处理单个文件的逻辑
func parseFile(filePath string) *Info {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	buff, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	title := parseTitle(string(buff))
	Content := parseContent(string(buff))
	Url := parseUrl(filePath)

	return &Info{
		Title:   title,
		Content: Content,
		Url:     Url,
	}
}

// Parse 递归遍历目录中的所有文件并调用 parseFile 处理每个文件
func Parse(dir string) []Info {
	infos := make([]Info, 0)
	fmt.Printf("dir: %v\n", dir)
	info, err := os.Stat(dir)
	if err != nil {
		return infos
	}

	if !info.IsDir() {
		info1 := parseFile(dir)
		infos = append(infos, *info1)
		return infos
	}
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return infos
	}

	for _, entry := range dirs {
		entryPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// 如果是目录，递归调用 Parse
			lists := Parse(entryPath)
			infos = append(infos, lists...)
		} else {
			// 如果是文件，调用 parseFile
			if path.Ext(entry.Name()) == ".html" {
				info := parseFile(entryPath)
				infos = append(infos, *info)
			}
		}
	}

	return infos
}
