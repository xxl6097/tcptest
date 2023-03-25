package entity

type LengthFieldEntity struct {
	LengthFieldOffset   int
	LengthFieldLength   int
	LengthAdjustment    int
	InitialBytesToStrip int
	MaxFrameLength      uint64
	Order               string
}

type WorkEntity struct {
	Workmode   string
	Localport  int
	Destip     string
	Destport   int
	Stickymode bool
}

type RecvEntity struct {
	Hex bool
}

type SendEntity struct {
	Hex        bool
	Stickymode bool
	Sendmode   bool
	Interval   uint64
}

type Entity struct {
	Work        WorkEntity
	Lengthfield LengthFieldEntity
	Send        SendEntity
	Recv        RecvEntity
	IsCache     bool
}
