package pojo

type QualificationRemark struct {
	AccountId  string `json:"accountId"`
	PicType    string `json:"picType"`
	BankCardNo string `json:"bankCardNo,omitempty"`
}
