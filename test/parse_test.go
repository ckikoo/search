package test

import (
	"ckikoo/search/model/index"
	"encoding/json"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	src := "../data"
	index.GetInstance().BuildIndex(src)

	buff, err := json.Marshal(index.GetInstance())
	if err != nil {
		panic(err)
	}

	os.WriteFile("lo.json", buff, 0644)

	// fmt.Printf("index.GetInstance(): %+v\n", index.GetInstance().GetInvertedList("固"))
	// fmt.Printf("index.GetInstance(): %+v\n", index.GetInstance())
}
