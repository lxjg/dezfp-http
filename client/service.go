package client

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"dezfp-http/tools"
)

// URL 电子发票路由
const (
	// URL 请求地址
	URL = "http://fw1test.shdzfp.com:9000/sajt-shdzfp-sl-http/SvrServlet"

	// 接口编码
	// MakeOutBill 开具发票
	MakeOutBill = "ECXML.FPKJ.BC.E_INV"
	// DownloadBill 下载发票
	DownloadBill = "ECXML.FPXZ.CX.E_INV"
)

// GlobalInfo 全局数据项
type GlobalInfo struct {
	TerminalCode  string `xml:"terminalCode"`
	AppID         string `xml:"appId"`
	Version       string `xml:"version"`
	InterfaceCode string `xml:"interfaceCode"`

	UserName          string `xml:"userName"`
	PassWord          string `xml:"passWord"`
	TaxpayerID        string `xml:"taxpayerId"`
	AuthorizationCode string `xml:"authorizationCode"`
	RequestCode       string `xml:"requestCode"`
	RequestTime       string `xml:"requestTime"`
	ResponseCode      string `xml:"responseCode"`
	DataExchangeID    string `xml:"dataExchangeId"`
}

// ReturnStateInfo 数据交换请求返回状态信息
type ReturnStateInfo struct {
	ReturnCode    string `xml:"returnCode"`
	ReturnMessage string `xml:"returnMessage"`
}

// DataDescription 交换数据属性描述
type DataDescription struct {
	ZipCode     string `xml:"zipCode"`
	EncryptCode string `xml:"encryptCode"`
	CodeType    string `xml:"codeType"`
}

// FPTXX 订单信息
type FPTXX struct {
	// 必填字段
	FPQQLSH   string `xml:"FPQQLSH"`    // FPQQLSH 发票请求唯一流水号
	DSPTBM    string `xml:"DSPTBM"`     // 平台编码
	NSRSBH    string `xml:"NSRSBH"`     // 开票方识别号
	NSRMC     string `xml:"NSRMC"`      // 开票方名称
	DKBZ      string `xml:"DKBZ"`       // 代开标志
	KPXM      string `xml:"KPXM"`       // 主要开票商品，或者第一条商品
	BMBBBH    string `xml:"BMB_BBH"`    // 编码表版本号
	XHFNSRSBH string `xml:"XHF_NSRSBH"` // 销方识别号
	XHFMC     string `xml:"XHFMC"`      // 销方名称
	XHFDZ     string `xml:"XHF_DZ"`     // 销方地址
	XHFDH     string `xml:"XHF_DH"`     // 销方电话
	GHFMC     string `xml:"GHFMC"`      // 购货方名称
	GHFQYLX   string `xml:"GHFQYLX"`    // 购货方企业类型
	KPY       string `xml:"KPY"`        // 开票员
	KPLX      string `xml:"KPLX"`       // 开票类型
	CZDM      string `xml:"CZDM"`       // 操作代码 10正票正常开具20退货折让红票
	QDBZ      string `xml:"QD_BZ"`      // 清单标志 默认为0
	KPHJJE    string `xml:"KPHJJE"`     // 价税合计金额

	// 非必填
	XHFYHZH  string `xml:"XHF_YHZH"`  // 销方银行、账号
	GHFDZ    string `xml:"GHF_DZ"`    // 购货方地址
	GHFGDDH  string `xml:"GHF_GDDH"`  // 购货方固定电话
	GHFSJ    string `xml:"GHF_SJ"`    // 购货方手机
	GHFEMAIL string `xml:"GHF_EMAIL"` // 购货方邮箱
	GHFYHZH  string `xml:"GHF_YHZH"`  // 购货方银行、账号
	HYDM     string `xml:"HY_DM"`     // 行业代码
	HYMC     string `xml:"HY_MC"`     // 行业名称
	SKY      string `xml:"SKY"`       // 收款员
	FHR      string `xml:"FHR"`       // 复核人
	KPRQ     string `xml:"KPRQ"`      // 开票日期
	YFPDM    string `xml:"YFP_DM"`    // 原发票代码
	YFPHM    string `xml:"YFP_HM"`    // 原发票号码
	QDXMMC   string `xml:"QDXMMC"`    // 清单发票项目名称
	CHYY     string `xml:"CHYY"`      // 冲红原因
	TSCHBZ   string `xml:"TSCHBZ"`    // 特殊冲红标志 0正常冲红(电子发票) 1特殊冲红(冲红纸质等)
}

