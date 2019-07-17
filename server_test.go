// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"net/http"
	"testing"
)

func TestListenAndServeTLS(t *testing.T) {
	getFunc := func(t *testing.T) {
		c := http.DefaultClient
		resp, err := c.Get("localhost:8080")
		if err != nil {
			t.Fail()
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fail()
		}

	}
	s := NewServer()
	s.ListenAndServeTLS()
	go getFunc(t)
}
