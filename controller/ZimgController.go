package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sadlil/gologger"
	"os"
	"path/filepath"
	"ququ.im/common"
	"ququ.im/common/utils"
	"ququ.im/zimg-go/service"
	"ququ.im/zimg-go/zimg"
	"strconv"
	"time"
)

var logger = gologger.GetLogger()

// SaveImageFromUrl	godoc
// @Summary		转存网络图片到Zimg图床
// @Description	转存网络图片到Zimg图床
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	imgUrl formData string true "网络图片URL"
// @Param	path formData string false "文件路径"
// @Param	remark formData string false "图片备注描述"
// @Success 200 {string} string	"ok"
// @Router	/upload/from [post][get]
func SaveImageFromUrl(params map[string]string) common.Result {
	fileId, err := service.SaveImageFromUrl(params["imgUrl"], params["path"], params["remark"])
	if err != nil {
		return *common.Error(-1, err.Error())
	}
	result := make(map[string]string)
	result["id"] = fileId
	result["url"] = zimg.ZIMG_HOST + fileId
	return *common.Success(result)
}

// GetPictureResized	godoc
// @Summary		获取限定文件大小的图片
// @Description	获取限定文件大小的图片
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	id formData string true "图片编号id"
// @Param	maxSize formData int false "图片大小上限，单位为kb，默认为最大10Mb"
// @Success 200 {string} string	"ok"
// @Router	/pic [get]
func GetPictureResized(params map[string]string) string {
	if !utils.Exists(params, "id") {
		return ""
	}
	limit := 10 * 1024 //默认最大10M
	if utils.Exists(params, "maxSize") {
		limit, _ = strconv.Atoi(params["maxSize"])
	}
	picUrl, err := service.GetPictureResized(params["id"], limit)
	if err != nil {
		return ""
	}
	return picUrl
}

// DeleteImage	godoc
// @Summary		删除指定图片
// @Description	删除指定图片
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	id formData string true "图片编号id"
// @Success 200 {string} string	"ok"
// @Router	/del [post][get]
func DeleteImage(params map[string]string) common.Result {
	return service.DeleteImage(params["id"])
}

// ListImageByPath	godoc
// @Summary		按指定路径查询图片列表
// @Description	按指定路径查询图片列表
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	path formData string true "图片路径,可以使用*通配符"
// @Success 200 {string} string	"ok"
// @Router	/list/path [post][get]
func ListImageByPath(params map[string]string) common.Result {
	return service.ListImageByPath(params["path"])
}

// ListImageByUrl	godoc
// @Summary		按指定图片源地址查询
// @Description	按指定图片源地址查询
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	img formData string true "图片源地址"
// @Success 200 {string} string	"ok"
// @Router	/list/url [post][get]
func ListImageByUrl(params map[string]string) common.Result {
	return service.ListImageByUrl(params["img"])
}

// ListSourceUrlByUploadTime	godoc
// @Summary		按起止时间获取源地址清单
// @Description	按起止时间获取源地址清单
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	beginDate formData string false "开始日期，格式yyyy-MM-dd，默认为当天"
// @Param	endDate formData string false "结束日期，格式yyyy-MM-dd，默认为明天"
// @Success 200 {string} string	"ok"
// @Router	/list/url/bytime [post][get]
func ListSourceUrlByUploadTime(params map[string]string) common.Result {
	return service.ListSourceUrlByUploadTime(params["beginDate"], params["endDate"])
}

// ListImageByRemark	godoc
// @Summary		按图片备注查询
// @Description	按图片备注查询
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	remark formData string true "图片备注信息,可以使用*通配符"
// @Success 200 {string} string	"ok"
// @Router	/list/remark [post][get]
func ListImageByRemark(params map[string]string) common.Result {
	return service.ListImageByRemark(params["remark"])
}

// UploadImage	godoc
// @Summary		上传图片
// @Description	上传图片
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	image formData file true "图片文件数据"
// @Param	path formData string false "图片路径"
// @Param	remark formData string false "图片备注信息"
// @Param	fileName formData string true "本地图片文件名"
// @Success 200 {string} string	"ok"
// @Router	/upload [post]
func UploadImage(c *gin.Context) common.Result {
	imgFile, err := c.FormFile("image")
	if err != nil {
		logger.Error("获取上传文件失败:" + err.Error())
		return *common.Error(-1, "获取上传文件失败:"+err.Error())
	}
	path := c.PostForm("path")
	if path == "" {
		path = "/tmp/" + time.Now().Format("20060102")
	}
	remark := c.PostForm("remark")
	fileName := c.PostForm("fileName")
	localpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	localFileName := localpath + zimg.IMAGE_PATH + fileName
	err = c.SaveUploadedFile(imgFile, localFileName)
	if err != nil {
		return *common.Error(-1, "保存上传文件"+localFileName+"失败："+err.Error())
	}
	return service.UploadImageFile(localFileName, path, remark, fileName)
}
