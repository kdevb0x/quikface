// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"errors"
	"syscall/js"
)

func getUserMedia(dstTag string) error {
	var video = js.Global().Get("document").Call("querySelector", "#"+dstTag)
	if gum := js.Global().Get("navigator").Get("mediadevices").Call("getUserMedia"); gum {
		gum.Call("onSuccess", js.FuncOf(jsStreamCallback))
		gum.Call("onError", js.FuncOf(jsErrorCallback))
	} else if !gum {
		return errors.New("error: call to getUserMedia failed")
	}

	// TODO: Check the that the init code above for the DOM api is correct,
	// then finish setup.
	return nil
}

func jsStreamCallback(this js.Value, args []js.Value) interface{} {

}

func jsErrorCallback(this js.Value, args []js.Value) interface{} {

}
