package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	factory "sensitive/controllers/factory"
	"sensitive/controllers/schemes"
	utils "sensitive/controllers/utils"
)

func SenFilterCreate(c *gin.Context) {
	info := schemes.SensitiveStringCreate{}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 100422, Msg: err.Error()})
		return
	}

	var documents []interface{}
	for _, v := range info.Text {
		sensitiveType := info.SensitiveType
		if sensitiveType == "" {
			sensitiveType = "未分类"
		}
		factory.DfaInstance.AddWord(v)

		document := map[string]interface{}{"text": v, "sensitive_type": sensitiveType}
		documents = append(documents, document)

	}
	// insert mongo
	insertResult, err := factory.MongoInstance.InsertMany(context.TODO(), "public_info", "sensitive", documents)
	if err != nil {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 100500, Msg: "敏感词入库失败"})
		return
	}
	c.JSON(http.StatusOK, &utils.ResponseContent{Code: 200, Msg: "没有敏感词", Data: insertResult})
}

func SenFilterQuery(c *gin.Context) {
	query := schemes.SensitiveStringQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 100422, Msg: err.Error()})
		return
	}
	isChineseSensitive := factory.DfaInstance.CheckChinese(query.Text)
	isEnglishSensitive := factory.DfaInstance.CheckEnglish(query.Text)

	if isChineseSensitive || isEnglishSensitive {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 200, Msg: "存在敏感词"})
	} else {
		c.JSON(http.StatusOK, &utils.ResponseContent{Code: 200, Msg: "没有敏感词"})
	}
}
