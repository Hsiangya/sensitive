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

func CreateApp(env string) {
	initConfig(env)
	initPlugin()
	initValidation()
	app := createGinApp()
	fmt.Println(viper.GetString("application.port"))
	err := app.Run(fmt.Sprintf(":%s", viper.GetString("application.port")))
	if err != nil {
		return
	}
}

func initConfig(env string) {
	configMapping := map[string]string{
		"dev":  "./configs/app-dev.yaml",
		"test": "./configs/app-test.yaml",
		"pro":  "./configs/app-pro.yaml",
	}
	configPath, exist := configMapping[env]
	if !exist {
		configPath = configMapping["dev"]
	}
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}
}

func createGinApp() *gin.Engine {
	g := gin.Default()
	g.MaxMultipartMemory = 24 << 20 // 8 MiB
	response := &utils.ResponseContent{Code: 100404, Msg: "路由不存在"}
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, response)
	})
	g.POST("/sensitive", api.SenFilterCreate)
	g.GET("/sensitive", api.SenFilterQuery)
	return g
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
	factory.CreateDFATree()

	//factory.CreateMongoApp(mongoURL)
}
