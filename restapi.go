package main

import (
	"encoding/json"
	"log"
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
var idCount = 2

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

	// creating dummy data
	movies = append(movies, Movie{ID: "1", Name: "Titanic", Director : "cameroon"})
	movies = append(movies, Movie{ID: "2", Name: "Avengers", Director : "unknown"})

	
	r.HandleFunc("/movies", getMovie).Methods("GET")
	r.HandleFunc("/movies", deleteAllMovie).Methods("DELETE")
	r.HandleFunc("/movies", postMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", putMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
