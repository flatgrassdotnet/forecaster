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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flatgrassdotnet/forecaster/ingame/browser"
)

func main() {
	port := flag.Int("port", 80, "web server listen port")
	flag.Parse()

	if os.Getenv("API_URL") == "" {
		os.Setenv("API_URL", "https://api.cl0udb0x.com")
	}

	// static assets
	http.Handle("GET /client/", http.StripPrefix("/client/", http.FileServer(http.Dir("data/client"))))

	// toybox.garrysmod.com
	http.HandleFunc("GET /", browser.Home)
	http.HandleFunc("GET /{type}", browser.Browser)
	http.HandleFunc("GET /search", browser.Search)

	// redundant?
	//http.HandleFunc("GET /IG/{show}", browser.Browser)

	// redirect
	http.HandleFunc("GET /ingame/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("GET /IG/maps/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/maps", http.StatusSeeOther)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error while serving: %s", err)
	}
}
