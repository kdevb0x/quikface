// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	_ "net"
	_ "net/http"

	"honnef.co/go/js/dom/v2"
)

func getUserMedia(c *Client) error {
	doc := dom.GetWindow().Document()

}
