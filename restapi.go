package main

import (
	"context"
	"os"
    "os/signal"
    "time"
	"flag"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)


type ACT struct {
	First string `json:"First"`
	Last  string `json:"last"`
}

type Movie struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Actor    []ACT  `json:"actor"`
}


var movies []Movie
var idCount = 1

func ReplaceMovie(id string, mv Movie){
	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			movies = append(movies, mv)
			return
		}
	}
}

func DeleteMv(id string){
	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			return
		}
	}
}


func getAllActors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item.Actor)
			return
		}
	}
}

func postActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			var act ACT
			json.NewDecoder(r.Body).Decode(&act)
			item.Actor = append(item.Actor,act)

			movies = append(movies[:index], movies[index+1:]...)
			movies = append(movies, item)

			json.NewEncoder(w).Encode(item)
			return
		}
	}
}


func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}



func postMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err!= nil{
		log.Fatal(err)
	}
	idCount = idCount+1
	movie.ID = strconv.Itoa(idCount) 
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}


func putMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err!= nil{
		log.Fatal(err)
	}
	movie.ID = params["id"]
	ReplaceMovie(movie.ID,movie)
	json.NewEncoder(w).Encode(movie)

}


func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	DeleteMv(params["id"])
	json.NewEncoder(w).Encode(movies)
}

func deleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movies = movies[:0]
	json.NewEncoder(w).Encode(movies)
}




func main() {
	// router initialization
	r := mux.NewRouter()

	var act []ACT

	act = append(act, ACT{ First : "Leo", Last : "DeCaprio" } )
	act = append(act, ACT{ First : "abc", Last : "def" } )

	movies = append(movies, Movie{ ID : "1", Name : "Titanic", Actor : act } )


	r.HandleFunc("/movies", getMovie).Methods("GET")
	r.HandleFunc("/movies/{id}/actors", getAllActors).Methods("GET")
	r.HandleFunc("/movies", deleteAllMovie).Methods("DELETE")
	r.HandleFunc("/movies", postMovie).Methods("POST")
	r.HandleFunc("/movies/{id}/actors", postActor).Methods("POST")
	r.HandleFunc("/movies/{id}", putMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	//log.Fatal(http.ListenAndServe(":8000", r))


	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()
	srv := &http.Server{
        Addr:         "0.0.0.0:8080",
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

    // Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}
