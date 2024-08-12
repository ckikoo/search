package router

import (
	"ckikoo/search/model/search"
	"ckikoo/search/util/log"

	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) *gin.Engine {
	engine.Static("/data", "data/")
	engine.GET("/query", func(ctx *gin.Context) {
		queryWords := ctx.Query("word")

		log.Info("user query word: [%v]", queryWords)
		var t search.Search
		res := t.S(queryWords)

		ctx.JSON(200, gin.H{
			"c": res,
		})
	})

	return engine
}
