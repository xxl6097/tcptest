package view

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net"
	"net/url"
	"strconv"
	"time"
)

type Left struct {
	rightview                                                                                   RightView
	mainview                                                                                    *MainView
	localportEntry, destipEntry, destportEntry, sendEverySecEntry                               *widget.Entry
	workmodeSel, orderSel                                                                       *widget.Select
	startBtn, clearLogBtn, saverecvbtn                                                          *widget.Button
	recvHexCheck, sendHexCheck, sendIntervalCheck, sendStickyCheck, setStickyCheck              *widget.Check
	rootview                                                                                    *fyne.Container
	workmodeview                                                                                *fyne.Container
	lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip, maxFrameLength *widget.Entry
	destiplabel, destportlabel, localportlabel, worklabel                                       *widget.Label
	workhbox, clienthbox, serverhbox, lenthfield                                                *fyne.Container
}

type LeftView interface {
	GetView() fyne.CanvasObject
	OnConnStatus(bool)
}

func NewLeftView(_mainview *MainView, _rightview RightView) LeftView {
	view := &Left{
		mainview:  _mainview,
		rightview: _rightview,
	}
	view.initView()
	view.initData()
	return view
}

func (this *Left) initView() *fyne.Container {
	setview := this.settingView()
	recvview := this.recvView()
	sendview := this.sendView()
	disview := this.displayView()

	label1 := container.NewGridWrap(fyne.NewSize(200, 10), widget.NewLabel("通讯设置"))
	label2 := container.NewGridWrap(fyne.NewSize(200, 10), widget.NewLabel("接收区设置"))
	label3 := container.NewGridWrap(fyne.NewSize(200, 10), widget.NewLabel("发送区设置"))
	this.rootview = container.NewVBox(label1, setview, label2, recvview, label3, sendview, disview)
	return this.rootview
}

func (this *Left) settingView() *fyne.Container {
	this.workmodeSel = widget.NewSelect([]string{"TCP客户端", "TCP服务器"}, func(item string) {
		this.workmodeview.RemoveAll()
		this.workmodeview.Add(this.workhbox)
		if item == "TCP服务器" {
			this.workmodeview.Add(this.serverhbox)
		} else {
			this.workmodeview.Add(this.clienthbox)
		}
		this.workmodeview.Add(this.startBtn)
		this.setStickyCheck.SetChecked(false)
		this.rightview.OnWorkModeSelect(item)
	})
	this.workmodeSel.Selected = AppData.Work.Workmode
	this.localportEntry = widget.NewEntry()
	this.destipEntry = widget.NewEntry()
	this.destportEntry = widget.NewEntry()
	this.destipEntry.SetPlaceHolder("目标IP或域名")
	this.destportEntry.SetPlaceHolder("目标端口")
	this.localportlabel = widget.NewLabel("本地端口：")

	this.worklabel = widget.NewLabel("工作模式：")
	this.destiplabel = widget.NewLabel("目的IP：      ")
	this.destportlabel = widget.NewLabel("目的端口：")

	this.startBtn = widget.NewButton("打开", this.onClick(0))
	//this.startBtn.Importance = widget.WarningImportance

	this.setStickyCheck = widget.NewCheck("断粘包设置", func(b bool) {
		if b {
			this.serverhbox.Add(this.lenthfield)
		} else {
			this.serverhbox.Remove(this.lenthfield)
		}
	})

	labels1 := container.NewVBox(
		widget.NewLabel("lengthFieldOffset："),
		widget.NewLabel("lengthFieldLength："),
		widget.NewLabel("lengthAdjustment："),
		widget.NewLabel("initialBytesToStrip："),
		widget.NewLabel("maxFrameLength："),
		widget.NewLabel("大小端："))

	//https://blog.csdn.net/weixin_45271492/article/details/125347939
	this.lengthFieldOffset = widget.NewEntry()
	this.lengthFieldLength = widget.NewEntry()
	this.lengthAdjustment = widget.NewEntry()
	this.initialBytesToStrip = widget.NewEntry()
	this.maxFrameLength = widget.NewEntry()
	this.lengthFieldOffset.PlaceHolder = "长度字段偏移量"
	this.lengthFieldLength.PlaceHolder = "长度字段所占的字节数"
	this.lengthAdjustment.PlaceHolder = "长度的调整值"
	this.initialBytesToStrip.PlaceHolder = "解码后跳过的字节数"
	this.maxFrameLength.PlaceHolder = "数据包最大长度"
	this.orderSel = widget.NewSelect([]string{"大端模式", "小端模式"}, func(item string) {
		AppData.Lengthfield.Order = item
	})

	entrys1 := container.NewVBox(this.lengthFieldOffset, this.lengthFieldLength, this.lengthAdjustment, this.initialBytesToStrip, this.maxFrameLength, this.orderSel)
	this.lenthfield = container.NewHBox(labels1, entrys1)

	this.workhbox = container.New(layout.NewFormLayout(), this.worklabel, this.workmodeSel)
	hbox1 := container.New(layout.NewFormLayout(), this.destiplabel, this.destipEntry)
	hbox2 := container.New(layout.NewFormLayout(), this.destportlabel, this.destportEntry)
	this.clienthbox = container.NewVBox(hbox1, hbox2)

	hbox3 := container.New(layout.NewFormLayout(), this.localportlabel, this.localportEntry)
	url, _ := url.Parse("https://netty.io/4.0/api/io/netty/handler/codec/LengthFieldBasedFrameDecoder.html")
	hbox4 := container.NewHBox(this.setStickyCheck, widget.NewHyperlink("断粘包规则", url))
	this.serverhbox = container.NewVBox(hbox3, hbox4)

	this.workmodeview = container.NewVBox(this.workhbox)
	return container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), this.workmodeview)
}

