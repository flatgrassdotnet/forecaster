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

package common

import (
	"time"
)

type Package struct {
	ID          int       `json:"id"`
	Revision    int       `json:"rev"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Dataname    string    `json:"dataname,omitempty"`
	Author      string    `json:"author,omitempty"`
	AuthorName  string    `json:"authorname,omitempty"`
	AuthorIcon  string    `json:"authoricon,omitempty"`
	Description string    `json:"description,omitempty"`
	Data        []byte    `json:"data,omitempty"`
	Content     []Content `json:"content,omitempty"`
	Includes    []Include `json:"includes,omitempty"`
	Uploaded    time.Time `json:"uploaded,omitempty"`

	Downloads int `json:"downloads,omitempty"`
	Favorites int `json:"favorites,omitempty"`
	Goods     int `json:"goods,omitempty"`
	Bads      int `json:"bads,omitempty"`
}

type Content struct {
	ID       int    `json:"id"`
	Revision int    `json:"rev"`
	Path     string `json:"path"`
	Size     int    `json:"size"`
	PSize    int    `json:"psize"`
}

type Include struct {
	ID       int    `json:"id"`
	Revision int    `json:"rev"`
	Type     string `json:"type"`
}
