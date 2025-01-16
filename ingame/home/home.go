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

package home

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/flatgrassdotnet/forecaster/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

type Home struct {
	InGame       bool
	GM13         bool
	LoggedIn     bool
	HomePage     bool
	PageType     string
	IsDarkMode   bool
	MapName      string
	Search       string
	Sort         string
	Category     string
	Packages     []common.Package
	PrevLink     string
	NextLink     string
	News         []common.NewsEntry
	PopularENTs  []common.Package
	PopularSWEPs []common.Package
	PopularMaps  []common.Package
}

const itemsPerPage = 50

var (
	pagetypes = map[string]string{
		"home":   "home",
		"news":   "news",
		"info":   "info",
		"search": "search",
		"error":  "error",
	}
	t = template.Must(template.New("browser.html").Funcs(template.FuncMap{"StripHTTPS": func(url string) string { s, _ := strings.CutPrefix(url, "https:"); return s }}).ParseGlob("data/templates/browser/*.html"))
)

func Handle(w http.ResponseWriter, r *http.Request) {
	pagetype, ok := pagetypes[r.PathValue("pagetype")]
	if !ok {
		http.Error(w, "unknown page", http.StatusNotFound)
		return
	}

	var steamid []byte
	if r.Header.Get("TICKET") != "" {
		v := make(url.Values)
		v.Set("ticket", r.Header.Get("TICKET"))

		resp, err := http.Get(fmt.Sprintf("%s/auth/getid?%s", os.Getenv("API_URL"), v.Encode()))
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

	var darkmode bool
	darkCookie, err := r.Cookie("darkmode")
	if err == nil {
		darkmode = darkCookie.Value == "true"
	}

	v := make(url.Values)
	v.Set("count", "15")
	v.Set("sort", "popular")

	if r.Host == "safe.cl0udb0x.com" {
		v.Set("safemode", "true")
	}

	var ents []common.Package
	var sweps []common.Package
	var maps []common.Package
	if pagetype == "home" {
		v.Set("type", "entity")
		resp, err := http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&ents)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
			return
		}

		v.Set("type", "weapon")
		resp, err = http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&sweps)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
			return
		}

		v.Set("type", "map")
		resp, err = http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&maps)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
			return
		}
	}

	var news []common.NewsEntry
	if pagetype == "home" || pagetype == "news" {
		resp, err := http.Get(fmt.Sprintf("%s/news/list", os.Getenv("API_URL")))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get news list: %s", err))
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&news)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode news list: %s", err))
			return
		}

		slices.Reverse(news)
	}

	ingame := strings.Contains(r.UserAgent(), "Valve") || r.Host == "toybox.garrysmod.com" || r.Host == "ingame.cl0udb0x.com" || r.Host == "safe.cl0udb0x.com"

	err = t.Execute(w, Home{
		InGame:       ingame,
		GM13:         ingame && r.Host != "toybox.garrysmod.com",
		LoggedIn:     steamid != nil,
		HomePage:     true,
		PageType:     pagetype,
		IsDarkMode:   darkmode,
		Search:       r.URL.Query().Get("search"),
		News:         news,
		PopularENTs:  ents,
		PopularSWEPs: sweps,
		PopularMaps:  maps,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
