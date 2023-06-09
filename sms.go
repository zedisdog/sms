package sms

type Request struct {
	Mobiles         []string //手机号,支持多个,阿里云
	SignName        string   //签名,阿里云
	TemplateCode    string   //模板代码,阿里云
	Content         any      //内容 map or string   //模板参数,阿里云
	SmsUpExtendCode string   //上行短信扩展码,阿里云
	OutId           string   //外部流水id,阿里云
}

type Resposne struct {
	Raw string //原始数据
}

type IDriver interface {
	Send(request Request) (Resposne, error)
}

type Sms struct {
	drivers map[string]IDriver
}

func (s *Sms) Send(request Request, driver ...string) (response Resposne, err error) {
	var drv IDriver
	if len(driver) > 0 {
		d := driver[0]
		drv = s.drivers[d]
	} else {
		for _, d := range s.drivers {
			drv = d
			break
		}
	}

	return drv.Send(request)
}
