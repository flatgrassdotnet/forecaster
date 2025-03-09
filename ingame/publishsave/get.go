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
	"strconv"

	"github.com/flatgrassdotnet/forecaster/utils"
)

type SaveData struct {
	ID  int
	SID int
	Map string
}

var ts = template.Must(template.New("save.html").ParseFiles("data/templates/publishsave/save.html"))

func Save(w http.ResponseWriter, r *http.Request) {
	sd := SaveData{
		Map: r.Header.Get("MAP"),
	}

	var err error
	sd.ID, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	sd.SID, err = strconv.Atoi(r.URL.Query().Get("sid"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse sid value: %s", err))
		return
	}

	err = ts.Execute(w, sd)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