func (this *Left) GetView() fyne.CanvasObject {
	return this.rootview
}

func (this *Left) disable() {
	this.workmodeSel.Disable()
	this.localportEntry.Disable()
	this.destipEntry.Disable()
	this.destportEntry.Disable()
	this.setStickyCheck.Disable()

	this.lengthFieldOffset.Disable()
	this.lengthFieldLength.Disable()
	this.lengthAdjustment.Disable()
	this.initialBytesToStrip.Disable()
	this.maxFrameLength.Disable()
	this.orderSel.Disable()

	//this.recvHexCheck.Disable()

	//this.sendHexCheck.Disable()
	this.sendStickyCheck.Disable()
	this.sendIntervalCheck.Disable()
	this.sendEverySecEntry.Disable()
}
func (this *Left) enable() {
	this.workmodeSel.Enable()
	this.localportEntry.Enable()
	this.destipEntry.Enable()
	this.destportEntry.Enable()
	this.setStickyCheck.Enable()

	this.lengthFieldOffset.Enable()
	this.lengthFieldLength.Enable()
	this.lengthAdjustment.Enable()
	this.initialBytesToStrip.Enable()
	this.maxFrameLength.Enable()
	this.orderSel.Enable()

	//this.recvHexCheck.Enable()

	//this.sendHexCheck.Enable()
	this.sendStickyCheck.Enable()
	this.sendIntervalCheck.Enable()
	this.sendEverySecEntry.Enable()
}

func (this *Left) initData() {
	if !AppData.IsCache {
		return
	}
	this.workmodeSel.SetSelected(AppData.Work.Workmode)
	this.localportEntry.SetText(strconv.Itoa(AppData.Work.Localport))
	this.destipEntry.SetText(AppData.Work.Destip)
	this.destportEntry.SetText(strconv.Itoa(AppData.Work.Destport))

	this.setStickyCheck.SetChecked(AppData.Work.Stickymode)
	this.lengthFieldOffset.SetText(strconv.Itoa(AppData.Lengthfield.LengthFieldOffset))
	this.lengthFieldLength.SetText(strconv.Itoa(AppData.Lengthfield.LengthFieldLength))
	this.lengthAdjustment.SetText(strconv.Itoa(AppData.Lengthfield.LengthAdjustment))
	this.initialBytesToStrip.SetText(strconv.Itoa(AppData.Lengthfield.InitialBytesToStrip))
	this.maxFrameLength.SetText(strconv.FormatUint(AppData.Lengthfield.MaxFrameLength, 10))
	this.orderSel.SetSelected(AppData.Lengthfield.Order)

	this.recvHexCheck.SetChecked(AppData.Recv.Hex)

	this.sendHexCheck.SetChecked(AppData.Send.Hex)
	this.sendStickyCheck.SetChecked(AppData.Send.Stickymode)
	this.sendEverySecEntry.SetText(strconv.FormatUint(AppData.Send.Interval, 10))
	this.sendIntervalCheck.SetChecked(AppData.Send.Sendmode)

}

func (this *Left) saveData() {
	AppData.Work.Workmode = this.workmodeSel.Selected
	value, err := strconv.Atoi(this.localportEntry.Text)
	if err == nil {
		AppData.Work.Localport = value
	}
	AppData.Work.Destip = this.destipEntry.Text

	value, err = strconv.Atoi(this.destportEntry.Text)
	if err == nil {
		AppData.Work.Destport = value
	}

	AppData.Work.Stickymode = this.setStickyCheck.Checked
	value, err = strconv.Atoi(this.lengthFieldOffset.Text)
	if err == nil {
		AppData.Lengthfield.LengthFieldOffset = value
	}
	value, err = strconv.Atoi(this.lengthFieldLength.Text)
	if err == nil {
		AppData.Lengthfield.LengthFieldLength = value
	}
	value, err = strconv.Atoi(this.lengthAdjustment.Text)
	if err == nil {
		AppData.Lengthfield.LengthAdjustment = value
	}
	value, err = strconv.Atoi(this.initialBytesToStrip.Text)
	if err == nil {
		AppData.Lengthfield.InitialBytesToStrip = value
	}
	u_value_64, err1 := strconv.ParseUint(this.maxFrameLength.Text, 10, 64)
	if err1 == nil {
		AppData.Lengthfield.MaxFrameLength = u_value_64
	}
	AppData.Lengthfield.Order = this.orderSel.Selected
	AppData.Recv.Hex = this.recvHexCheck.Checked
	AppData.Send.Hex = this.sendHexCheck.Checked
	AppData.Send.Stickymode = this.sendStickyCheck.Checked
	AppData.Send.Sendmode = this.sendIntervalCheck.Checked
	u_value_64, err1 = strconv.ParseUint(this.sendEverySecEntry.Text, 10, 64)
	if err1 == nil {
		AppData.Send.Interval = u_value_64
	}

}

