#!/bin/bash
#https://www.topgoer.cn/docs/fyne/fyne-1e2fi58566ilg
function make_windows() {
  fyne-cross windows -arch=amd64,386
#  fyne-cross windows -arch=*
}
function make_macos() {
  #sudo xattr -r -d io.github.bili
  fyne-cross darwin -arch=amd64,arm64 -app-id=io.github.bili
  #fyne-cross darwin -arch=amd64 -app-id=io.github.bili
  # fyne-cross ios -app-id=io.github.bili
#  fyne-cross darwin -arch=*
# fyne package -os darwin -appID io.github.bili -icon Icon.png

# CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build *.go
}
function make_android() {
#  fyne package -os android -appID com.example.bili
#  fyne-cross android -arch=*
#fyne-cross ios -arch=arm,arm64 -app-id=io.github.bili
  fyne-cross android -arch=arm,arm64 -app-id=io.github.bili
}

function make_all() {
    fyne-cross windows -arch=amd64,386
    fyne-cross darwin -arch=amd64,arm64 -app-id=io.github.bili
    fyne-cross android -arch=arm,arm64 -app-id=io.github.bili
}
function menu() {
  echo "1. 编译 Windows"
  echo "2. 编译 MacOS"
  echo "3. 编译 Android"
  echo "4. 编译全平台"
  echo "请输入编号:"
  read index

  case "$index" in
  [1]) (make_windows) ;;
  [2]) (make_macos) ;;
  [3]) (make_android) ;;
  [4]) (make_all) ;;
  *) echo "exit" ;;
  esac
}

menu


#CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CC=/opt/homebrew/Cellar/mingw-w64/10.0.0_4/toolchain-x86_64/bin/x86_64-w64-mingw32-gcc go build *.go
