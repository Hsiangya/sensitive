package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"sensitive/controllers/factory"
	schemes "sensitive/controllers/schemes"
	services "sensitive/controllers/services"
	utils "sensitive/controllers/utils"
)

func SenFilterCreate(c *gin.Context) {
	info := schemes.SensitiveStringCreate{}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 100422, Msg: err.Error()})
		return
	}

	// construct documents
	ctx := context.Background()
	var documents []interface{}
	for _, v := range info.Text {
		sensitiveType := info.SensitiveType
		if sensitiveType == "" {
			sensitiveType = "未分类"
		}
		err := services.InsertWord(ctx, v)
		if err != nil {
			return
		}
		document := map[string]interface{}{"text": v, "sensitive_type": sensitiveType}
		documents = append(documents, document)

	}
	// insert mongo
	many, err := factory.MongoInstance.InsertMany(context.TODO(), "public_info", "sensitive", documents)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, &utils.ResponseContent{Code: 200, Msg: "没有敏感词", Data: many})
}

func SenFilterQuery(c *gin.Context) {
	info := schemes.SensitiveStringQuery{}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 100422, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &utils.ResponseContent{Code: 200, Msg: "没有敏感词"})
}
