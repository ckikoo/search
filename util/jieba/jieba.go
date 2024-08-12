package jieba

import (
	"strings"

	"github.com/go-ego/gse"
)

var GlobalSega gse.Segmenter

func init() {
	newGse, _ := gse.New()
	GlobalSega = newGse
}

func CutString(str string) []string {
	str = ignoredChar(str)
	// 使用搜索引擎模式进行分词
	return GlobalSega.CutSearch(str, false)
}

func CutStringForSe(str string) []string {
	str = ignoredChar(str)
	// 使用精准模式进行分词
	return GlobalSega.CutSearch(str, false)
}
func ignoredChar(str string) string {
	for _, c := range str {
		switch c {
		case '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', '【', '】', ']', '“', '”', '「', '」', '★', '^', '·', '_', '`', '{', '|', '}', '~', '《', '》', '：',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'（', '）', 0x3000, 0x3001, 0x3002, 0xFF01, 0xFF0C, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), " ")
		}
	}
	return str
}
