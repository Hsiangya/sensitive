package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"sensitive/controllers/api"
	"sensitive/controllers/factory"
	"sensitive/controllers/utils"
)

func Run() {
	initConfig()
	port := viper.GetString("server.port")
	router := gin.Default()
	router.MaxMultipartMemory = 24 << 20 // 8 MiB
	initPlugin()
	initValidation()
	InitRouter(router)
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		return
	}
}

func initConfig() {
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}
}

func InitRouter(g *gin.Engine) {
	//注册404路由
	response := &utils.ResponseContent{Code: 100404, Msg: "路由不存在"}
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, response)
	})
	g.POST("/sensitive", api.SenFilterCreate)
	g.GET("/sensitive", api.SenFilterQuery)

}

func initValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidators := map[string]validator.Func{
			//"TargetLangValidator": schemes.TargetLangValidator,
		}
		for tag, validationFunc := range customValidators {
			err := v.RegisterValidation(tag, validationFunc)
			if err != nil {
				return
			}
		}
	}
}

func initPlugin() {
	mongoURL := viper.GetString("mongo.url")
	factory.CreateMongoApp(mongoURL)
}
