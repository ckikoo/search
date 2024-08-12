package index

import (
	"ckikoo/search/model/parse"
	"ckikoo/search/util/jieba"
	"ckikoo/search/util/log"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type DocInfo struct { // 文档信息
	Title   string // 标题
	Content string // 正文
	Url     string // url
	DocId   uint64 // 文档编号
}

type InvertedElem struct {
	DocId  uint64 // 编号
	Word   string // 单词
	Weight int64  // 权重
}

type InvertedList_t []InvertedElem

type Index struct {
	Forward_index  []DocInfo
	Inverted_index map[string]InvertedList_t
	locker         sync.Mutex
}

var (
	instance *Index
	once     sync.Once
)

// GetInstance 返回 Index 的单例实例
func GetInstance() *Index {
	once.Do(func() {
		instance = &Index{
			Forward_index:  make([]DocInfo, 0),
			Inverted_index: make(map[string]InvertedList_t),
		}
	})
	return instance
}

func (index *Index) GetForwardIndex(doc_id int) *DocInfo {
	if doc_id > len(index.Forward_index) {
		return nil
	}

	return &index.Forward_index[doc_id]
}
func (index *Index) GetInvertedList(keyWord string) *InvertedList_t {

	res, ok := index.Inverted_index[keyWord]
	if !ok {
		return nil
	}

	return &res
}

// BuildIndex 读取文件并构建索引
func (index *Index) BuildIndex(dir string) bool {
	infos := parse.Parse(dir)
	pool, err := ants.NewPool(1024)
	if err != nil {
		log.Error("创建池子失败")
		return false
	}
	defer pool.Release()
	wg := sync.WaitGroup{}
	fmt.Printf("len(infos): %v\n", len(infos))
	wg.Add(len(infos))
	for _, info := range infos {
		pool.Submit(func() {
			defer wg.Done()
			fmt.Printf("info.Title: %v\n", info.Title)
			info1 := index.buildForWardIndex(info)

			if nil == info1 {
				log.Warning("buildForWardIndex nil, content:%v", info.Title)
				return
			}
			fmt.Printf("info1.Title: %v\n", info1.Title)
			index.buildInvertedIndex(info1)
		})
	}

	wg.Wait()
	fmt.Printf("len(index.forward_index): %v\n", len(index.Forward_index))
	fmt.Printf("\"构建成功\": %v\n", "构建成功")
	return true
}

func (index *Index) buildForWardIndex(line parse.Info) *DocInfo {

	index.locker.Lock()
	defer index.locker.Unlock()
	doc := DocInfo{
		Title:   line.Title,
		Content: line.Content,
		Url:     line.Url,
		DocId:   uint64(len(index.Forward_index)), // 使用当前索引长度作为 doc_id
	}

	index.Forward_index = append(index.Forward_index, doc)
	return &doc
}

func (index *Index) buildInvertedIndex(info *DocInfo) bool {
	type word_cnt struct {
		title_cnt   int64
		content_cnt int64
	}
	wordCountMap := make(map[string]*word_cnt)
	titleWords := jieba.CutString(info.Title)
	fmt.Printf("titleWords: %v\n", titleWords)
	for _, word := range titleWords {
		if len(word) == 0 || word == " " {
			continue
		}
		if _, exists := wordCountMap[word]; !exists {
			wordCountMap[word] = &word_cnt{}
		}
		wordCountMap[word].title_cnt++
	}

	// 处理正文中的词汇
	contentWords := jieba.CutString(info.Content)
	for _, word := range contentWords {
		if len(word) == 0 || word == " " {
			continue
		}
		if _, exists := wordCountMap[word]; !exists {
			wordCountMap[word] = &word_cnt{}
		}
		wordCountMap[word].content_cnt++
	}

	index.locker.Lock()         // 手动锁定
	defer index.locker.Unlock() // 在函数结束时解锁
	// 根据词汇出现的频率来构建倒排索引
	for word, count := range wordCountMap {
		weight := count.title_cnt*123456 + count.content_cnt // 这里根据需要调整权重计算方式
		elem := InvertedElem{
			DocId:  info.DocId,
			Word:   word,
			Weight: weight,
		}

		index.Inverted_index[word] = append(index.Inverted_index[word], elem)
	}

	return true
}
