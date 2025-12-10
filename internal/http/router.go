package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type PingResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// fake in-memory storage
var students = []Student{}
var idCounter = 1

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/students", studentsHandler)

	return mux
}

func pingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
	res := PingResponse{
		Message: "pong",
		Status:  "success",
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprint(w, "pong")

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
	}

}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createStudent(w, r)
		return
	case http.MethodGet:
		listStudents(w, r)
		return
	case http.MethodPut:
		updateStudent(w, r)
		return
	case http.MethodDelete:
		deleteStudent(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	if input.Name == "" || input.Age == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("all fields are required"))
		return
	}

	stud := Student{
		ID:   idCounter,
		Name: input.Name,
		Age:  input.Age,
	}
	idCounter++

	students = append(students, stud)

	res := APIResponse{
		Data: stud,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	errs := json.NewEncoder(w).Encode(res)
	if errs != nil {
		log.Fatal(errs)
	}
}

// func listStudents(w http.ResponseWriter) {
// 	res := APIResponse{
// 		Data: students,
// 	}
// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	err := json.NewEncoder(w).Encode(res)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func listStudents(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	idParam := qry.Get("id")

	if idParam == "" {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(APIResponse{Data: students})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(APIResponse{Error: "invalid id"})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, s := range students {
		if s.ID == id {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(APIResponse{Data: s})
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	errs := json.NewEncoder(w).Encode(APIResponse{Error: "student not found"})
	if errs != nil {
		log.Fatal(errs)
	}
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	idParam := qry.Get("id")

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "id required"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(APIResponse{Error: "invalid id"})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	var input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	errs := json.NewDecoder(r.Body).Decode(&input)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(APIResponse{Error: "invalid json"})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if input.Name == "" || input.Age == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(APIResponse{Error: "all fields are required"})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for i, s := range students {
		if s.ID == id {
			students[i].Name = input.Name
			students[i].Age = input.Age
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(APIResponse{Data: students[i]})
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	ers := json.NewEncoder(w).Encode(APIResponse{Error: "student not found"})
	if ers != nil {
		log.Fatal(ers)
	}

}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	idParam := qry.Get("id")

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "id required"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(APIResponse{Error: "invalid id"})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for i, s := range students {
		if s.ID == id {
			students = append(students[:i], students[i+1:]...)
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(APIResponse{Data: fmt.Sprintf("student deleted %d", id)})
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	ers := json.NewEncoder(w).Encode(APIResponse{Error: "student not found"})
	if ers != nil {
		log.Fatal(ers)
	}

}
