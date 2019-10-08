// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/andlabs/ui"
)

type video_viewport struct {
	imgwidget *ui.Image
}

func ExecOnNewThread(f unsafe.Pointer) {
	runtime.LockOSThread()

}

func buildGUI() {
	err := ui.Main(buildGUI)
	if err != nil {
		panic(err)
	}

	// var oswidth, osheight int
	for _, s := range os.Environ() {
		println(s)
	}
	mainWin := ui.NewWindow("quikface vc", 640, 480, false)
	mainWin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWin.Destroy()
		return true
	})
	vb := ui.NewVerticalBox()
	mainWin.SetChild(vb)
	mainWin.SetMargined(false)

	hb := ui.NewHorizontalBox()
	vb.Append(hb)

	create := ui.NewButton("Create Session")
	hb.Append(create, false)
}

func main() {
	buildGUI()
}
