package view

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
	"runtime"
	vlog "tcptest/util"
)

type Logview interface {
	Print(level vlog.Level, log string)
	GetView() fyne.CanvasObject
	SetLogLineSize(maxSize int)
	GetText() *bytes.Buffer
	Clear()
}

var (
	red    = color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}
	blue   = color.RGBA{R: 65, G: 105, B: 225, A: 0xff}
	orange = color.RGBA{R: 0xff, G: 97, B: 0, A: 0xff}
	green  = color.RGBA{R: 0, G: 0xff, B: 0, A: 0xff}
)

type Data struct {
	_autoScroll  bool
	_scroll      *container.Scroll
	_container   *fyne.Container
	_window      fyne.Window
	_logLineSize int
}

func (this *Data) Clear() {
	this._container.RemoveAll()
}

func (this *Data) SetLogLineSize(maxSize int) {
	if maxSize <= 100 {
		return
	}
	this._logLineSize = maxSize
}

func (this *Data) GetView() fyne.CanvasObject {
	return this._scroll
}

func NewLogView(window fyne.Window) Logview {
	c := &Data{
		_window:      window,
		_autoScroll:  true,
		_logLineSize: 2000,
	}
	c.initView()
	return c
}

func (this *Data) getSize() int {
	size := len(this._container.Objects)
	return size
}

func (this *Data) GetText() *bytes.Buffer {
	buff := new(bytes.Buffer)
	for _, item := range this._container.Objects {
		var str = item.(*canvas.Text).Text
		buff.WriteString(str)
		buff.WriteString("\n")
	}
	return buff
}

func (this *Data) handleTypedKey(ke *fyne.KeyEvent) {
	if this.getSize() <= 0 {
		return
	}
	delta := this._container.Objects[0].Size().Height
	switch ke.Name {
	case fyne.KeyUp:
		this._scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: delta}})
	case fyne.KeyDown:
		this._scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: -delta}})
	default:
		return
	}
	this._autoScroll = false
}

func (this *Data) handleTypedRune(r rune) {
	switch r {
	case 't':
		this._scroll.ScrollToTop()
		this._autoScroll = false
	case 'b':
		this._scroll.ScrollToBottom()
		this._autoScroll = true
	case 'c':
		this._container.RemoveAll()
	}
}

func (this *Data) Print(level vlog.Level, line string) {
	switch level {
	case vlog.INFO:
		this._container.Add(canvas.NewText(line, blue))
	case vlog.ERROR:
		this._container.Add(canvas.NewText(line, red))
	case vlog.DEBUG:
		this._container.Add(canvas.NewText(line, color.White))
	case vlog.WARM:
		this._container.Add(canvas.NewText(line, orange))
	default:
		this._container.Add(canvas.NewText(line, green))
	}
	this._container.Refresh()
	if this.getSize() > this._logLineSize {
		this._container.RemoveAll()
	}

	title := fmt.Sprintf("TCP/IP工具 [%s/%s] 日志行数:%d", runtime.GOOS, runtime.GOARCH, this.getSize())
	//title := fmt.Sprintf("碧丽水机测试工具[基于GO语言] %s %s 日志行数:%d ", runtime.GOOS, runtime.GOARCH, this.getSize())
	this._window.SetTitle(title)
	this._scroll.ScrollToBottom()
}

func (this *Data) initView() {
	this._container = fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	this._scroll = container.NewScroll(this._container)
	this._window.Canvas().SetOnTypedKey(this.handleTypedKey)
	this._window.Canvas().SetOnTypedRune(this.handleTypedRune)
}
