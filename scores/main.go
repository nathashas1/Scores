
package main
import(
	"github.com/r3labs/sse"
	"scores/student"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
  "net/http"
	"fmt"
)
var students   = make(map[string][]student.Student)
var studentTotal = make(map[string]float32)
var exams = make(map[string][]student.Student)
var examTotal = make(map[string]float32)

type StudentResponse struct {
    Average   float32
    Marks     []student.Student
}

type ExamResponse struct {
    Average   float32
    Marks     []student.Student
}

func main() {
	go subScribeToStream()
	setUpEndPoints() ;
}

func subScribeToStream(){
	client := sse.NewClient("http://live-test-scores.herokuapp.com/scores");
	var stud student.Student
	client.Subscribe("data", func(msg *sse.Event) {
		err := json.Unmarshal([]byte(msg.Data), &stud)
		if err == nil {
			storeStudents(stud)
			storeExams(stud)
		}
	})
}

func setUpEndPoints() {
	router := mux.NewRouter()
	router.HandleFunc("/students", GetAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", GetStudent).Methods("GET")
	router.HandleFunc("/exams", GetAllExams).Methods("GET")
	router.HandleFunc("/exams/{number}", GetExam).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func GetAllStudents(w http.ResponseWriter, r *http.Request){
	    var allStudents = make([]string, len(students))
	    index := 0
		for name := range students {
			allStudents[index] = name
			index++
		}
	    json.NewEncoder(w).Encode(allStudents)
}


func GetAllExams(w http.ResponseWriter, r *http.Request){
			var allExams = make([]string, len(exams))
			index := 0
		for name := range exams {
			allExams[index] = name
			index++
		}
	    json.NewEncoder(w).Encode(allExams)
}


func GetExam(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
      number := params["number"]
			sum,ok := examTotal[number]
			if ok {
					allMarks,_ := exams[number]
					avg  := sum/float32(len(allMarks))
					resp := & ExamResponse{
						Average: avg,
						Marks :allMarks,
					}

      	  json.NewEncoder(w).Encode(resp)
			}else{
				  json.NewEncoder(w).Encode("Exam not found!")
			}
}



func GetStudent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
      name := params["id"]
			sum,ok := studentTotal[name]
			if ok {
					allMarks,_ := students[name]
				  	avg  := sum/float32(len(allMarks))
						resp := & StudentResponse{
							Average: avg,
							Marks :allMarks,
						}

      	  json.NewEncoder(w).Encode(resp)
			}else{
				  json.NewEncoder(w).Encode("Student Not found")
			}
}


func storeStudents(stud student.Student){
			studList, ok := students[stud.StudentId]
			if ok {
					studList =append(studList,stud)
					students[stud.StudentId]=studList
					studentTotal[stud.StudentId] +=stud.Score;
			} else {
					 var studentList []student.Student
					 studentList=append(studentList,stud)
					 students[stud.StudentId]=studentList
					 studentTotal[stud.StudentId] = stud.Score
			}

}

func storeExams(stud student.Student){
	  exam, ok := exams[fmt.Sprint(stud.Exam)]
		if ok{
				exam = append(exam,stud)
				exams[fmt.Sprint(stud.Exam)] = exam;
				examTotal[fmt.Sprint(stud.Exam)] +=stud.Score
		}else{
				var examList []student.Student
				examList = append(examList,stud)
				exams[fmt.Sprint(stud.Exam)] = examList
				examTotal[fmt.Sprint(stud.Exam)] = stud.Score
		}
}
