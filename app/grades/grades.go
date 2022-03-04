package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, g := range s.Grades {
		result += g.Score
	}
	return result / float32(len(s.Grades))
}

type GradeType string
type Students []Student

func (s Students) GetById(id int) (*Student, error) {
	for i := range s {
		if s[i].ID == id {
			return &s[i], nil
		}
	}
	return nil, fmt.Errorf("No student found with id %v", id)
}

var (
	students      Students
	studentsMutex sync.Mutex
)

const (
	GradeQuiz     = GradeType("Quiz")
	GradeTest     = GradeType("Test")
	GradeHomework = GradeType("Homework")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
