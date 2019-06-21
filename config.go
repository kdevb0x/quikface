// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

func createConfigDB(filename string) (*sqlite.DB, error) {
	sqlitex.NewFile(nil)
	sqlite.OpenConn()
}
