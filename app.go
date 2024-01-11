package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/angelofallars/htmx-go"
	"html/template"
	"io"
	"net/http"
	"strings"
)

func main() {

	validierungsHandler := func(w http.ResponseWriter, r *http.Request) {
		if !htmx.IsHTMX(r) {
			fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
			return
		}

		// create new context for cloned request
		ctx := context.WithValue(context.Background(), "key", "value")

		// clone request to access it freely
		requestClone := r.Clone(ctx) // we need a context here

		// read cloned request // byte array b
		// request.Body.src.R.buf -> []uint8 => string
		b, _ := io.ReadAll(requestClone.Body)

		tmp := strings.Split(string(b), "=")

		tmp2 := tmp[1]

		tmp2 = strings.ReplaceAll(tmp2, "%5Cn", "\n")
		tmp2 = strings.ReplaceAll(tmp2, "%20", " ")
		tmp2 = strings.ReplaceAll(tmp2, "%3F", "?")
		tmp2 = strings.ReplaceAll(tmp2, "%2C", ",")

		alleFilme := strings.SplitN(tmp2, "\n", -1)

		// write alleFilme to database
		const database string = "movies.sqlite"
		db, err := sql.Open("sqlite3", database)

		const createTable string = `CREATE TABLE IF NOT EXISTS movies (
									id INTEGER NOT NULL PRIMARY KEY,
									name TEXT NOT NULL );`
	}

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/glücksrad", startPageHandler)
	multiplexer.HandleFunc("/filmselektion", filmSelektionHandler)
	multiplexer.HandleFunc("/validierung", validierungsHandler)

	http.ListenAndServe(":8080", multiplexer)

}

func validierungsHandler(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
		return
	}

}

func filmSelektionHandler(w http.ResponseWriter, request *http.Request) {
	if !htmx.IsHTMX(request) {
		fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
		return
	}

	dat := template.Must(template.ParseFiles("C:\\Users\\maria\\GolandProjects\\awesomeProject\\routes\\filmselektion.html"))
	dat.Execute(w, nil)
}

func startPageHandler(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		// fmt.Fprintf(w, "hello World")
		http.ServeFile(w, r, "routes/start.html")
	}

	if htmx.IsBoosted(r) {
		// logic for boosted
	}

	// Basic usage, page refresh
	writer := htmx.NewResponse().Refresh(true)
	writer.Write(w)

	// RETARGET response
	htmx.NewResponse().
		Reswap(htmx.SwapBeforeEnd). // fine tune swap behavior
		Retarget("#errors").
		ReplaceURL("/errors").
		Write(w)
}
