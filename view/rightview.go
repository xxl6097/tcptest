package view

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/simplifiedchinese"
	vlog "tcptest/util"
	"time"
)

type Right struct {
	mainview                                                  *MainView
	localport, destip, destport                               *widget.Entry
	firstEntry, secondEntry, thirdEntry                       *widget.Entry
	workmode, chsClientBtn                                    *widget.Select
	startBtn, sendBtn, stopBtn, clearBtn                      *widget.Button
	hexCheck, sendCheck                                       *widget.Check
	recvlogview, sendlogview                                  Logview
	rootview, upview, downview, sendbtnview, sendandclearview *fyne.Container
	sendview, stickyview                                      *fyne.Container
	isHexDisplay, isHexSend, bIsSticky                        bool
	clientsFunc                                               func() map[string]uint64
	clientsSetFunc                                            func([]uint64)
	bIsSendInterval                                           bool
	sendInterval                                              int
	ctx                                                       context.Context
	cancel                                                    context.CancelFunc
}

type RightView interface {
	GetView() fyne.CanvasObject
	GetRecvData() *bytes.Buffer
	OnWorkModeSelect(string)
	OnStickyClick(b bool)
	OnStartClick(status bool)
	OnSendClearClick()
	OnRecvClearClick()
	OnHexCheck(status bool)
	OnSendDataHex(status bool)
	OnMessageReceiver([]byte)
	SetClientsFunc(fc func() map[string]uint64)
	SetclientsSetFunc(fc func([]uint64))
	OnConnDetailLog(vlog.Level, string)
	SetSendInterval(bool, int)
}

func NewRightView(_mainview *MainView) RightView {
	view := &Right{
		mainview: _mainview,
	}
	view.recvlogview = NewLogView(_mainview.GetWindow())
	view.sendlogview = NewLogView(_mainview.GetWindow())
	view.initView()
	return view
}

func (this *Right) OnConnDetailLog(level vlog.Level, s string) {
	this.sendlogview.Print(level, s)
}
func (this *Right) SetClientsFunc(fc func() map[string]uint64) {
	this.clientsFunc = fc
}
func (this *Right) SetclientsSetFunc(fc func([]uint64)) {
	this.clientsSetFunc = fc
}

func (this *Right) GetView() fyne.CanvasObject {
	return this.rootview
}
func (this *Right) GetRecvData() *bytes.Buffer {
	return this.recvlogview.GetText()
}
func (this *Right) initView() {
	setview := this.recvlogview.GetView()
	this.recvlogview.Print(vlog.ERROR, "hello wolrd.....")

	this.firstEntry = widget.NewMultiLineEntry()
	this.secondEntry = widget.NewMultiLineEntry()
	this.thirdEntry = widget.NewMultiLineEntry()
	this.firstEntry.PlaceHolder = "请输入要发送的数据内容"

	this.upview = container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), setview)

	this.sendBtn = widget.NewButton("发送", this.onSendClick)
	this.sendBtn.Disable()
	this.stopBtn = widget.NewButton("停止", func() {
		this.sendBtn.Enable()
		this.cancel()
	})
	this.clearBtn = widget.NewButton("清空信息", func() {
		this.sendlogview.Clear()
	})
	this.chsClientBtn = widget.NewSelect([]string{"全部客户端", "选择..."}, this.onGetClientsClick)
	this.chsClientBtn.SetSelected("全部客户端")
	this.chsClientBtn.Disable()

	this.sendandclearview = container.NewVBox(this.sendBtn, this.clearBtn)

	this.sendbtnview = container.NewBorder(nil, nil, nil, this.sendandclearview, this.firstEntry)

	this.stickyview = container.NewGridWithRows(3, this.sendbtnview, this.secondEntry, this.thirdEntry)

	sendlogdetailview := container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), this.sendlogview.GetView())
	this.sendview = container.NewGridWithRows(2, this.sendbtnview, sendlogdetailview)

	this.rootview = container.NewGridWithRows(2, this.upview, this.sendview)
}
func (this *Right) OnMessageReceiver(bytes []byte) {
	if this.isHexDisplay {
		this.recvlogview.Print(vlog.DEBUG, hex.EncodeToString(bytes))
	} else {
		buf, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bytes)
		this.recvlogview.Print(vlog.DEBUG, string(buf))
	}
}

