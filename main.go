package main

import(
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

type Employee struct {
	gorm.Model
	Name   string `gorm:"unique" json:"name"`
	City   string `json:"city"`
	Age    int    `json:"age"`
	Status bool   `json:"status"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/gosunshine?charset=utf8&parseTime=True")

	if err!=nil{
		log.Println("Connection Failed to Open")
	}else{
		log.Println("Connection Established")
	}

	db.AutoMigrate(Employee{})

	router := mux.NewRouter()

	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/employees", GetAll).Methods("GET")
	router.HandleFunc("/employees/{name}", GetById).Methods("GET")
	router.HandleFunc("/employees", CreateEmployee).Methods("POST")

	log.Fatal(http.ListenAndServe(":8099", router))

}

func Home(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode("Hello brow")

}

func CreateEmployee(w http.ResponseWriter, r *http.Request)  {
	employee := Employee{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	if errSave := db.Save(&employee).Error; errSave != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errSave)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(employee)

}

func GetById(w http.ResponseWriter, r *http.Request)  {
	employee := Employee{}

	//vars := mux.Vars(r)
	//name := vars["name"]

	name := r.Header.Get("name")

	if errGet := db.Where(Employee{Name: name}).First(&employee); errGet != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errGet)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(employee)
}

func GetAll(w http.ResponseWriter, r *http.Request)  {
	var employees []Employee

	if errGet := db.Find(&employees).Error; errGet != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errGet)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(employees)
}