package service

import (
	"errors"
	"gitee.com/vchuzu/imagedecode"
	"github.com/levigross/grequests"
	"github.com/sadlil/gologger"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"ququ.im/common"
	"ququ.im/zimg-go/dao"
	"ququ.im/zimg-go/pojo"
	"ququ.im/zimg-go/zimg"
	"strings"
	"time"
)

var logger = gologger.GetLogger()

/**
从网络图片转存到zimg
*/
func SaveImageFromUrl(imgUrl, path, remark string) (string, error) {
	if imgUrl == "" {
		return "", errors.New("缺少图片URL")
	}
	u, _ := url.Parse(imgUrl)
	if path == "" {
		path = u.Host
	}
	if strings.Contains(imgUrl, zimg.ZIMG_HOST) {
		return "", errors.New("源图地址不能是zimg地址")
	}
	if strings.Contains(imgUrl, "?") {
		idx := strings.Index(imgUrl, "?")
		imgUrl = imgUrl[0:idx]
	}
	filename := strings.TrimLeft(u.Path, "/")
	localpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := localpath + zimg.IMAGE_PATH + filename

	//下载网络图片
	resp, err := grequests.Get(imgUrl, nil)
	if err != nil {
		logger.Error("源文件下载失败:" + err.Error())
		return "", err
	}
	//if resp.Ok {
	err = resp.DownloadToFile(fileName)
	if err != nil {
		logger.Error("源文件下载失败:" + err.Error())
		return "", err
	}
	//} else {
	//	logger.Error("源文件下载失败:" + resp.String())
	//	return "", errors.New("源文件下载失败:" + resp.String())
	//}

	//上传文件到Zimg
	fileId, err := zimg.Upload(fileName)
	if err != nil {
		logger.Error("文件上传到zimg失败:" + err.Error())
		return "", err
	}
	//删除本地缓存文件
	err = os.Remove(fileName)
	if err != nil {
		logger.Error("删除本地缓存文件失败:" + err.Error())
		return "", err
	}

	//重新下载到本地检验，并获取文件大小
	size, _, err := zimg.Download(fileId, fileName, 0, 0, -1, 0, 0, 0, 0)
	if err != nil {
		logger.Error("从zimg下载文件失败:" + err.Error())
		return "", err
	}

	imgType, err := imagedecode.GetImageType(fileName)
	if err != nil {
		logger.Error("不是图片类型:" + err.Error())
		return "", errors.New("不是图片类型:" + err.Error())
	}
	ext := ""
	switch imgType {
	case imagedecode.GIF_TYPE:
		ext = "gif"
	case imagedecode.JPG_TYPE:
		ext = "jpg"
	case imagedecode.PNG_TYPE:
		ext = "png"
	}

	//删除本地缓存文件
	err = os.Remove(fileName)
	if err != nil {
		logger.Error("删除本地缓存文件失败:" + err.Error())
		return "", err
	}

	//保存到MySQL
	imageFile := dao.GetImage(fileId)
	if imageFile != nil {
		//已存在，直接返回
		return fileId, nil
	}
	imageFile = new(pojo.ImageFiles)
	imageFile.Id = fileId
	imageFile.SourceUrl = imgUrl
	imageFile.Size = size
	imageFile.ImgType = ext
	imageFile.Remark = remark
	imageFile.ServerPath = path
	imageFile.FileName = fileName
	imageFile, err = dao.SaveImage(imageFile)
	if err != nil {
		logger.Error("图片信息入库失败:" + err.Error())
		return "", err
	}
	return fileId, nil
}

/**
按限定文件大小(kb)自动缩放获取相应图片URL
*/
func GetPictureResized(fileId string, maxSize int) (string, error) {
	imageFile := dao.GetImage(fileId)
	if imageFile == nil {
		logger.Error(fileId + "图片不存在")
		return "", errors.New("图片不存在")
	}
	if maxSize*1024 >= imageFile.Size {
		return zimg.GetPictureUrl(fileId, 0, 0, -1, 0, 0, 0, 0), nil
	} else {
		scale := int(100*math.Sqrt(float64(maxSize*1024)/float64(imageFile.Size))) - 2
		return zimg.GetPictureUrl(fileId, scale, scale, 3, 0, 0, 0, 0), nil
	}
}

