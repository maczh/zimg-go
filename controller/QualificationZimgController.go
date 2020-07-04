package controller

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"ququ.im/common"
	"ququ.im/zimg-go/service"
	"ququ.im/zimg-go/zimg"
)

// UploadQualificationToZimg	godoc
// @Summary		单个生意宝资质上传到Zimg服务器
// @Description	单个生意宝资质上传到Zimg服务器
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	accountId formData string true "生意宝账户ID"
// @Success 200 {string} string	"ok"
// @Router	/qua/upload/from [post][get]
func UploadQualificationToZimg(params map[string]string) common.Result {
	return service.UploadQualificationToZimg(params["accountId"])
}

// UploadQualificationToZimg	godoc
// @Summary		全部生意宝资质上传到Zimg服务器
// @Description	全部生意宝资质上传到Zimg服务器
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	start formData string false "起始生意宝账户ID"
// @Success 200 {string} string	"ok"
// @Router	/qua/upload/all [get]
func UploadQualificationToZimgAll(params map[string]string) common.Result {
	return service.UploadQualificationToZimgAll(params["start"])
}

// UploadQualificationImage	godoc
// @Summary		上传生意宝资质图片
// @Description	上传生意宝资质图片
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	image formData file true "图片文件数据"
// @Param	accountId formData string true "生意宝账号"
// @Param	picType formData string true "资质图片类型"
// @Param	bankCardNo formData string false "银行卡号码，若是传银行卡正面照片时必填"
// @Success 200 {string} string	"ok"
// @Router	/qua/upload [post]
func UploadQualificationImage(c *gin.Context) common.Result {
	imgFile, err := c.FormFile("image")
	if err != nil {
		return *common.Error(-1, "获取上传文件失败:"+err.Error())
	}
	accountId := c.PostForm("accountId")
	if accountId == "" {
		return *common.Error(-1, "缺少资质生意宝账号ID")
	}
	picType := c.PostForm("picType")
	if picType == "" {
		return *common.Error(-1, "缺少资质类型")
	}
	bankCardNo := c.PostForm("bankCardNo")
	fileName := picType + "-" + accountId + ".jpg"
	localpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	localFileName := localpath + zimg.IMAGE_PATH + fileName
	err = c.SaveUploadedFile(imgFile, localFileName)
	if err != nil {
		return *common.Error(-1, "保存上传文件"+localFileName+"失败："+err.Error())
	}
	return service.UploadQualificationImage(accountId, picType, bankCardNo, localFileName)
}
