package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	companyHandler "github.com/aditi-zs/Placement-API/delivery/company"
	studentHandler "github.com/aditi-zs/Placement-API/delivery/student"
	"github.com/aditi-zs/Placement-API/driver"
	companyService "github.com/aditi-zs/Placement-API/service/company"
	studentService "github.com/aditi-zs/Placement-API/service/student"
	"github.com/aditi-zs/Placement-API/store/company"
	"github.com/aditi-zs/Placement-API/store/student"
)

func main() {
	db, err := driver.DBConnection("mysql", "root:password@tcp(localhost:3306)/placement")

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
