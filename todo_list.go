package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// A Todo represents a task in our todo list
type Todo struct {
	ID     int
	Task   string
	Status bool
}

// This slice will act as our in-memory database for now
var todoList = []Todo{
	{ID: 1, Task: "Hi, Go", Status: false},
	{ID: 2, Task: "Im go", Status: false},
}

type statusBool struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

func showTodoList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, todoList)
}

func addTodoForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(386)
	if err != nil {
		http.Error(w, "Form is empty", http.StatusBadRequest)
		return
	}

	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("%s\n", task)

	// Add the new task to our in-memory database
	newID := len(todoList) + 1
	todoAdd := Todo{ID: newID, Task: task, Status: false}
	todoList = append(todoList, todoAdd)

	j, _ := json.Marshal(todoAdd)
	w.Write(j)
}

func statusUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields
	s := &statusBool{}
	err := d.Decode(s)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update status
	todoList[s.ID-1].Status = s.Status

	fmt.Println(s.ID-1, s.Status)
	j, _ := json.Marshal(s)
	w.Write(j)
}

func main() {
	http.HandleFunc("/", showTodoList)
	http.HandleFunc("/status_update", statusUpdate)
	http.HandleFunc("/add_form", addTodoForm)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
