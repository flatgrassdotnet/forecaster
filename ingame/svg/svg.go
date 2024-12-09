/*
	forecaster - cloudbox frontend
	Copyright (C) 2024  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package svg

import (
	_ "embed"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/flatgrassdotnet/forecaster/utils"
)

type SVG struct {
	Name   string
	Fill   string
	Stroke string
}

var (
	t       = template.Must(template.New("svgbase.svg").Funcs(template.FuncMap{"StripHTTPS": func(url string) string { s, _ := strings.CutPrefix(url, "https:"); return s }}).ParseGlob("data/templates/svg/*.svg"))
	isColor = regexp.MustCompile(`^[0-9a-fA-F]{6}$`).MatchString
)

func Handle(w http.ResponseWriter, r *http.Request) {
	fill := "none"
	if isColor(r.URL.Query().Get("fill")) {
		fill = "#" + r.URL.Query().Get("fill")
	}

	stroke := "none"
	if isColor(r.URL.Query().Get("stroke")) {
		stroke = "#" + r.URL.Query().Get("stroke")
	}

	w.Header().Set("Content-Type", "image/svg+xml")

	err := t.Execute(w, SVG{
		Name:   r.PathValue("id"),
		Fill:   fill,
		Stroke: stroke,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
