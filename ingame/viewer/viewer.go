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

package viewer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/flatgrassdotnet/forecaster/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

var t = template.Must(template.New("viewer.html").ParseFiles("data/templates/viewer/viewer.html"))

type Viewer struct {
	Package  common.Package
	Category string
}

func Handle(w http.ResponseWriter, r *http.Request) {
	v := make(url.Values)
	v.Set("id", r.PathValue("id"))

	resp, err := http.Get(fmt.Sprintf("%s/packages/get?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
		return
	}

	var pkg common.Package
	err = json.NewDecoder(resp.Body).Decode(&pkg)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
		return
	}

	err = t.Execute(w, Viewer{
		Package:  pkg,
		Category: r.URL.Query().Get("show"),
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
