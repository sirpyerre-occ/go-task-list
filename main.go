package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:name`
	Content string `json:content`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task 1",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been delete succesfully", taskID)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Welcome to my API")
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	var updateTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}

	json.Unmarshal(reqBody, &updateTask)

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]... )
			updateTask.ID = taskID
			tasks = append(tasks, updateTask)

			fmt.Fprintf(w, "The task with ID %q has been updated successfully", taskID)
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid task")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
