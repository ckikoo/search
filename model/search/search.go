package search

import (
	"ckikoo/search/model/index"
	"ckikoo/search/util/jieba"
	"fmt"
	"sort"
)

type InvertedElemPrint struct {
	DocId  int
	Weight int
}

type res struct {
	Url   string
	Title string
}

type Search struct{}

// S 方法用于通过关键词搜索相关文档
func (s *Search) S(word string) []res {
	// 使用 jieba 对输入词语进行分词
	words := jieba.CutStringForSe(word)
	fmt.Printf("words: %v\n", words)

	// 用于存储所有的倒排索引结果
	var invertedListAll []InvertedElemPrint
	// 用于记录每个文档的累计权重
	tokensMap := make(map[int]*InvertedElemPrint)

	// 遍历每个分词结果
	for _, word := range words {
		fmt.Printf("word: %v\n", word)
		invertedList := index.GetInstance().GetInvertedList(word)
		if invertedList == nil {
			continue // 如果当前词没有对应的倒排索引，跳过
		}
		fmt.Printf("invertedList: %v\n", invertedList)
		for _, elem := range *invertedList {
			// 查找 tokensMap 中是否已存在该文档
			item, ok := tokensMap[int(elem.DocId)]
			if !ok {
				// 如果不存在，则创建一个新的 InvertedElemPrint 对象
				item = &InvertedElemPrint{
					DocId:  int(elem.DocId),
					Weight: 0,
				}
				tokensMap[int(elem.DocId)] = item
			}

			// 累计当前文档的权重
			item.Weight += int(elem.Weight)
		}
	}

	// 预先分配切片的容量
	invertedListAll = make([]InvertedElemPrint, len(tokensMap))

	// 将 map 中的结果转为 slice 以便排序
	i := 0
	for _, item := range tokensMap {
		invertedListAll[i] = *item
		i++
	}

	// 按权重从高到低排序
	sort.Slice(invertedListAll, func(i, j int) bool {
		return invertedListAll[i].Weight > invertedListAll[j].Weight
	})

	r := make([]res, len(invertedListAll))
	// 输出搜索结果
	i = 0
	for _, item := range invertedListAll {
		info := index.GetInstance().GetForwardIndex(item.DocId)
		r[i] = res{Url: info.Url, Title: info.Title}
		i++
	}

	return r
}
