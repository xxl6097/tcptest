package tcp

import (
	"encoding/binary"
	"github.com/aceld/zinx/ziface"
)

type IBaseTcp interface {
	SetOnConnStart(func(ziface.IConnection)) //设置该Server的连接创建时Hook函数
	SetOnConnStop(func(ziface.IConnection))  //设置该Server的连接断开时的Hook函数
	SetDecoder(lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int, maxFrameLength uint64, order binary.ByteOrder)
	Run()
	Start()
	Stop()
	SetRecv(recv func([]byte))
	Send([]byte)
	GetMaps() map[string]uint64
	SetClientArray([]uint64)
	IsServer() bool
}
