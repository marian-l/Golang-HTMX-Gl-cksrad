package main

import (
	"fmt"
	"github.com/angelofallars/htmx-go"
	"html/template"
	"net/http"
)

func main() {

	filmSelektionHandler := func(w http.ResponseWriter, request *http.Request) {
		if !htmx.IsHTMX(request) {
			fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
			return
		}

		dat := template.Must(template.ParseFiles("C:\\Users\\maria\\GolandProjects\\awesomeProject\\routes\\filmselektion.html"))
		dat.Execute(w, nil)
	}

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/glücksrad", startPageHandler)
	multiplexer.HandleFunc("/filmselektion", filmSelektionHandler)

	http.ListenAndServe(":8080", multiplexer)

}

func filmSelektionHandler(w http.ResponseWriter, request *http.Request) {
	if !htmx.IsHTMX(request) {
		fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
		return
	}

	dat := template.Must(template.ParseFiles("C:\\Users\\maria\\GolandProjects\\awesomeProject\\routes\\filmselektion.html"))
	dat.Execute(w, nil)

	fmt.Fprintf(w, "HTMX-Anfrage")
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
