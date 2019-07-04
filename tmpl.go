// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"html/template"
	"net/http"
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

type ClientInfoForm struct {
	DisplayName string

}

func renderTmpl(w http.ResponseWriter, tmpl template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
