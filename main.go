package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	companyHandler "github.com/Zopsmart-Training/student-recruitment-system/delivery/company"
	studentHandler "github.com/Zopsmart-Training/student-recruitment-system/delivery/student"
	"github.com/Zopsmart-Training/student-recruitment-system/driver"
	companyService "github.com/Zopsmart-Training/student-recruitment-system/service/company"
	studentService "github.com/Zopsmart-Training/student-recruitment-system/service/student"
	"github.com/Zopsmart-Training/student-recruitment-system/store/company"
	"github.com/Zopsmart-Training/student-recruitment-system/store/student"
)

func main() {
	db, err := driver.DBConnection("mysql", "root:Aditi#2#@tcp(127.0.0.1:3306)/placement")

	if err != nil {
		log.Println(err)
		return
	}

	companyStore := company.New(db)
	studentStore := student.New(db)

	svcCmp := companyService.New(companyStore)
	svcStu := studentService.New(studentStore)

	cmpHandler := companyHandler.New(svcCmp)
	stuHandler := studentHandler.New(svcStu)

	router := mux.NewRouter()
	router.HandleFunc("/companies", middleware(cmpHandler.Get)).Methods("GET")
	router.HandleFunc("/companies/{id}", middleware(cmpHandler.GetByID)).Methods("GET")
	router.HandleFunc("/companies", middleware(cmpHandler.Create)).Methods("POST")
	router.HandleFunc("/companies/{id}", middleware(cmpHandler.Update)).Methods("PUT")
	router.HandleFunc("/companies/{id}", middleware(cmpHandler.Delete)).Methods("DELETE")

	router.HandleFunc("/students", middleware(stuHandler.Get)).Methods("GET")
	router.HandleFunc("/students/{id}", middleware(stuHandler.GetByID)).Methods("GET")
	router.HandleFunc("/students", middleware(stuHandler.Create)).Methods("POST")
	router.HandleFunc("/students/{id}", middleware(stuHandler.Update)).Methods("PUT")
	router.HandleFunc("/students/{id}", middleware(stuHandler.Delete)).Methods("DELETE")

	const timeoutVar = 3

	server := &http.Server{
		Addr:              ":8000",
		ReadHeaderTimeout: timeoutVar * time.Second,
		Handler:           router,
	}

	fmt.Println("server at port 8000")
	log.Fatal(server.ListenAndServe())
}
