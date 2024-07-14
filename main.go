package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/exp/rand"
)

type User struct {
	Id         int    `json:"id"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserStatus string `json:"user_status"`
	Department string `json:"department"`
}

type CreateUserInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Department string `json:"department"`
}

type UpdateUserInput struct {
	Id         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserStatus string `json:"user_status"`
	Department string `json:"department"`
}

func main() {
	// connect to db
	connStr := "user=user password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	handler := c.Handler(mux)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetAllUsers(w, r, db)
			break
		case "POST":
			CreateUser(w, r, db)
			break
		default:
			return
		}
	})

	mux.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			UpdateUser(w, r, db)
			break
		case "DELETE":
			DeleteUser(w, r, db)
			break
		}
	})

	log.Fatal(http.ListenAndServe(":8081", handler))
}

func GetAllUsers(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		respString(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var user User
	var users []User = make([]User, 0)
	for rows.Next() {
		rows.Scan(&user.Id, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.UserStatus, &user.Department)
		users = append(users, user)
	}

	fmt.Println(users)

	respJSON(w, users, http.StatusOK)
}

func CreateUser(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var user CreateUserInput
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("new user %#v", user)

	user_name := fmt.Sprintf("%s%s%s", strings.ToLower(user.FirstName), strings.ToLower(user.LastName), generateUserNameSuffix())
	email := fmt.Sprintf("%s@integra.com", user_name)

	_, err = db.Exec(`INSERT INTO users(user_name, first_name, last_name, email, user_status, department)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user_name, user.FirstName, user.LastName, email, "A", user.Department,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respSuccess(w)
}

func UpdateUser(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil && id > 0 {
		respString(w, "invalid id given; must be a positive integer", http.StatusBadRequest)
		return
	}

	var user UpdateUserInput
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("USER: %#v \n", user)

	res, err := db.Exec(`UPDATE users 
				SET first_name = $2, last_name = $3, 
					email = $4, user_status = $5, department = $6 
				WHERE id = $1`,
		id, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nrows, _ := res.RowsAffected(); nrows == 0 {
		respString(w, "user does not exist", http.StatusNotFound)
		return
	}

	respSuccess(w)
}

func DeleteUser(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rawId := req.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil && id > 0 {
		respString(w, "invalid id given; must be a positive integer", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nrows, _ := res.RowsAffected(); nrows == 0 {
		respString(w, "user does not exist", http.StatusNotFound)
		return
	}

	respSuccess(w)
}

func wrapWithDB(handle func(w http.ResponseWriter, r *http.Request, db *sql.DB), db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, db)
	}
}

func respJSON(w http.ResponseWriter, value any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func respString(w http.ResponseWriter, value string, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	io.WriteString(w, value)
}

func respSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "")
}

func generateUserNameSuffix() string {
	const LENGTH = 6
	const charset = "0123456789"

	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	random := rand.New(seed)

	result := make([]byte, LENGTH)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}
