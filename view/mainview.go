package view

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/aceld/zinx/ziface"
	"runtime"
	"tcptest/entity"
	"tcptest/tcp"
	"tcptest/the"
	vlog "tcptest/util"
)

var _width float32 = 800
var _height float32 = 600
var AppData *entity.Entity

type MainView struct {
	_app      fyne.App
	_win      fyne.Window
	rightview RightView
	leftview  LeftView
	tcptest   tcp.TcpTest
}

func NewView() *MainView {
	this := &MainView{}
	return this
}

func (this *MainView) GetWindow() fyne.Window {
	return this._win
}

func (this *MainView) makeView(leftView, rightView fyne.CanvasObject) *container.Split {
	_view := container.NewHSplit(leftView, rightView)
	_view.SetOffset(0.25)
	return _view
}

func (this *MainView) Start() {
	arr, err := json.Marshal(AppData)
	if err == nil {
		this._app.Preferences().SetString("data", string(arr))
	}

	if this.tcptest == nil {
		this.tcptest = tcp.NewTcpTest()
	}

	this.tcptest.Init(AppData)
	this.tcptest.GetTcp().SetRecv(this.rightview.OnMessageReceiver)
	this.rightview.SetClientsFunc(this.tcptest.GetTcp().GetMaps)
	this.rightview.SetclientsSetFunc(this.tcptest.GetTcp().SetClientArray)
	this.tcptest.GetTcp().SetOnConnStart(this.onConnStart)
	this.tcptest.GetTcp().SetOnConnStop(this.onConnStop)
	this.tcptest.GetTcp().Start()
	if this.tcptest.GetTcp().IsServer() {
		this.leftview.OnConnStatus(true)
	}
}

func (this *MainView) onConnStart(connection ziface.IConnection) {
	this.rightview.OnConnDetailLog(vlog.INFO, fmt.Sprintf("成功连接 %s", connection.GetConnection().RemoteAddr()))
	if !this.tcptest.GetTcp().IsServer() {
		this.leftview.OnConnStatus(true)
	}

}

func (this *MainView) onConnStop(connection ziface.IConnection) {
	this.rightview.OnConnDetailLog(vlog.ERROR, fmt.Sprintf("连接断开 %s", connection.GetConnection().RemoteAddr()))
	if !this.tcptest.GetTcp().IsServer() {
		this.leftview.OnConnStatus(false)
	}
}

func (this *MainView) Stop() {
	if this.tcptest != nil {
		this.tcptest.GetTcp().Stop()
	}
}

func (this *MainView) Send(data [][]byte) {
	for i, arr := range data {
		fmt.Println(i, arr)
		this.tcptest.GetTcp().Send(arr)
	}
}

func (this *MainView) RunFyneView() {
	this._app = app.NewWithID("io.github.tcp.ip")
	this._app.Settings().SetTheme(&the.MyTheme{})
	title := fmt.Sprintf("TCP/IP工具 [%s/%s]", runtime.GOOS, runtime.GOARCH)
	this._win = this._app.NewWindow(title)
	this._win.Resize(fyne.NewSize(_width, _height))
	_data := this._app.Preferences().String("data")
	AppData = new(entity.Entity)
	if _data != "" {
		err := json.Unmarshal([]byte(_data), &AppData)
		if err == nil {
			fmt.Println("有缓存...")
			AppData.IsCache = true
		} else {
			fmt.Println("没有缓存...")
			AppData.IsCache = false
		}
	} else {
		fmt.Println("没有缓存...")
		AppData.IsCache = false
	}

	this.rightview = NewRightView(this)
	this.leftview = NewLeftView(this, this.rightview)

	splitView := this.makeView(this.leftview.GetView(), this.rightview.GetView())
	this._win.SetContent(splitView)
	this._win.ShowAndRun()
}
