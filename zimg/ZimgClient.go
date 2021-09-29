package zimg

import (
	"bytes"
	"errors"
	"gitee.com/vchuzu/imagedecode"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty"
	"github.com/levigross/grequests"
	"io/ioutil"
	"os"
	"path/filepath"
	"ququ.im/common/logs"
	"strconv"
	"strings"
)

func Upload(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("文件打开错误:{}", err.Error())
		return "", err
	}
	imgType, err := imagedecode.GetImageType(fileName)
	if err != nil {
		logs.Error("不是图片类型:{}", err.Error())
		return "", errors.New("不是图片类型:" + err.Error())
	}
	ext := ""
	switch imgType {
	case imagedecode.GIF_TYPE:
		ext = ".gif"
	case imagedecode.JPG_TYPE:
		ext = ".jpg"
	case imagedecode.PNG_TYPE:
		ext = ".png"
	}
	fn := filepath.Base(fileName)
	if filepath.Ext(fn) == "" {
		fn = fn + ext
	}
	client := resty.New()
	resp, err := client.R().
		SetMultipartField("userfile", fn, imgType, bytes.NewReader(file)).
		Post(ZIMG_HOST_LAN + API_URL_UPLOAD)
	if err != nil {
		logs.Error("文件上传失败:{}", err.Error())
		return "", err
	}
	if resp.StatusCode() == 200 {
		//logs.Debug("上传文件返回结果: " + resp.String())
		if strings.Contains(resp.String(), "Failed") {
			return "", errors.New("文件上传失败:" + resp.String())
		}
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		a := doc.Find("a").First()
		if a == nil {
			return "", errors.New("文件上传失败:" + resp.String())
		}
		fileId, _ := a.Attr("href")
		fileId = strings.ReplaceAll(fileId, "/", "")
		return fileId, nil
	} else {
		return "", errors.New("文件上传失败:" + resp.String())
	}
}

func Download(fileId, fileName string, w, h, p, x, y, g, q int) (int, string, error) {
	if fileName[0:1] != "/" {
		path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fileName = path + IMAGE_PATH + fileName
	}
	logs.Debug("下载文件保存到:{}", fileName)
	url := ZIMG_HOST + fileId
	params := make(map[string]string)
	if w != 0 {
		params["w"] = strconv.Itoa(w)
	}
	if h != 0 {
		params["h"] = strconv.Itoa(h)
	}
	if p != -1 {
		params["p"] = strconv.Itoa(p)
	}
	if x != 0 {
		params["x"] = strconv.Itoa(x)
	}
	if y != 0 {
		params["y"] = strconv.Itoa(y)
	}
	if g != 0 {
		params["g"] = strconv.Itoa(g)
	}
	if q != 0 {
		params["q"] = strconv.Itoa(q)
	}
	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		logs.Error("文件下载失败:{}", err.Error())
		return 0, fileName, err
	}
	if resp.Ok {
		err = resp.DownloadToFile(fileName)
		if err != nil {
			logs.Error("文件下载失败:{}", err.Error())
			return 0, fileName, err
		}
		return getFileSize(fileName), fileName, nil
	} else {
		logs.Error("文件下载失败:{}", resp.String())
		return 0, fileName, errors.New("文件下载失败:" + resp.String())
	}
}

func Delete(fileId string) error {
	url := ZIMG_HOST_LAN + API_URL_ADMIN
	params := make(map[string]string)
	params["md5"] = fileId
	params["t"] = "1"
	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		logs.Error("远程文件删除失败:{}", err.Error())
		return err
	}
	if resp.Ok {
		if strings.Contains(resp.String(), "Successful") {
			return nil
		} else {
			return errors.New(resp.String())
		}
	} else if resp.StatusCode == 404 {
		return errors.New("远程文件不存在")
	} else {
		return errors.New(resp.String())
	}
}

func GetPictureUrl(fileId string, w, h, p, x, y, g, q int) string {
	url := ZIMG_HOST + fileId + "?"
	if w != 0 {
		url = url + "w=" + strconv.Itoa(w) + "&"
	}
	if h != 0 {
		url = url + "h=" + strconv.Itoa(h) + "&"
	}
	if p != -1 {
		url = url + "p=" + strconv.Itoa(p) + "&"
	}
	if x != 0 {
		url = url + "x=" + strconv.Itoa(x) + "&"
	}
	if y != 0 {
		url = url + "y=" + strconv.Itoa(y) + "&"
	}
	if g != 0 {
		url = url + "g=" + strconv.Itoa(g) + "&"
	}
	if q != 0 {
		url = url + "q=" + strconv.Itoa(q) + "&"
	}
	url = strings.TrimRight(url, "&")
	url = strings.TrimRight(url, "?")
	return url
}

func getFileSize(filename string) int {
	var result int
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = int(f.Size())
		return nil
	})
	return result
}
