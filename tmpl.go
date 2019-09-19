// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"html"
	"html/template"
	"net/http"
)

var _ = html.EscapeString("")

var _ = template.New("quikface")

func loadTemplateFile(file ...string) (*template.Template, error) {
	return template.ParseFiles(file...)
}

type ClientInfoForm struct {
	DisplayName string
}

func renderTmpl(w http.ResponseWriter, tmpl template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
