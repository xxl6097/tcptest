package tcp

import (
	"encoding/binary"
	"fmt"
	"tcptest/entity"
)

type TcpTest interface {
	Init(_data *entity.Entity)
	GetTcp() IBaseTcp
}

type Tcp struct {
	baseTcp IBaseTcp
}

func NewTcpTest() TcpTest {
	tcp := &Tcp{}
	return tcp
}

func (this *Tcp) Init(_data *entity.Entity) {
	if _data == nil {
		return
	}
	if _data.Work.Workmode == "" {
		return
	}
	if _data.Work.Workmode == "TCP客户端" {
		this.baseTcp = NewClient(_data.Work.Destip, _data.Work.Destport)
	} else {
		this.baseTcp = NewServer(_data.Work.Localport)
		if _data.Work.Stickymode {
			var order binary.ByteOrder
			if _data.Lengthfield.Order == "小端模式" {
				order = binary.LittleEndian
			} else {
				order = binary.BigEndian
			}
			this.baseTcp.SetDecoder(
				_data.Lengthfield.LengthFieldOffset,
				_data.Lengthfield.LengthFieldLength,
				_data.Lengthfield.LengthAdjustment,
				_data.Lengthfield.InitialBytesToStrip,
				_data.Lengthfield.MaxFrameLength,
				order)
		}
	}
	fmt.Println("... start")
}

func (this *Tcp) GetTcp() IBaseTcp {
	return this.baseTcp
}
