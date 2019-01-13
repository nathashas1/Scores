
## To get the project running:

1. Go to $GOPATH
2. Run go get github.com/gorilla/mux (library used to implement rest endpoints)
3. Run go get github.com/r3labs/sse (library used to listen to SSE)
4. Clone the repo https://github.com/nathashas1/StudentScores.git
5. Copy 'scores' directory from StudentScores/ to current directory
6. Run go install scores
7. Go to $GOPATH/bin
8. Run ./scores


### Testing Restful Routes:

1. To get list of all users that have received at least one test score
 Navigate to http://localhost:9000/students

2. To get list of the test results for the specified student, and provides the student's average score across all exams
 Navigate to http://localhost:9000/students/{id}

3. To get lists of all the exams that have been recorded
 Navigate to http://localhost:9000/exams

4. To lists all the results for the specified exam, and provides the average score across all students
 Navigate to http://localhost:9000/exams/{id}
