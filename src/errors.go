// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
	// initialize ErrorLog and spawn goroutine to loop through the errors
	// sent to it, printing them to os.Stderr
	ErrorLog = make(Errchan, 1000)
	go func(c Errchan) {
		for range c {
			select {
			case e := <-c:
				// sometimes Close() or other func that return an error
				// will be sent (usually deferred), so we should
				// check for and discard the nil values.
				if e == nil {
					continue
				}

				go throwError(e, os.Stderr)
			}
		}
	}(ErrorLog)
}

// Errchan is an error queue, implemented using a chan.
type Errchan chan error

type ErrorType string

const (
	ErrHTTPRequestParseError ErrorType = "parsing error: unexpected or malformed HTTP Request."
)

// Error method implements the error interface.
func (et ErrorType) Error() string {
	return string(et)
}

// ErrorLog is the global error queue.
// Because the http.Handler interface has no return parameters, errors
// encountered within handlers can send them here.
//
// Likewise, any procedure that is unable to return an error (such as
// concurrently executed code ran with the `go` keyword) will send any
// errors encountered.
var ErrorLog Errchan

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
