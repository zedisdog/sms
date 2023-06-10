package qiniu

import (
	"encoding/json"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/qiniu/go-sdk/v7/auth"
	qiniusms "github.com/qiniu/go-sdk/v7/sms"
	"github.com/zedisdog/sms"
)

var _ sms.IDriver = (*Driver)(nil)

func NewDriver(accessKey string, secretKey string) *Driver {
	d := &Driver{}
	d.client = qiniusms.NewManager(&auth.Credentials{
		AccessKey: accessKey,
		SecretKey: []byte(secretKey),
	})
	return d
}

type Driver struct {
	client *qiniusms.Manager
}

// Send implements sms.IDriver.
func (d *Driver) Send(request sms.Request) (resp sms.Resposne, err error) {
	req := qiniusms.MessagesRequest{
		SignatureID: request.SignName,
		TemplateID:  request.TemplateCode,
		Mobiles:     request.Mobiles,
	}

	if request.Content != nil {
		req.Parameters, err = d.convertContent(request.Content)
		if err != nil {
			return
		}
	}

	r, err := d.client.SendMessage(req)
	if err != nil {
		return
	}
	resp.Raw = r.JobID
	return
}

func (d *Driver) convertContent(content any) (params map[string]any, err error) {
	switch c := content.(type) {
	case string:
		err = json.Unmarshal([]byte(c), &params)
		if err != nil {
			return
		}
	default:
		params = gconv.Map(c)
	}
	return
}