// XMXX 项目信息
type XMXX struct {
	XMLName xml.Name `xml:"FPKJXX_XMXX"`
	XMMC    string   `xml:"XMMC"`  // 项目名称
	XMDW    string   `xml:"XMDW"`  // 项目单位(比如科、门)
	XMSL    string   `xml:"XMSL"`  // 项目数量
	HSBZ    string   `xml:"HSBZ"`  // 含税标志 。0表示都不含税，1 表示都含税
	SPBM    string   `xml:"SPBM"`  // 商品编码
	XMJE    string   `xml:"XMDJ"`  // 项目金额
	SL      string   `xml:"SL"`    // 税率 0标识免税
	FPHXZ   string   `xml:"FPHXZ"` // 发票行性质 0正常行1折扣行2被折扣行
}

// XMXXS 项目组
type XMXXS struct {
	Items []XMXX
}

// DDXX 订单信息
type DDXX struct {
	DDH    string `xml:"DDH"`    // 订单号
	THDH   string `xml:"THDH"`   // 退单号
	DDDATE string `xml:"DDDATE"` // 订单时间
}

// MakeRequestContent 开具发票交换数据内容
type MakeRequestContent struct {
	XMLName xml.Name `xml:"REQUEST_FPKJXX"`
	FPTXX   FPTXX    `xml:"FPKJXX_FPTXX"`
	XMXXS   XMXXS    `xml:"FPKJXX_XMXXS"`
	DDXX    DDXX     `xml:"FPKJXX_DDXX"`
}

// DownloadRequestContent 下载发票交换数据内容
type DownloadRequestContent struct {
	XMLName xml.Name `xml:"REQUEST_FPXXXZ_NEW"`
	FPQQLSH string   `xml:"FPQQLSH"`
	DSPTBM  string   `xml:"DSPTBM"`
	NSRSBH  string   `xml:"NSRSBH"`
	DDH     string   `xml:"DDH"`
	PDFXZFS string   `xml:"PDF_XZFS"`
}

// ResponseContent 交换数据内容
type ResponseContent struct {
	XMLName       xml.Name `xml:"REQUEST_FPKJXX_FPJGXX_NEW"`
	Fpqqlsh       string   `xml:"FPQQLSH"`
	Ddh           string   `xml:"DDH"`
	Kplsh         string   `xml:"KPLSH"`
	Fwm           string   `xml:"FWM"`
	Ewn           string   `xml:"EWM"`
	FpzlDm        string   `xml:"FPZL_DM"`
	FpDm          string   `xml:"FP_DM"`
	FpHm          string   `xml:"FP_HM"`
	KPRQ          string   `xml:"KPRQ"`
	KPLX          string   `xml:"KPLX"`
	HJBHSJE       string   `xml:"HJBHSJE"`
	KPHJSE        string   `xml:"KPHJSE"`
	PdfFile       string   `xml:"PDF_FILE"`
	PdfURL        string   `xml:"PDF_URL"`
	CZDM          string   `xml:"CZDM"`
	RETURNCODE    string   `xml:"RETURNCODE"`
	RETURNMESSAGE string   `xml:"RETURNMESSAGE"`
}

// RequestData 交换数据
type RequestData struct {
	Description *DataDescription `xml:"dataDescription"`
	// EncryptContent 根据Content加密生成
	EncryptContent string `xml:"content"`
	// Content 交换数据内容明文，必需
	Content interface{} `xml:"-"`
	// ActionName 区别请求action
	ActionName string `xml:"-"`
}

// ResponseData 返回数据
type ResponseData struct {
	Content string `xml:"content"`
}

// BillClient xml请求
type BillClient struct {
	XMLName     xml.Name         `xml:"interface"`
	Global      *GlobalInfo      `xml:"globalInfo"`
	ReturnState *ReturnStateInfo `xml:"returnStateInfo"`
	RequestData *RequestData     `xml:"Data"`
	Key         string           `xml:"-"`
}

// SecBillClient 返回数据接收器
type SecBillClient struct {
	XMLName      xml.Name         `xml:"interface"`
	ReturnState  *ReturnStateInfo `xml:"returnStateInfo"`
	ResponseData *ResponseData    `xml:"Data"`
}

