package service

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"path/filepath"
	"ququ.im/common"
	"ququ.im/common/config"
	"ququ.im/common/utils"
	"ququ.im/syb/common/nacos"
	pojo2 "ququ.im/syb/common/pojo"
	"ququ.im/zimg-go/pojo"
	"ququ.im/zimg-go/zimg"
)

func UploadQualificationToZimg(accountId string) common.Result {
	if accountId == "" {
		return *common.Error(-1, "缺少生意宝账户id")
	}
	qualification := nacos.GetQualificationUrls(accountId)
	if qualification == nil {
		return *common.Error(-1, "此账户无资质数据")
	}
	path := "/qualification/" + accountId
	remark := new(pojo.QualificationRemark)
	remark.AccountId = accountId
	//身份证正面
	if qualification.Idcardfront != "" {
		remark.PicType = "idCardFront"
		id, err := SaveImageFromUrl(qualification.Idcardfront, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Idcardfront = zimg.ZIMG_HOST + id
		}
	}

	//身份证反面
	if qualification.Idcardback != "" {
		remark.PicType = "idCardBack"
		id, err := SaveImageFromUrl(qualification.Idcardback, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Idcardback = zimg.ZIMG_HOST + id
		}
	}

	//银行卡正面
	bankCardNo := ""
	if qualification.Bankcardfront != nil || len(qualification.Bankcardfront) != 0 {
		remark.PicType = "bankCardFront"
		for k, v := range qualification.Bankcardfront {
			remark.BankCardNo = k
			id, err := SaveImageFromUrl(v, path, utils.ToJSON(remark))
			if err != nil {
				if err.Error() != "源图地址不能是zimg地址" {
					return *common.Error(-1, "上传身份证正面错误:"+err.Error())
				}
			} else {
				qualification.Bankcardfront[k] = zimg.ZIMG_HOST + id
			}
			bankCardNo = k
		}
	}
	remark.BankCardNo = ""

	//手持身份证
	if qualification.Personidcard != "" {
		remark.PicType = "personIdCard"
		id, err := SaveImageFromUrl(qualification.Personidcard, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Personidcard = zimg.ZIMG_HOST + id
		}
	}

	//营业执照
	if qualification.Companylicense != "" {
		remark.PicType = "companyLicense"
		id, err := SaveImageFromUrl(qualification.Companylicense, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Companylicense = zimg.ZIMG_HOST + id
		}
	}

	//店头照
	if qualification.Shopface != "" {
		remark.PicType = "shopFace"
		id, err := SaveImageFromUrl(qualification.Shopface, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Shopface = zimg.ZIMG_HOST + id
		}
	}

	//内景照
	if qualification.Shoproom != "" {
		remark.PicType = "shopRoom"
		id, err := SaveImageFromUrl(qualification.Shoproom, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Shoproom = zimg.ZIMG_HOST + id
		}
	}

	//法人身份证正面
	if qualification.Legalpersonidcardfront != "" {
		remark.PicType = "legalPersonIdCardFront"
		id, err := SaveImageFromUrl(qualification.Legalpersonidcardfront, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Legalpersonidcardfront = zimg.ZIMG_HOST + id
		}
	}

	//法人身份证反面
	if qualification.Legalpersonidcardback != "" {
		remark.PicType = "legalPersonIdCardBack"
		id, err := SaveImageFromUrl(qualification.Legalpersonidcardback, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Legalpersonidcardback = zimg.ZIMG_HOST + id
		}
	}

	//手持银行卡
	if qualification.Personbankcard != "" {
		remark.PicType = "personBankCard"
		id, err := SaveImageFromUrl(qualification.Personbankcard, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Personbankcard = zimg.ZIMG_HOST + id
		}
	}

	//法人手持银行卡
	if qualification.Legalpersonbankcard != "" {
		remark.PicType = "legalPersonBankCard"
		id, err := SaveImageFromUrl(qualification.Legalpersonbankcard, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Legalpersonbankcard = zimg.ZIMG_HOST + id
		}
	}

	//税务登记证照片
	if qualification.Companytax != "" {
		remark.PicType = "companyTax"
		id, err := SaveImageFromUrl(qualification.Companytax, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Companytax = zimg.ZIMG_HOST + id
		}
	}

	//组织机构代码证
	if qualification.Companyorg != "" {
		remark.PicType = "companyOrg"
		id, err := SaveImageFromUrl(qualification.Companyorg, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Companyorg = zimg.ZIMG_HOST + id
		}
	}

	//银行开户许可证
	if qualification.Companybank != "" {
		remark.PicType = "companyBank"
		id, err := SaveImageFromUrl(qualification.Companybank, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Companybank = zimg.ZIMG_HOST + id
		}
	}

	//行业特许证
	if qualification.Tradelicense != "" {
		remark.PicType = "tradeLicense"
		id, err := SaveImageFromUrl(qualification.Tradelicense, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Tradelicense = zimg.ZIMG_HOST + id
		}
	}

	//收银台照片
	if qualification.Cashier != "" {
		remark.PicType = "cashier"
		id, err := SaveImageFromUrl(qualification.Cashier, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Cashier = zimg.ZIMG_HOST + id
		}
	}

	//手持营业执照门头照合影
	if qualification.Legalpersonlicenseshopface != "" {
		remark.PicType = "legalPersonLicenseShopFace"
		id, err := SaveImageFromUrl(qualification.Legalpersonlicenseshopface, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Legalpersonlicenseshopface = zimg.ZIMG_HOST + id
		}
	}

	//授权书
	if qualification.Authorization != "" {
		remark.PicType = "authorization"
		id, err := SaveImageFromUrl(qualification.Authorization, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Authorization = zimg.ZIMG_HOST + id
		}
	}

	//商户风险承诺书
	if qualification.Pledge != "" {
		remark.PicType = "pledge"
		id, err := SaveImageFromUrl(qualification.Pledge, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Pledge = zimg.ZIMG_HOST + id
		}
	}

	//法人个人签名照片
	if qualification.Personsign != "" {
		remark.PicType = "personSign"
		id, err := SaveImageFromUrl(qualification.Personsign, path, utils.ToJSON(remark))
		if err != nil {
			if err.Error() != "源图地址不能是zimg地址" {
				return *common.Error(-1, "上传身份证正面错误:"+err.Error())
			}
		} else {
			qualification.Personsign = zimg.ZIMG_HOST + id
		}
	}

	nacos.SaveQualificationUrls(*qualification, bankCardNo)
	return *common.Success(qualification)
}

