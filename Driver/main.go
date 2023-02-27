package ETIAssignment1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for driver
var drivers Driver
var db *sql.DB

type Driver struct {
	DriverID     string
	FirstName    string
	LastName     string
	MobileNumber string
	EmailAddress string
	CarPlateNo   string
	isAvailable  bool
}

//																//
//	Functions for Database										//
//
// ///////////////////////////////////////////////////////////////////////////////////////
// Registering new driver
func InitNewDriver(db *sql.DB, d Driver) {
	query := fmt.Sprintf("INSERT INTO Drivers VALUES ('%s', '%s', '%s','%s', '%s', '%s', %t)",
		d.DriverID, d.FirstName, d.LastName, d.MobileNumber, d.EmailAddress, d.CarPlateNo, d.isAvailable)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

// Updating existing driver information
func UpdateDriver(db *sql.DB, d Driver) {
	fmt.Println(d)
	query := fmt.Sprintf("UPDATE Drivers SET FirstName='%s', LastName='%s', MobileNumber='%s', EmailAddress='%s', CarPlateNo='%s', isAvailable=%t WHERE DriverID='%s'",
		d.FirstName, d.LastName, d.MobileNumber, d.EmailAddress, d.CarPlateNo, d.isAvailable, d.DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// Driver using mobile phone number to login
func LoginDriverAccount(db *sql.DB, MobileNumber string) (Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Drivers WHERE MobileNumber = '%s'", MobileNumber)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&drivers.DriverID, &drivers.FirstName, &drivers.LastName, &drivers.MobileNumber, &drivers.EmailAddress, &drivers.CarPlateNo, &drivers.isAvailable); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Driver login failed."
	case nil:
	default:
		panic(err.Error())
	}

	return drivers, errMsg
}

//																								//
//									Functions for HTTP request and response						//
//////////////////////////////////////////////////////////////////////////////////////////////////

func driverInit(resp http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported Media Type"))
		return
	}

	var newDriver Driver
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	err = json.Unmarshal(reqBody, &newDriver)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	if newDriver.DriverID == "" || newDriver.FirstName == "" || newDriver.LastName == "" || newDriver.MobileNumber == "" || newDriver.EmailAddress == "" || newDriver.CarPlateNo == "" {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		resp.Write([]byte("422 - Please supply driver " + "information " + "in JSON format"))
		return
	}
	InitNewDriver(db, newDriver)
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("201 - Successfully created driver's information"))
}

func driverUpdate(resp http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported Media Type"))
		return
	}

	var driver Driver
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	err = json.Unmarshal(reqBody, &driver)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	UpdateDriver(db, driver)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("200 - Successfully updated driver's information"))
}

func driverRetrieve(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	MobileNumber := params["mobile_no"]
	driver, _ := LoginDriverAccount(db, MobileNumber)
	if driver.MobileNumber == "" {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("404 - Driver not found"))
		return
	}
	json.NewEncoder(resp).Encode(driver)
}

func driverDelete(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusForbidden)
	resp.Write([]byte("403 - For audit purposes, driver's account cannot be deleted."))
}

func driver(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		driverInit(resp, req)
	case "PUT":
		driverUpdate(resp, req)
	case "GET":
		driverRetrieve(resp, req)
	case "DELETE":
		driverDelete(resp, req)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		resp.Write([]byte("405 - Method Not Allowed"))
	}
}

// Assign driver to trip
func main() {
	// instantiate drivers
	hailgrab_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hailgrab_db")

	db = hailgrab_db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/drivers", driver).Methods("GET", "PUT", "POST", "DELETE")
	fmt.Println("Driver --> Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

	defer db.Close()
}
