package main

import (
	"database/sql"
	"net/http"
	"simple-api/api"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_database sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	api.SetupAPIRoutes(router, db)

	http.ListenAndServe(":8080", router)
}