func UploadQualificationToZimgAll(startAccountId string) common.Result {
	var qualifications []pojo2.QualificationUrls
	count, _ := config.Mgo.C("QualificationUrls").Count()
	errorAccountIds := []string{}
	size := 100
	lastPage := count/size + 1
	up := true
	if startAccountId != "" {
		up = false
	}
	for page := 1; page <= lastPage; page++ {
		config.Mgo.C("QualificationUrls").Find(bson.M{}).Sort("accountId").Skip((page - 1) * size).Limit(size).All(&qualifications)
		for _, qualification := range qualifications {
			fmt.Println("========================================================================================================")
			logger.Debug("正在上传" + qualification.Accountid + "到zimg服务器")
			if !up && qualification.Accountid != startAccountId {
				logger.Debug("跳过当前账户")
				continue
			}
			up = true
			result := UploadQualificationToZimg(qualification.Accountid)
			if result.Status != 1 {
				logger.Error(qualification.Accountid + "资质上传错误:" + result.Msg + "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				errorAccountIds = append(errorAccountIds, qualification.Accountid)
			}
		}
	}
	logger.Debug("全部资质上传到zimg完成！")
	if len(errorAccountIds) > 0 {
		logger.Error("发现资质上传失败的账户列表:\n" + utils.ToJSON(errorAccountIds))
	}
	return *common.Success(errorAccountIds)
}

func UploadQualificationImage(accountId, picType, bankCardNo, localFileName string) common.Result {
	path := "/qualification/" + accountId
	remark := new(pojo.QualificationRemark)
	remark.AccountId = accountId
	remark.PicType = picType
	remark.BankCardNo = bankCardNo
	result := UploadImageFile(localFileName, path, utils.ToJSON(remark), filepath.Base(localFileName))
	if result.Status != 1 {
		return result
	}
	imageFile := result.Data.(*pojo.ImageFiles)
	params := make(map[string]string)
	params["accountId"] = accountId
	params[picType] = zimg.ZIMG_HOST + imageFile.Id
	params["bankCardNo"] = bankCardNo
	r, err := utils.CallNacos(nacos.QUALIFICATION_SERVICE, nacos.SAVE_QUALIFICATION, params)
	if err != nil {
		logger.Error("qualification微服务访问异常:" + err.Error())
		return *common.Error(-1, "资质微服务访问异常:"+err.Error())
	}
	utils.FromJSON(r, &result)
	return result
}