// ToString 转成字符串
func (c *RequestData) encrypt(key []byte) error {
	code, _ := xml.Marshal(c.Content)
	str := string(code)
	if c.ActionName == MakeOutBill {
		requestType := `<REQUEST_FPKJXX class="REQUEST_FPKJXX">`
		str = strings.Replace(str, "<REQUEST_FPKJXX>", requestType, -1)
		fptxx := `<FPKJXX_FPTXX class="FPKJXX_FPTXX">`
		str = strings.Replace(str, "<FPKJXX_FPTXX>", fptxx, -1)
		xmxxs := `<FPKJXX_XMXXS class="FPKJXX_XMXX;" size="1">`
		str = strings.Replace(str, "<FPKJXX_XMXXS>", xmxxs, -1)
		ddxx := `<FPKJXX_DDXX class="FPKJXX_DDXX">`
		str = strings.Replace(str, "<FPKJXX_DDXX>", ddxx, -1)
	}

	if c.ActionName == DownloadBill {
		requestType := `<REQUEST_FPXXXZ_NEW class='REQUEST_FPXXXZ_NEW'>`
		str = strings.Replace(str, "<REQUEST_FPXXXZ_NEW>", requestType, -1)
	}
	res, err := tools.TripleDesECBEncrypt([]byte(str), key)

	if err != nil {
		return err
	}

	c.EncryptContent = base64.StdEncoding.EncodeToString(res)

	return nil
}

func (c *RequestData) defaultDescription() {
	c.Description = &DataDescription{
		ZipCode:     "0",
		EncryptCode: "1",
		CodeType:    "3DES",
	}
}

// ToString 转成字符串
func (s *BillClient) toString() string {
	code, _ := xml.Marshal(s)
	interfactType := `<?xml version="1.0" encoding="utf-8"?>
	<interface xmlns="" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.chinatax.gov.cn/tirip/dataspec/interfaces.xsd"
	version="DZFP1.0">`
	return strings.Replace(string(code), "<interface>", interfactType, -1)
}

// NewBillClient 实例化发票对象
func NewBillClient() *BillClient {
	return &BillClient{
		Global: &GlobalInfo{
			Version: "2.0",
		},
		ReturnState: &ReturnStateInfo{},
	}
}

// MakeOut 开具发票
func (s *BillClient) MakeOut() (interface{}, error) {
	// 发起请求
	return s.doAction(MakeOutBill)
}

// Download 开具发票
func (s *BillClient) Download() (interface{}, error) {
	// 发起请求
	return s.doAction(DownloadBill)
}

func (s *BillClient) init(interfaceCode string) []byte {
	s.setInterfaceCode(interfaceCode)
	if interfaceCode == MakeOutBill {
		s.RequestData.ActionName = MakeOutBill
	}
	if interfaceCode == DownloadBill {
		s.RequestData.ActionName = DownloadBill
	}
	s.RequestData.encrypt([]byte(s.Key))
	if s.RequestData.Description == nil {
		s.RequestData.defaultDescription()
	}

	return []byte(s.toString())
}

func (s *BillClient) getInterfaceCode() string {
	return s.Global.InterfaceCode
}

func (s *BillClient) setInterfaceCode(code string) {
	s.Global.InterfaceCode = code
}

func (s *BillClient) doAction(interfaceCode string) (interface{}, error) {
	// BillClient初始化
	xmlStr := s.init(interfaceCode)
	//发送请求.
	req, err := http.NewRequest("POST", URL, bytes.NewReader(xmlStr))
	if err != nil {
		return nil, err
	}

	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Add("charset", "utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if s.getInterfaceCode() == MakeOutBill {
		var xmlRe BillClient
		err = xml.Unmarshal(body, &xmlRe)
		return xmlRe.ReturnState, err
	}
	var xmlRe SecBillClient
	err = xml.Unmarshal(body, &xmlRe)

	if xmlRe.ReturnState.ReturnCode != "0000" {
		return xmlRe.ReturnState, errors.New(xmlRe.ReturnState.ReturnCode)
	}
	// base64解码
	crypted, err := base64.StdEncoding.DecodeString(xmlRe.ResponseData.Content)
	if err != nil {
		return nil, err
	}
	// 3des解密
	res, err := tools.TripleDesECBDecrypt(crypted, []byte(s.Key))
	if err != nil {
		return nil, err
	}
	// 解析content
	var responseContent ResponseContent
	err = xml.Unmarshal(res, &responseContent)
	return &responseContent, err
}
