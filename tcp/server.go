package tcp

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type ServerData struct {
	server  ziface.IServer
	port    int
	onRecv  func([]byte)
	clients []uint64
}

func (this *ServerData) SetOnConnStart(f func(ziface.IConnection)) {
	this.server.SetOnConnStart(f)
}

func (this *ServerData) SetOnConnStop(f func(ziface.IConnection)) {
	this.server.SetOnConnStop(f)
}

func (this *ServerData) Send(bytes []byte) {
	if this.clients == nil || len(this.clients) == 0 {
		for _, index := range this.server.GetConnMgr().GetAllConnID() {
			conn, err := this.server.GetConnMgr().Get(index)
			if err == nil {
				conn.Send(bytes)
			}
		}
	} else {
		for _, id := range this.clients {
			conn, err := this.server.GetConnMgr().Get(id)
			if err == nil {
				conn.Send(bytes)
			}
		}
	}

}

func NewServer(port int) IBaseTcp {
	this := &ServerData{
		port: port,
		server: znet.NewServer(func(s *znet.Server) {
			s.Port = port
		}),
	}
	this.server.SetDecoder(nil)
	this.server.AddInterceptor(this)

	return this
}
func (this *ServerData) SetClientArray(uint64s []uint64) {
	this.clients = uint64s
}

func (this *ServerData) SetRecv(recv func([]byte)) {
	this.onRecv = recv
}
func (this *ServerData) SetDecoder(lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int, maxFrameLength uint64, order binary.ByteOrder) {
	if this.server == nil {
		return
	}
	if lengthFieldLength > 0 {
		field := ziface.LengthField{
			LengthFieldOffset:   lengthFieldOffset,
			LengthFieldLength:   lengthFieldLength,
			LengthAdjustment:    lengthAdjustment,
			InitialBytesToStrip: initialBytesToStrip,
			MaxFrameLength:      maxFrameLength,
			Order:               order,
		}
		this.server.SetDecoder(NewDecoder(&field))
	} else {
		this.server.SetDecoder(nil)
	}
}

func (this *ServerData) Run() {
	this.server.Serve()
}

func (this *ServerData) Stop() {
	this.server.Stop()
}

func (this *ServerData) Start() {
	this.server.Start()
}

func (this *ServerData) GetMaps() map[string]uint64 {
	stringMap := make(map[string]uint64, 10)
	for _, index := range this.server.GetConnMgr().GetAllConnID() {
		conn, err := this.server.GetConnMgr().Get(index)
		if err == nil {
			stringMap[conn.GetConnection().RemoteAddr().String()] = index
		}
	}
	return stringMap
}
func (this *ServerData) Intercept(chain ziface.Chain) ziface.Response {
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

func (this *ServerData) IsServer() bool {
	return true
}
