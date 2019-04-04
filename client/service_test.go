package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MakeOut(t *testing.T) {
	client := NewBillClient("111MFWIK", "310101000000090", "3100000090", "92884519", "9oyKs7cVo1yYzkuisP9bhA==")
	client.Global.PassWord = "12345678909oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.MakeOut("上海爱信诺航天信息有限公司", "上海市浦东新区", "13810189561", "10.0", time.Now())

	if assert.Nil(t, err) {
		if data, ok := result.(*ReturnStateInfo); ok {
			assert.Equal(t, "0000", data.ReturnCode)
		}
	}
}

func Test_Download(t *testing.T) {
	client := NewBillClient("111MFWIK", "310101000000090", "3100000090", "92884519", "9oyKs7cVo1yYzkuisP9bhA==")
	client.Global.PassWord = "12345678909oyKs7cVo1yYzkuisP9bhA=="
	result, err := client.Download("111MFWIK201904041554354107")
	if assert.Nil(t, err) {
		if data, ok := result.(*ResponseContent); ok {
			assert.NotNil(t, data.Fpqqlsh)
		}
	}
}