/**
删除一个图片
*/
func DeleteImage(fileId string) common.Result {
	if fileId == "" {
		return *common.Error(-1, "缺少图片编号")
	}
	imageFile := dao.GetImage(fileId)
	if imageFile == nil {
		logger.Error(fileId + "图片不存在")
		return *common.Error(-1, "图片不存在")
	}
	err := zimg.Delete(fileId)
	if err != nil {
		logger.Error("图片删除错误:" + err.Error())
		return *common.Error(-1, "图片删除错误:"+err.Error())
	}
	err = dao.DeleteImage(fileId)
	if err != nil {
		logger.Error("图片库删除错误:" + err.Error())
		return *common.Error(-1, "图片库删除错误:"+err.Error())
	}
	return *common.Success(nil)
}

func ListImageByPath(path string) common.Result {
	if path == "" {
		return *common.Error(-1, "缺少路径")
	}
	imageFiles := dao.ListImageByPath(path)
	if imageFiles == nil || len(imageFiles) == 0 {
		return *common.Error(-1, "无此图片")
	} else {
		return *common.Success(imageFiles)
	}
}

func ListImageByUrl(imgUrl string) common.Result {
	if imgUrl == "" {
		return *common.Error(-1, "缺少源图片网址")
	}
	imageFiles := dao.ListImageByUrl(imgUrl)
	if imageFiles == nil || len(imageFiles) == 0 {
		return *common.Error(-1, "无此图片")
	} else {
		return *common.Success(imageFiles)
	}
}

func ListImageByRemark(remark string) common.Result {
	if remark == "" {
		return *common.Error(-1, "缺少图片说明")
	}
	imageFiles := dao.ListImageByRemark(remark)
	if imageFiles == nil || len(imageFiles) == 0 {
		return *common.Error(-1, "无此图片")
	} else {
		return *common.Success(imageFiles)
	}
}

func ListSourceUrlByUploadTime(beginDate, endDate string) common.Result {
	if beginDate == "" {
		beginDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Add(24 * time.Hour).Format("2006-01-02")
	}
	sourceUrls := dao.ListSourceUrlByUploadTime(beginDate, endDate)
	if sourceUrls == nil {
		return *common.Error(-1, "查无数据")
	}
	return *common.Success(sourceUrls)
}

func UploadImageFile(localFileName, path, remark, fileName string) common.Result {
	//上传文件到Zimg
	fileId, err := zimg.Upload(localFileName)
	if err != nil {
		logger.Error("文件上传到zimg失败:" + err.Error())
		return *common.Error(-1, "文件上传到zimg失败:"+err.Error())
	}
	//删除本地缓存文件
	err = os.Remove(localFileName)
	if err != nil {
		logger.Error("删除本地缓存文件失败:" + err.Error())
		return *common.Error(-1, "删除本地缓存文件失败:"+err.Error())
	}

	//重新下载到本地检验，并获取文件大小
	size, _, err := zimg.Download(fileId, localFileName, 0, 0, -1, 0, 0, 0, 0)
	if err != nil {
		logger.Error("从zimg下载文件失败:" + err.Error())
		return *common.Error(-1, "从zimg下载文件失败:"+err.Error())
	}

	imgType, err := imagedecode.GetImageType(localFileName)
	if err != nil {
		logger.Error("不是图片类型:" + err.Error())
		return *common.Error(-1, "不是图片类型:"+err.Error())
	}
	ext := ""
	switch imgType {
	case imagedecode.GIF_TYPE:
		ext = "gif"
	case imagedecode.JPG_TYPE:
		ext = "jpg"
	case imagedecode.PNG_TYPE:
		ext = "png"
	}

	//删除本地缓存文件
	err = os.Remove(localFileName)
	if err != nil {
		logger.Error("删除本地缓存文件失败:" + err.Error())
		return *common.Error(-1, "删除本地缓存文件失败:"+err.Error())
	}
	result := make(map[string]string)

	//保存到MySQL
	imageFile := dao.GetImage(fileId)
	if imageFile != nil {
		//已存在，直接返回
		result["imgUrl"] = zimg.ZIMG_HOST + imageFile.Id
		return *common.Success(result)
	}
	imageFile = new(pojo.ImageFiles)
	imageFile.Id = fileId
	imageFile.SourceUrl = localFileName
	imageFile.Size = size
	imageFile.ImgType = ext
	imageFile.Remark = remark
	imageFile.ServerPath = path
	imageFile.FileName = fileName
	imageFile, err = dao.SaveImage(imageFile)
	if err != nil {
		logger.Error("图片信息入库失败:" + err.Error())
		return *common.Error(-1, "图片信息入库失败:"+err.Error())
	}
	result["imgUrl"] = zimg.ZIMG_HOST + fileId
	return *common.Success(result)
}
