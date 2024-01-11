package main

import (
	"context"
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
		b, _ := io.ReadAll(requestClone.Body)

		// alternative body cloning i found
		// seems bugged or sth
		body, err := r.GetBody()
		if err != nil {
			return
		}

		defer body.Close()
		buf := make([]byte, len(b))
		n, err := body.Read(buf)

		fmt.Printf("%s\n", buf[:n])

		// copy the byte array from the request somewhere

		// body := b

		// body = string(body)

		alleFilme := strings.Split(string(b), "=")

		alleFilme = strings.SplitN(alleFilme[1], "%5Cn", -1)
		fmt.Println(alleFilme)

		// currently there is no newline character inside the list
		// reading the bytes from the body, we probably need a better string conversion.

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
