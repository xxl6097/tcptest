package view

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

type HexWindow interface {
	Show()
}

type HexData struct {
	_win          fyne.Window
	enter, result *widget.Entry
	menu          *widget.Select
	selectType    int
}

var opts = []string{"十六进制转十进制", "十进制转十六进制", "字节数组转int", "无符号[]byte转hex", "有符号[]byte转hex", "hex追加0x前缀"}

func NewWindow() HexWindow {
	win := &HexData{
		_win: fyne.CurrentApp().NewWindow("hex转换工具"),
	}
	return win
}

func (this *HexData) hexToInteger(s string) {
	size := len(s)
	bytesize := size / 2
	if size%2 == 0 {
		arr, err := hex.DecodeString(s)
		if err == nil {
			//var x uint
			//for i, u := range arr {
			//	x |= uint(u) << ((len(arr) - 1 - i) * 8)
			//	fmt.Println(x)
			//}
			var x uint64
			for i, u := range arr {
				offset := len(arr) - 1 - i
				if i == 0 {
					x = uint64(u) << (offset * 8)
				} else {
					y := uint64(u) << (offset * 8)
					x |= y
				}
			}

			var y interface{}
			var z interface{}
			switch bytesize {
			case 1:
				y = int8(x)
				z = uint8(x)
			case 2:
				y = int16(x)
				z = uint16(x)
			case 3:
				y = int32(x)
				z = uint32(x)
			case 4:
				y = int32(x)
				z = uint32(x)
			default:
				y = int64(x)
				z = uint64(x)
			}

			str := fmt.Sprintf("大小:%d字节\n无符号:%d\n有符号:%d", bytesize, z, y)
			this.result.SetText(str)
		}
	}
}

func (this *HexData) arrToInt(ss string) {
	arr := strings.Split(ss, " ")
	var x uint64
	bytesize := len(arr)
	for i, s := range arr {
		u, _ := strconv.Atoi(strings.TrimSpace(s))
		offset := len(arr) - 1 - i
		if i == 0 {
			x = uint64(u) << (offset * 8)
		} else {
			y := uint64(u) << (offset * 8)
			x |= y
		}

	}

	var y interface{}
	var z interface{}
	switch bytesize {
	case 1:
		y = int8(x)
		z = uint8(x)
	case 2:
		y = int16(x)
		z = uint16(x)
	case 3:
		y = int32(x)
		z = uint32(x)
	case 4:
		y = int32(x)
		z = uint32(x)
	default:
		y = int64(x)
		z = uint64(x)
	}
	str := fmt.Sprintf("大小:%d字节\n无符号:%d\n有符号:%d", bytesize, z, y)
	this.result.SetText(str)
}

func (this *HexData) integerToHex(s string) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		this.result.SetText(fmt.Sprintf("%X", uint16(v)))
	}
}
func (this *HexData) unbytesArrToHexArr(s string) {
	string_slice := strings.Split(s, " ")
	var bf bytes.Buffer
	for _, s := range string_slice {
		v, _ := strconv.Atoi(strings.TrimSpace(s))
		bf.WriteString(fmt.Sprintf("%X", uint8(v)))
	}
	this.result.SetText(bf.String())
}
func (this *HexData) bytesArrToHexArr(ss string) {
	var bf bytes.Buffer
	arr := strings.Split(ss, ",")
	for _, s := range arr {
		v, _ := strconv.Atoi(strings.TrimSpace(s))
		bf.WriteString(fmt.Sprintf("%X", uint8(v)))
	}
	this.result.SetText(bf.String())
}
func (this *HexData) hexArrAdd0x(s string) {
	arr, err := hex.DecodeString(s)
	if err == nil {
		var bf bytes.Buffer
		for i, b := range arr {
			if i == len(arr)-1 {
				bf.WriteString(fmt.Sprintf("0x%X", b))
			} else {
				bf.WriteString(fmt.Sprintf("0x%X,", b))
			}
		}
		this.result.SetText(bf.String())
	}
}
func (this *HexData) onTextChange(s string) {
	switch this.menu.SelectedIndex() {
	case 0:
		this.hexToInteger(this.enter.Text)
	case 1:
		this.integerToHex(this.enter.Text)
	case 2:
		this.arrToInt(this.enter.Text)
	case 3:
		this.unbytesArrToHexArr(this.enter.Text)
	case 4:
		this.bytesArrToHexArr(this.enter.Text)
	case 5:
		this.hexArrAdd0x(this.enter.Text)
	}
}

func (this *HexData) onSelect(s string) {
	this.selectType = this.menu.SelectedIndex()
	this.enter.SetText("")
	this.result.SetText("")
	switch this.menu.SelectedIndex() {
	case 0:
		this.enter.SetPlaceHolder("请输入十六进制数据")
	case 1:
		this.enter.SetPlaceHolder("请输入十进制数据")
	case 2:
		this.enter.SetPlaceHolder("字节数组转int\n转换示例：11 22 ==> 2838")
	case 3:
		this.enter.SetPlaceHolder("请输入无符号[]byte数据\n转换示例：161 228 22 35 ==> A1E41623")
	case 4:
		this.enter.SetPlaceHolder("请输入有符号[]byte数据\n转换示例：-94, -30, 22, 16 ==> A2E21610")
	case 5:
		this.enter.SetPlaceHolder("请输入十六进制数组\n转换示例：A1E41623 ==> 0xA1,0xE4,0x16,0x23")
	}
}

func (this *HexData) makeView() *fyne.Container {
	this.enter = widget.NewMultiLineEntry()
	this.enter.SetMinRowsVisible(3)
	this.enter.SetPlaceHolder("请严格输入十六进制数据")
	this.enter.OnChanged = this.onTextChange

	this.result = widget.NewMultiLineEntry()
	this.result.SetMinRowsVisible(10)
	this.result.SetPlaceHolder("显示转换结果")

	this.menu = widget.NewSelect(opts, this.onSelect)
	this.menu.SetSelectedIndex(0)

	return container.NewVBox(this.enter, this.menu, this.result)
}

func (this *HexData) Show() {
	this._win.SetContent(this.makeView())
	this._win.Resize(fyne.NewSize(350, 180))
	//this._win.SetFixedSize(true)
	this._win.Show()
}
