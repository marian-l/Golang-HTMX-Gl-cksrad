package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/angelofallars/htmx-go"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var db, _ = Connect()
var (
	name string
)

func main() {
	setupLogging()

	// setup server routes
	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/glücksrad", startPageHandler)
	multiplexer.HandleFunc("/filmselektion", filmSelektionHandler)
	multiplexer.HandleFunc("/validierung", validierungsHandler)

	// setup server
	err := http.ListenAndServe(":8080", multiplexer)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLogging() {
	// setup logging
	file, err := os.OpenFile("errorLog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func selectMoviesDB() (movieList []string) {
	// get movie list
	query := "SELECT movies.name from main.movies;"
	rows, err := db.Query(query)

	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	alleFilme := []string{}

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		err = rows.Err()
		alleFilme = append(alleFilme, name)
	}

	if err != nil {
		log.Fatal(err)
	}

	return alleFilme
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database/movies.sqlite")
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return db, err
}

func writeToDatabase(filme []string) {
	// write alleFilme to database

	query := "BEGIN TRANSACTION;"

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("locking failed: %s", err)
	}

	query = "INSERT INTO 'movies' (name) VALUES (?)"

	for _, movie := range filme {

		insertResult, err := db.ExecContext(context.Background(), query, movie)

		if err != nil {
			log.Fatalf("insert failed: %s", err)
		}

		id, err := insertResult.LastInsertId()
		if err != nil {
			log.Fatalf("impossible to retrieve last inserted id: %s", err)
		}
		log.Printf("inserted id: %d", id)
	}

	query = "COMMIT;"

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("commiting failed: %s", err)
	}
}

func processRequest(clone *http.Request) []string {
	b, err := io.ReadAll(clone.Body)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	tmp := strings.Split(string(b), "=")

	tmp2 := tmp[1]

	tmp2 = strings.ReplaceAll(tmp2, "%5Cn", "\n")
	tmp2 = strings.ReplaceAll(tmp2, "%20", " ")
	tmp2 = strings.ReplaceAll(tmp2, "%3F", "?")
	tmp2 = strings.ReplaceAll(tmp2, "%2C", ",")

	alleFilme := strings.SplitN(tmp2, "\n", -1)
	fmt.Println(alleFilme)

	for index, _ := range alleFilme {
		alleFilme[index] = string(alleFilme[index])

		if strings.HasPrefix(alleFilme[index], " ") {
			alleFilme[index], _ = strings.CutPrefix(alleFilme[index], " ")
		}
	}
	return alleFilme
	// request.Body.src.R.buf -> []uint8 => string
}

func setupDatabase() {
	const createTable string = `CREATE TABLE IF NOT EXISTS 'movies' (
									id INTEGER NOT NULL PRIMARY KEY,
									name TEXT NOT NULL );`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
		return
	}
}

func validierungsHandler(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
		return
	}

	// create new context for cloned request
	ctx := context.WithValue(context.Background(), "key", "value")

	// clone request to access it freely
	requestClone := r.Clone(ctx) // we need a context here

	alleFilme := processRequest(requestClone)

	writeToDatabase(alleFilme)
}

func filmSelektionHandler(w http.ResponseWriter, request *http.Request) {
	if !htmx.IsHTMX(request) {
		fmt.Fprintf(w, "keine HTMX-Anfrage, kein Content!")
		return
	}

	template := template.Must(template.ParseFiles("C:\\Users\\maria\\GolandProjects\\awesomeProject\\routes\\filmselektion.html"))
	err := template.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startPageHandler(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		// serve website
		http.ServeFile(w, r, "routes/start.html")
	}

	if htmx.IsBoosted(r) {
		// logic for boosted
	}

	// Basic usage, page refresh
	writer := htmx.NewResponse().Refresh(true)
	err := writer.Write(w)
	if err != nil {
		log.Fatal(err)
	}

	// RETARGET response
	err = htmx.NewResponse().
		Reswap(htmx.SwapBeforeEnd). // fine tune swap behavior
		Retarget("#errors").
		ReplaceURL("/errors").
		Write(w)

	if err != nil {
		log.Fatal(err)
	}

}
