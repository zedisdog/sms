package aliyun

import (
	"encoding/json"
	"errors"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zedisdog/sms"
)

func NewDriver(accessKeyID string, accessKeySecret string, RegionID string) (driver sms.IDriver, err error) {
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(RegionID),
	})
	if err != nil {
		return
	}

	driver = &Driver{
		client: client,
	}

	return
}

var _ sms.IDriver = (*Driver)(nil)

type Driver struct {
	client *dysmsapi.Client
}

// Send implements sms.IDriver.
func (d *Driver) Send(request sms.Request) (response sms.Resposne, err error) {
	err = d.CheckRequest(request)
	if err != nil {
		return
	}

	req := dysmsapi.SendSmsRequest{
		PhoneNumbers: tea.String(strings.Join(request.Mobiles, ",")),
		SignName:     tea.String(request.SignName),
		TemplateCode: tea.String(request.TemplateCode),
	}

	params, err := d.converContent(request.Content)
	if err != nil {
		return
	}

	if params != "" {
		req.TemplateParam = tea.String(params)
	}

	resp, err := d.client.SendSms(&req)

	response = sms.Resposne{
		Raw: resp.GoString(),
	}

	return
}

func (d *Driver) converContent(content any) (result string, err error) {
	switch c := content.(type) {
	case nil:
		return "", nil
	case string:
		return c, nil
	case map[string]any:
		var b []byte
		b, err = json.Marshal(c)
		if err != nil {
			return
		}
		result = string(b)
	default:
		err = errors.New("not support type")
	}

	return
}

func (d *Driver) CheckRequest(request sms.Request) (err error) {
	if len(request.Mobiles) <= 0 {
		return errors.New("Mobile is required")
	}

	if request.SignName == "" {
		return errors.New("SignName is required")
	}

	if request.TemplateCode == "" {
		return errors.New("TemplateCode is required")
	}

	return
}
