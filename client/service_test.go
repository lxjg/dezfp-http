package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_encrypt()
func Test_MakeOut(t *testing.T) {

	data := &Data{
		Content: &RequestContent{
			Fpqqlsh: "111MFWIK201805081541348952",
			Dsptbm:  "111MFWIK",
			Nsrsbh:  "310101000000090",
			Ddh:     "201805081541348952",
			PdfXzfs: "2",
		},
	}

	client := NewBillClient()
	client.Global = &GlobalInfo{
		TerminalCode:      "0",
		AppID:             "ZZS_PT_DZFP",
		Version:           "2.0",
		InterfaceCode:     "ECXML.FPXZ.CX.E_INV",
		UserName:          "111MFWIK",
		PassWord:          "12345678909oyKs7cVo1yYzkuisP9bhA==",
		TaxpayerID:        "310101000000090",
		AuthorizationCode: "3100000090",
		RequestCode:       "111MFWIK",
		RequestTime:       "2016-11-14 14:11:30 301", //time.Now().Format("2006-01-02 15:04:05 700"),
		ResponseCode:      "121",
		DataExchangeID:    "111MFWIK20180508154134634",
	}
	client.RequestData = data
	client.ReturnState = &ReturnStateInfo{}
	client.Key = "9oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.MakeOut()

	if assert.Nil(t, err) {
		if data, ok := result.(*ResponseData); ok {
			assert.NotNil(t, data.Content)
		}

		if state, ok := result.(*ReturnStateInfo); ok {
			if !assert.Equal(t, "0000", state.ReturnCode) {
				panic(state.ReturnMessage)
			}
		}
	}
}

func Test_Download(t *testing.T) {
	data := &Data{
		Content: &RequestContent{
			Fpqqlsh: "111MFWIK201805081541348952",
			Dsptbm:  "111MFWIK",
			Nsrsbh:  "310101000000090",
			Ddh:     "201805081541348952",
			PdfXzfs: "2",
		},
	}

	client := NewBillClient()
	client.Global = &GlobalInfo{
		TerminalCode:      "0",
		AppID:             "ZZS_PT_DZFP",
		Version:           "2.0",
		InterfaceCode:     "ECXML.FPXZ.CX.E_INV",
		UserName:          "111MFWIK",
		PassWord:          "12345678909oyKs7cVo1yYzkuisP9bhA==",
		TaxpayerID:        "310101000000090",
		AuthorizationCode: "3100000090",
		RequestCode:       "111MFWIK",
		RequestTime:       time.Now().Format("2006-01-02 15:04:05 700"),
		ResponseCode:      "121",
		DataExchangeID:    "111MFWIK20180508154134634",
	}
	client.RequestData = data
	client.ReturnState = &ReturnStateInfo{}
	client.Key = "9oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.Download()

	if assert.Nil(t, err) {
		if data, ok := result.(*ResponseData); ok {
			assert.NotNil(t, data.Content)
		}

		if state, ok := result.(*ReturnStateInfo); ok {
			if !assert.Equal(t, "0000", state.ReturnCode) {
				panic(state.ReturnMessage)
			}
		}
	}
}
