package student

import(
	// "fmt"
)
type Student struct {
    Exam      	 int
    StudentId    string
    Score        float32
}

func (stud Student) GetScore() {
	// fmt.Printf("\nStudent id=%s  Exam Id=%d  Score=%f",stud.StudentId , stud.Exam,stud.Score)
}
