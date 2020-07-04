package pojo

type ImageFiles struct {
	Id         string `json:"id" gorm:"colume:id,primary_key"`       //MD5键值
	ImgType    string `json:"img_type" gorm:"column:img_type"`       //图片类型 jpg/png/gif
	Size       int    `json:"size" gorm:"column:size"`               //默认75%质量的文件大小
	ServerPath string `json:"server_path" gorm:"server_path"`        //服务器上路径
	FileName   string `json:"file_name" gorm:"column:file_name"`     //本地未上传时文件名
	UploadTime string `json:"upload_time" gorm:"column:upload_time"` //上传时间
	SourceUrl  string `json:"source_url" gorm:"column:source_url"`   //来源URL，如七牛URL、OSS的URL
	Remark     string `json:"remark" gorm:"column:remark"`           //文件备注，建议json格式
}
