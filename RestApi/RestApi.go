package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)



type Movie struct {
	ID        string  `json:"id"`
	Name  	  string  `json:"name"`
	Director  string  `json:"director"`
}

var movies []Movie




func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}



func postMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100)) // Mock ID - not safe
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}


func putMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}


func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}




func main() {
	// router initialization
	r := mux.NewRouter()

	// creating dummy data
	movies = append(movies, Movie{ID: "1", Name: "Titanic", Director : "cameroon"})
	movies = append(movies, Movie{ID: "2", Name: "Avengers", Director : "unknown"})

	
	r.HandleFunc("/movies", getMovie).Methods("GET")
	r.HandleFunc("/movies", postMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", putMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
