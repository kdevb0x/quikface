// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// errchan is an error queue, implemented using a chan.
type errchan chan error

type ErrorType string

const (
	ErrHTTPRequestParseError ErrorType = "parsing error: unexpected or malformed HTTP Request."
)

// Error method implements the error interface.
func (et ErrorType) Error() string {
	return string(et)
}

// errorLog is the global error queue.
// Because the http.Handler interface has no return parameters, errors
// encountered within handlers can send them here.
var errorLog = make(errchan, 1000)

// throwError writes err to log.Println, as well as all of the optional
// io.Writer parameters passed in.
func throwError(err error, out ...io.Writer) {
	log.Println(fmt.Errorf("ThrowError: %w\n", err))
	if len(out) > 0 {
		for _, w := range out {

			switch rw := w.(type) {
			case http.ResponseWriter:
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			default:
				w.Write([]byte(fmt.Errorf("ThrowError: %w\n", err).Error()))

			}
		}
	}
}
