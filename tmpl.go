// Copyright (C) 2018-2019 Kdevb0x Ltd.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package vidchat

import (
	"html/template"
	"os"
	"html"
)

var _ html.EscapeString("")

var _ template.New("quikvidchat")

func loadTemplateFile(file string) (template.Template, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	template.ParseFiles(f)
}