func (this *Left) OnConnStatus(b bool) {
	this.rightview.OnStartClick(b)
	if b {
		//this.startBtn.Importance = widget.DangerImportance
		this.startBtn.SetText("关闭")
		this.disable()
	} else {
		//this.startBtn.Importance = widget.HighImportance
		this.startBtn.SetText("打开")
		this.enable()
	}
}
func (this *Left) onClick(viewid int, args ...interface{}) func() {
	switch viewid {
	case 0:
		return func() {
			for _, arg := range args {
				fmt.Println(arg)
			}
			if this.startBtn.Text == "打开" {
				this.saveData()
				this.mainview.Start()
			} else {
				this.mainview.Stop()
				this.OnConnStatus(false)
			}
		}
	case 1:
		return func() {
			this.rightview.OnRecvClearClick()
		}
	case 2:
		return func() {
			this.rightview.OnSendClearClick()
		}
	case 3:
		return func() {
			_dialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, this.mainview.GetWindow())
					return
				}
				if writer == nil {
					return
				}
				defer writer.Close()
				if this.rightview != nil {
					buff := this.rightview.GetRecvData()
					if buff != nil {
						_, err = writer.Write(buff.Bytes())
						if err != nil {
							dialog.ShowError(err, this.mainview.GetWindow())
						}
					}
				}

			}, this.mainview.GetWindow())
			_dialog.SetFileName(time.Now().Format("tcp-recv-2006-01-02-15-04-05.txt"))
			_dialog.Show()
		}
	}
	return nil
}

func (this *Left) recvView() *fyne.Container {
	this.recvHexCheck = widget.NewCheck("十六进制接收", func(b bool) {
		this.rightview.OnHexCheck(b)
	})
	this.saverecvbtn = widget.NewButton("保存接收窗口数据", this.onClick(3))
	this.clearLogBtn = widget.NewButton("清空接收窗口数据", this.onClick(1))

	view := container.NewVBox(this.recvHexCheck, container.NewHBox(this.saverecvbtn, this.clearLogBtn))
	return container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), view)
}

func (this *Left) sendView() *fyne.Container {
	this.sendHexCheck = widget.NewCheck("十六进制接收", func(b bool) {
		this.rightview.OnSendDataHex(b)
	})
	this.sendStickyCheck = widget.NewCheck("分包设置(可模拟断粘包)", func(b bool) {
		if this.rightview != nil {
			if this.rightview.OnStickyClick != nil {
				this.rightview.OnStickyClick(b)
			}
		}
	})
	this.sendIntervalCheck = widget.NewCheck("每隔", func(b bool) {
		text := this.sendEverySecEntry.Text
		if text == "" {
			this.sendIntervalCheck.Checked = false
			dialog.ShowError(errors.New("请输入间隔时间"), this.mainview.GetWindow())
			return
		}
		inter, err := strconv.Atoi(text)
		if err != nil {
			this.sendIntervalCheck.Checked = false
			dialog.ShowError(err, this.mainview.GetWindow())
			return
		}
		if inter <= 0 {
			this.sendIntervalCheck.Checked = false
			dialog.ShowError(errors.New("间隔时间小于0"), this.mainview.GetWindow())
			return
		}
		this.rightview.SetSendInterval(b, inter)
	})
	this.clearLogBtn = widget.NewButton("清空发送窗口数据", this.onClick(2))
	hexbtn := widget.NewButton("hex", func() {
		NewWindow().Show()
	})
	this.sendEverySecEntry = widget.NewEntry()
	sendsec := container.NewHBox(this.sendIntervalCheck, container.NewGridWrap(fyne.NewSize(70, 25), this.sendEverySecEntry), widget.NewLabel("毫秒发送"))
	view := container.NewVBox(this.sendHexCheck, this.sendStickyCheck, sendsec, container.NewHBox(hexbtn, this.clearLogBtn))
	cc := container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), view)
	return cc
}

func getIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() && ipNet.IP.IsPrivate() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func getPrivateIp() string {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ""
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() && ipNet.IP.IsPrivate() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

func (this *Left) displayView() *fyne.Container {
	view := container.NewHBox(widget.NewLabel(fmt.Sprintf("本机IP：%s", getPrivateIp())))
	return container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), view)
}
