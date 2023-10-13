package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"simple-api/model"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupAPIRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/tasks", CreateTask(db)).Methods("POST")
	router.HandleFunc("/tasks", GetTasks(db)).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", GetTask(db)).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", UpdateTask(db)).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", DeleteTask(db)).Methods("DELETE")
}

func CreateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO task (title, isCompleted) VALUES ($1, $2)", task.Title, task.IsCompleted)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	}
}

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, isCompleted FROM task")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []model.Task
		for rows.Next() {
			var task model.Task
			if err := rows.Scan(&task.ID, &task.Title, &task.IsCompleted); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		taskID, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var task model.Task
		err = db.QueryRow("SELECT id, title, isCompleted FROM task WHERE id = $1", taskID).Scan(&task.ID, &task.Title, &task.IsCompleted)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		taskID, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var updatedTask model.Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE task SET title = $1, isCompleted = $2 WHERE id = $3", updatedTask.Title, updatedTask.IsCompleted, taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedTask)
	}
}

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		taskID, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM task WHERE id = $1", taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
