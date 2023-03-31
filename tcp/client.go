package tcp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type ClientData struct {
	client  ziface.IClient
	port    int
	address string
	onRecv  func([]byte)
}

func (this *ClientData) IsServer() bool {
	return false
}

func (this *ClientData) GetMaps() map[string]uint64 {
	return nil
}

func NewClient(address string, port int) IBaseTcp {
	this := &ClientData{
		address: address,
		port:    port,
		client:  znet.NewClient(address, port),
	}
	this.client.SetDecoder(nil)
	this.client.AddInterceptor(this)

	return this
}

func (this *ClientData) Run() {
}

func (this *ClientData) SetOnConnStart(f func(ziface.IConnection)) {
	this.client.SetOnConnStart(f)
}

func (this *ClientData) SetOnConnStop(f func(ziface.IConnection)) {
	this.client.SetOnConnStop(f)
}
func (this *ClientData) SetClientArray(uint64s []uint64) {
}
func (this *ClientData) SetDecoder(lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int, maxFrameLength uint64, order binary.ByteOrder) {
	if this.client == nil {
		return
	}
	field := ziface.LengthField{
		LengthFieldOffset:   lengthFieldOffset,
		LengthFieldLength:   lengthFieldLength,
		LengthAdjustment:    lengthAdjustment,
		InitialBytesToStrip: initialBytesToStrip,
		MaxFrameLength:      maxFrameLength,
		Order:               order,
	}
	this.client.SetDecoder(NewDecoder(&field))
}

func (this *ClientData) Start() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("second recover", r)
		}
	}()
	this.client.Start()
	fmt.Println("...start...")
}

func (this *ClientData) Stop() {
	this.client.Stop()
}

func (this *ClientData) SetRecv(recv func([]byte)) {
	this.onRecv = recv
}

func (this *ClientData) Send(bytes []byte) {
	this.client.Conn().Send(bytes)
}
func (this *ClientData) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			iMessage := iRequest.GetMessage()
			if iMessage != nil {
				data := iMessage.GetData()
				zlog.Ins().DebugF("RawData size:%d data:%s %s\n", len(data), hex.EncodeToString(data), string(data))
				this.onRecv(data)
			}
		}
	}
	return chain.Proceed(chain.Request())
}
