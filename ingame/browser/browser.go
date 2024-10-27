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

package browser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/flatgrassdotnet/forecaster/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

type Browser struct {
	InGame   bool
	LoggedIn bool
	MapName  string
	Search   string
	Sort     string
	Category string
	Packages []common.Package
	PrevLink string
	NextLink string
}

const itemsPerPage = 50

var (
	categories = map[string]string{
		"mine":     "mine",
		"entities": "entity",
		"weapons":  "weapon",
		"props":    "prop",
		"saves":    "savemap",
		"maps":     "map",
	}
	t = template.Must(template.New("browser.html").Funcs(template.FuncMap{"StripHTTPS": func(url string) string { s, _ := strings.CutPrefix(url, "https:"); return s }}).ParseFiles("data/templates/browser/browser.html"))
)

func Handle(w http.ResponseWriter, r *http.Request) {
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	var steamid []byte
	if r.Header.Get("TICKET") != "" {
		v := make(url.Values)
		v.Set("ticket", r.Header.Get("TICKET"))

		resp, err := http.Get(fmt.Sprintf("https://api.cl0udb0x.com/auth/getid?%s", v.Encode()))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get steamid: %s", err))
			return
		}

		steamid, err = io.ReadAll(resp.Body)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to read steamid: %s", err))
			return
		}
	}

	v := make(url.Values)
	v.Set("type", category)
	v.Set("offset", strconv.Itoa((page-1)*itemsPerPage))
	v.Set("count", strconv.Itoa(itemsPerPage))
	v.Set("sort", r.URL.Query().Get("sort"))
	if category == "mine" {
		v.Del("type")
		v.Set("author", string(steamid))
	}

	if r.URL.Query().Has("search") {
		v.Set("search", r.URL.Query().Get("search"))
	}

	resp, err := http.Get(fmt.Sprintf("https://api.cl0udb0x.com/packages/list?%s", v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
		return
	}

	var list []common.Package
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
		return
	}

	prev := fmt.Sprintf("?page=%d", page-1)
	if page <= 1 {
		prev = "#"
	}

	next := fmt.Sprintf("?page=%d", page+1)
	if len(list) < itemsPerPage {
		next = "#"
	}

	err = t.Execute(w, Browser{
		InGame:   strings.Contains(r.UserAgent(), "Valve"),
		LoggedIn: steamid != nil,
		MapName:  r.Header.Get("MAP"),
		Search:   r.URL.Query().Get("search"),
		Sort:     r.URL.Query().Get("sort"),
		Category: category,
		Packages: list,
		PrevLink: prev,
		NextLink: next,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
