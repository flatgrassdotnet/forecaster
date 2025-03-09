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

package publishsave

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/flatgrassdotnet/forecaster/utils"
)

var tp = template.Must(template.New("publish.html").ParseFiles("data/templates/publishsave/publish.html"))

func Publish(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse form data: %s", err))
		return
	}

	v := make(url.Values)
	v.Set("id", r.URL.Query().Get("id"))
	v.Set("sid", r.URL.Query().Get("sid"))
	v.Set("name", r.PostForm.Get("name"))
	v.Set("desc", r.PostForm.Get("desc"))
	v.Set("ticket", r.Header.Get("TICKET"))

	resp, err := http.Get(fmt.Sprintf("%s/package/publishsave?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to publish save: %s", err))
		return
	}

	if resp.StatusCode != http.StatusOK {
		utils.WriteError(w, r, fmt.Sprintf("failed to publish save: got code %d", resp.StatusCode))
		return
	}

	err = tp.Execute(w, nil)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
