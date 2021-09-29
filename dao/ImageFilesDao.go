package dao

import (
	"ququ.im/common/config"
	"ququ.im/common/logs"
	"ququ.im/zimg-go/pojo"
	"strings"
	"time"
)

const TABLE_IMAGE_FILES string = "image_files"

func SaveImage(imageFiles *pojo.ImageFiles) (*pojo.ImageFiles, error) {
	imageFiles.UploadTime = time.Now().Format("2006-01-02 15:04:05")
	tx := config.Mysql.Begin()
	err := tx.Table(TABLE_IMAGE_FILES).Create(imageFiles).Error
	if err != nil {
		logs.Error("插入数据错误:{}", err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return imageFiles, nil
}

func UpdateImageById(id string, updateFields map[string]interface{}) (*pojo.ImageFiles, error) {
	tx := config.Mysql.Begin()
	err := tx.Table(TABLE_IMAGE_FILES).Where("id = ?", id).Update(updateFields).Error
	if err != nil {
		logs.Error("更新数据错误:{}", err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return GetImage(id), nil
}

func DeleteImage(id string) error {
	tx := config.Mysql.Begin()
	err := tx.Table(TABLE_IMAGE_FILES).Delete(nil, "id = ?", id).Error
	if err != nil {
		logs.Error("更新删除错误:{}", err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetImage(id string) *pojo.ImageFiles {
	var imageFile pojo.ImageFiles
	config.Mysql.Table(TABLE_IMAGE_FILES).Where("id = ?", id).First(&imageFile)
	if imageFile.Id == "" {
		return nil
	} else {
		return &imageFile
	}
}

func ListImageByUrl(sourceUrl string) []pojo.ImageFiles {
	var imageFiles []pojo.ImageFiles
	config.Mysql.Table(TABLE_IMAGE_FILES).Where("source_url = ?", sourceUrl).Order("upload_time", true).Find(&imageFiles)
	if len(imageFiles) == 0 {
		return nil
	} else {
		return imageFiles
	}
}

func ListImageByPath(path string) []pojo.ImageFiles {
	var imageFiles []pojo.ImageFiles
	if strings.Contains(path, "*") {
		path = strings.ReplaceAll(path, "*", "%")
		config.Mysql.Table(TABLE_IMAGE_FILES).Where("server_path like ?", path).Order("upload_time", true).Find(&imageFiles)
	} else {
		config.Mysql.Table(TABLE_IMAGE_FILES).Where("server_path = ?", path).Order("upload_time", true).Find(&imageFiles)
	}
	if len(imageFiles) == 0 {
		return nil
	} else {
		return imageFiles
	}
}

func ListImageByRemark(remark string) []pojo.ImageFiles {
	var imageFiles []pojo.ImageFiles
	if strings.Contains(remark, "*") {
		remark = strings.ReplaceAll(remark, "*", "%")
		config.Mysql.Table(TABLE_IMAGE_FILES).Where("remark like ?", remark).Order("upload_time", true).Find(&imageFiles)
	} else {
		config.Mysql.Table(TABLE_IMAGE_FILES).Where("remark = ?", remark).Order("upload_time", true).Find(&imageFiles)
	}
	if len(imageFiles) == 0 {
		return nil
	} else {
		return imageFiles
	}
}

func ListSourceUrlByUploadTime(beginDate, endDate string) []string {
	var sourceUrls []string
	var imageFiles []pojo.ImageFiles
	config.Mysql.Table(TABLE_IMAGE_FILES).Where("upload_time between ? and ? and source_url like ?", beginDate, endDate, "http%").Scan(&imageFiles)
	if len(imageFiles) == 0 {
		return nil
	} else {
		for _, imageFile := range imageFiles {
			sourceUrls = append(sourceUrls, imageFile.SourceUrl)
		}
		return sourceUrls
	}
}
