
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
var studentAvgs = make(map[string]float32)
var examAvgs = make(map[string]float32)
var countExams = make(map[string]int)
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
			storeStudents(stud);
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
			 var keys []int
			var allExams = make(map[int]int)
			for _, studArray := range students {
				for _,studObj := range studArray {
					allExams[studObj.Exam] = 1
				}
			}
			for k := range allExams {
        keys = append(keys, k)
    }
	    json.NewEncoder(w).Encode(keys)
}


func GetExam(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
      exam := params["number"]
			sum,ok := examAvgs[exam]
			 fmt.Println(sum)
			 fmt.Println(countExams[exam])
			if ok {
					avg  := sum/float32(countExams[exam])
      	  json.NewEncoder(w).Encode(avg)
			}else{
				  json.NewEncoder(w).Encode(exam)
			}
}



func GetStudent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
		fmt.Printf("%s\n", params["id"])
      name := params["id"]
			sum,ok := studentAvgs[name]
			if ok {
					allMarks,_ := students[name]
					avg  := sum/float32(len(allMarks))
					  fmt.Printf("length= %f\n", float32(len(allMarks)))
						fmt.Printf("sum= %f\n", sum)
      	  json.NewEncoder(w).Encode(avg)
			}else{
				  json.NewEncoder(w).Encode("Student Not found")
			}
}


func storeStudents(stud student.Student){
			studList, ok := students[stud.StudentId]
			if ok {
					studList =append(studList,stud)
					for _, studObj := range studList {
							studObj.GetScore()
					}
					students[stud.StudentId]=studList
					studentAvgs[stud.StudentId] +=stud.Score;
					examAvgs[fmt.Sprint(stud.Exam)] +=stud.Score;
					countExams[fmt.Sprint(stud.Exam)] +=1;
			} else {
					 var studentList []student.Student
					 studentList=append(studentList,stud)
					 students[stud.StudentId]=studentList
					 studentAvgs[stud.StudentId] = stud.Score
					 examAvgs[fmt.Sprint(stud.Exam)] = stud.Score
					 countExams[fmt.Sprint(stud.Exam)] =1;
			}
}
