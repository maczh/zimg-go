package main

import (
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"ququ.im/common"
	"ququ.im/common/aop"
	"ququ.im/zimg-go/controller"
	_ "ququ.im/zimg-go/docs"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	engine := gin.Default()

	//添加跟踪日志
	engine.Use(aop.TraceId())

	//设置接口日志
	engine.Use(aop.SetRequestLogger())

	//处理全局异常
	engine.Use(nice.Recovery(recoveryHandler))

	//添加swagger支持
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//设置404返回的内容
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, *common.Error(-1, "请求的方法不存在"))
	})

	var result common.Result
	//添加所需的路由映射
	engine.GET("/pic", func(c *gin.Context) {
		imgUrl := controller.GetPictureResized(c.GetParamMap())
		if imgUrl == "" {
			c.JSON(http.StatusOK, *common.Error(-1, "图片不存在"))
		} else {
			c.Redirect(http.StatusMovedPermanently, imgUrl)
		}
	})
	engine.Any("/upload/from", func(c *gin.Context) {
		result = controller.SaveImageFromUrl(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.Any("/del", func(c *gin.Context) {
		result = controller.DeleteImage(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.Any("/list/path", func(c *gin.Context) {
		result = controller.ListImageByPath(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.Any("/list/url", func(c *gin.Context) {
		result = controller.ListImageByUrl(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.Any("/list/url/bytime", func(c *gin.Context) {
		result = controller.ListSourceUrlByUploadTime(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.Any("/list/remark", func(c *gin.Context) {
		result = controller.ListImageByRemark(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.POST("/upload", func(c *gin.Context) {
		result = controller.UploadImage(c)
		c.JSON(http.StatusOK, result)
	})

	//生意宝资质查相关
	engine.Any("/qua/upload/from", func(c *gin.Context) {
		result = controller.UploadQualificationToZimg(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.GET("/qua/upload/all", func(c *gin.Context) {
		result = controller.UploadQualificationToZimgAll(c.GetParamMap())
		c.JSON(http.StatusOK, result)
	})
	engine.POST("/qua/upload", func(c *gin.Context) {
		result = controller.UploadQualificationImage(c)
		c.JSON(http.StatusOK, result)
	})
	return engine
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, *common.Error(-1, "系统异常，请联系客服"))
}
