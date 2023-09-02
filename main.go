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
	db, err := driver.DBConnection("mysql", "root:password@tcp(172.17.0.2:3306)/placement")

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
	router.HandleFunc("/companies", cmpHandler.Get).Methods("GET")
	router.HandleFunc("/companies/{id}", cmpHandler.GetByID).Methods("GET")
	router.HandleFunc("/companies", cmpHandler.Create).Methods("POST")
	router.HandleFunc("/companies/{id}", cmpHandler.Update).Methods("PUT")
	router.HandleFunc("/companies/{id}", cmpHandler.Delete).Methods("DELETE")

	router.HandleFunc("/students", stuHandler.Get).Methods("GET")
	router.HandleFunc("/students/{id}", stuHandler.GetByID).Methods("GET")
	router.HandleFunc("/students", stuHandler.Create).Methods("POST")
	router.HandleFunc("/students/{id}", stuHandler.Update).Methods("PUT")
	router.HandleFunc("/students/{id}", stuHandler.Delete).Methods("DELETE")

	const timeoutVar = 3

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: timeoutVar * time.Second,
		Handler:           router,
	}

	fmt.Println("server at port 8080")
	log.Fatal(server.ListenAndServe())
}