func (this *Right) getSendData() [][]byte {
	if this.bIsSticky {
		data1 := this.firstEntry.Text
		data2 := this.secondEntry.Text
		data3 := this.thirdEntry.Text
		fmt.Println(data1, data2, data3)
		data := make([][]byte, 3)
		if this.isHexSend {
			h1, err := hex.DecodeString(data1)
			if err != nil {
				dialog.ShowError(errors.New(fmt.Sprintf("请输入十六进制数据\n错误信息如下:\n%s", err.Error())), this.mainview.GetWindow())
				return nil
			}
			data[0] = h1
			h2, err1 := hex.DecodeString(data2)
			if err1 != nil {
				dialog.ShowError(errors.New(fmt.Sprintf("请输入十六进制数据\n错误信息如下:\n%s", err1.Error())), this.mainview.GetWindow())
				return nil
			}
			data[1] = h2
			h3, err2 := hex.DecodeString(data3)
			if err2 != nil {
				dialog.ShowError(errors.New(fmt.Sprintf("请输入十六进制数据\n错误信息如下:\n%s", err2.Error())), this.mainview.GetWindow())
				return nil
			}
			data[2] = h3
		} else {
			data1, _ = simplifiedchinese.GBK.NewEncoder().String(data1)
			data2, _ = simplifiedchinese.GBK.NewEncoder().String(data2)
			data3, _ = simplifiedchinese.GBK.NewEncoder().String(data3)
			data[0] = []byte(data1)
			data[1] = []byte(data2)
			data[2] = []byte(data3)
		}
		//this.mainview.Send(data)
		return data
	} else {
		data1 := this.firstEntry.Text
		fmt.Println(data1)
		data := make([][]byte, 1)
		if this.isHexSend {
			h1, err := hex.DecodeString(data1)
			if err != nil {
				dialog.ShowError(errors.New(fmt.Sprintf("请输入十六进制数据\n错误信息如下:\n%s", err.Error())), this.mainview.GetWindow())
				return nil
			}
			data[0] = h1
		} else {
			data1, _ = simplifiedchinese.GBK.NewEncoder().String(data1)
			data[0] = []byte(data1)
		}
		return data
	}
	return nil
}

func (this *Right) SetSendInterval(b bool, i int) {
	this.bIsSendInterval = b
	this.sendInterval = i
	this.sendandclearview.Remove(this.clearBtn)
	this.sendandclearview.Remove(this.stopBtn)
	if b {
		this.sendandclearview.Add(this.stopBtn)
		this.sendandclearview.Add(this.clearBtn)
	} else {
		this.sendandclearview.Add(this.clearBtn)
	}
}

func (this *Right) OnWorkModeSelect(item string) {
	this.sendandclearview.RemoveAll()
	if item == "TCP服务器" {
		this.sendandclearview.Add(this.chsClientBtn)
		this.sendandclearview.Add(this.sendBtn)
		this.sendandclearview.Add(this.clearBtn)
	} else {
		this.sendandclearview.Add(this.sendBtn)
		this.sendandclearview.Add(this.clearBtn)
	}
	this.SetSendInterval(this.bIsSendInterval, this.sendInterval)
}

func (this *Right) onSendClick() {
	if this.bIsSendInterval {
		this.ctx, this.cancel = context.WithCancel(context.Background())
		this.sendBtn.Disable()
		go func() {
			defer this.cancel()
			for {
				select {
				case <-this.ctx.Done():
					return
				case <-time.After(time.Millisecond * time.Duration(this.sendInterval)):
					this.mainview.Send(this.getSendData())
				}
			}
		}()

	} else {
		this.mainview.Send(this.getSendData())
	}
}

func (this *Right) OnStickyClick(b bool) {
	this.sendandclearview.RemoveAll()
	this.bIsSticky = b
	if b {
		this.rootview.Remove(this.sendview)
		this.rootview.Add(this.stickyview)
		this.firstEntry.PlaceHolder = "请输入第一包数据（一般为一包半）"
		this.secondEntry.PlaceHolder = "请输入第二包数据（一般为上面剩余半包）"
		this.thirdEntry.PlaceHolder = "请输入第三包数据（完整的两包数据）"
		this.sendandclearview.Add(this.sendBtn)
	} else {
		this.rootview.Remove(this.stickyview)
		this.rootview.Add(this.sendview)
		this.firstEntry.PlaceHolder = "请输入要发送的数据内容"
		this.sendandclearview.Add(this.sendBtn)
		this.sendandclearview.Add(this.clearBtn)
	}
	this.rootview.Refresh()
}

func (this *Right) OnStartClick(status bool) {
	if status {
		this.sendBtn.Enable()
		this.chsClientBtn.Enable()
	} else {
		this.sendBtn.Disable()
		this.chsClientBtn.Disable()
		if this.cancel != nil {
			this.cancel()
		}
	}
}

func (this *Right) OnSendClearClick() {
	this.firstEntry.SetText("")
	this.secondEntry.SetText("")
	this.thirdEntry.SetText("")
}

func (this *Right) OnRecvClearClick() {
	this.recvlogview.Clear()
}
func (this *Right) OnHexCheck(status bool) {
	this.isHexDisplay = status
}

func (this *Right) OnSendDataHex(status bool) {
	this.isHexSend = status
}

func (this *Right) onGetClientsClick(sel string) {
	if this.clientsSetFunc == nil {
		return
	}
	if sel == "全部客户端" {
		this.clientsSetFunc(nil)
	} else {
		if this.clientsFunc != nil {
			maps := this.clientsFunc()
			var strArr []string
			for key, value := range maps {
				fmt.Println(key, value)
				strArr = append(strArr, key)
			}
			cp := widget.NewCheckGroup(strArr, nil)
			clientDlg := dialog.NewCustom("客户端列表", "确定", container.NewScroll(cp), this.mainview.GetWindow())
			clientDlg.SetOnClosed(func() {
				value := make([]uint64, 0)
				for _, item := range cp.Selected {
					a := maps[item]
					value = append(value, a)
				}
				this.clientsSetFunc(value)
			})
			clientDlg.Resize(fyne.NewSize(200, 300))
			clientDlg.Show()
		}
	}

}
