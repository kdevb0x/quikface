// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

// errchan is an error queue, implemented using a chan.
type errchan chan error

// errorLog is the global error queue.
// Because the http.Handler interface has no return parameters, errors
// encountered within handlers can send them here.
var errorLog = make(errchan, 1000)
