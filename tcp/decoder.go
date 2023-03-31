package tcp

import (
	"encoding/hex"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

type Data struct {
	field *ziface.LengthField
}

type Decoder interface {
	ziface.IDecoder
}

func NewDecoder(_field *ziface.LengthField) ziface.IDecoder {
	return &Data{
		field: _field,
	}
}

func (this *Data) GetLengthField() *ziface.LengthField {
	return this.field
}

func (this *Data) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			iMessage := iRequest.GetMessage()
			if iMessage != nil {
				data := iMessage.GetData()
				zlog.Ins().DebugF("HTLVCRC-RawData size:%d data:%s\n", len(data), hex.EncodeToString(data))
			}
		}
	}
	return chain.Proceed(chain.Request())
}
