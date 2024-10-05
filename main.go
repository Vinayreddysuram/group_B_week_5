package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type Movie struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Director    string `json:"director"`
    ReleaseYear int    `json:"release_year"`
    Genre       string `json:"genre"`
}

var movies []Movie
var nextID int = 1

// Create a new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var movie Movie
    json.NewDecoder(r.Body).Decode(&movie)
    movie.ID = nextID
    nextID++
    movies = append(movies, movie)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movie)
}

// Get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies)
}

// Get a single movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id, err := extractID(r.URL.Path)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }

    for _, movie := range movies {
        if movie.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(movie)
            return
        }
    }
    http.Error(w, "Movie not found", http.StatusNotFound)
}

// Update an existing movie by ID
func updateMovie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id, err := extractID(r.URL.Path)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }

    for i, movie := range movies {
        if movie.ID == id {
            json.NewDecoder(r.Body).Decode(&movie)
            movies[i] = movie
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(movie)
            return
        }
    }
    http.Error(w, "Movie not found", http.StatusNotFound)
}

// Delete a movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id, err := extractID(r.URL.Path)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }

    for i, movie := range movies {
        if movie.ID == id {
            movies = append(movies[:i], movies[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Movie not found", http.StatusNotFound)
}

// Extract ID from URL path (e.g., /movies/1)
func extractID(path string) (int, error) {
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        return 0, fmt.Errorf("invalid path")
    }
    return strconv.Atoi(parts[2])
}

func main() {
    http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getMovies(w, r)
        case http.MethodPost:
            createMovie(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getMovie(w, r)
        case http.MethodPut:
            updateMovie(w, r)
        case http.MethodDelete:
            deleteMovie(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Start the server
    fmt.Println("Server is running on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
