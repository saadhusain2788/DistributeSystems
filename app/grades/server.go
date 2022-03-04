package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

type studentHandler struct{}

func (s studentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	switch len(paths) {
	case 2:
		s.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(paths[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.AddGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func (sh studentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJson(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("COntent-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	data, err := sh.toJson(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("COntent-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) AddGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)

	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJson(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("COntent-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) toJson(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Faild to serialize ")
	}
	return b.Bytes(), nil

}
