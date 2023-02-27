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

// Global variable for passenger
var passengers Passenger
var db *sql.DB

type Passenger struct { // map this type to the record in the table
	PassengerID  string
	FirstName    string
	LastName     string
	MobileNumber string
	EmailAddress string
}

//																		  //
//	Functions for MySQL Database										  //
//																		  //
//
// //////////////////////////////////////////////////////////////////////////////////////////////////
// Registering new passenger
func InitNewPassenger(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("INSERT INTO Passengers VALUES ('%s', '%s', '%s', '%s', '%s')",
		p.PassengerID, p.FirstName, p.LastName, p.MobileNumber, p.EmailAddress)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

// Updating existing passenger information
func UpdatePassenger(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("UPDATE Passengers SET FirstName='%s', LastName='%s', MobileNumber='%s', EmailAddress='%s' WHERE PassengerID='%s'",
		p.FirstName, p.LastName, p.MobileNumber, p.EmailAddress, p.PassengerID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// Passenger using mobile phone number to login
func PassengerLogin(db *sql.DB, MobileNumber string) (Passenger, string) {
	query := fmt.Sprintf("SELECT * FROM Passengers WHERE MobileNumber = '%s'", MobileNumber)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&passengers.PassengerID, &passengers.FirstName, &passengers.LastName, &passengers.MobileNumber, &passengers.EmailAddress); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Passenger login failed."
	case nil:
	default:
		panic(err.Error())
	}

	return passengers, errMsg
}

// 																								//
// 									Functions for HTTP request and response						//
//////////////////////////////////////////////////////////////////////////////////////////////////

func passengerInit(resp http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported Media Type"))
		return
	}

	var newPassenger Passenger
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	err = json.Unmarshal(reqBody, &newPassenger)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	if newPassenger.PassengerID == "" || newPassenger.FirstName == "" || newPassenger.LastName == "" || newPassenger.MobileNumber == "" || newPassenger.EmailAddress == "" {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		resp.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
		return
	}
	InitNewPassenger(db, newPassenger)
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("201 - Successfully created passenger's information"))
}

func passengerUpdate(resp http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported Media Type"))
		return
	}

	var passenger Passenger
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	err = json.Unmarshal(reqBody, &passenger)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("400 - Bad Request"))
		return
	}
	UpdatePassenger(db, passenger)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("200 - Successfully updated passenger's information"))
}

func passengerRetrieve(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	MobileNumber := params["mobile_no"]
	passenger, _ := PassengerLogin(db, MobileNumber)
	if passenger.MobileNumber == "" {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("404 - Passenger not found"))
		return
	}
	json.NewEncoder(resp).Encode(passenger)
}

func passengerDelete(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusForbidden)
	resp.Write([]byte("403 - For audit purposes, passenger's account cannot be deleted."))
}

func passenger(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		passengerInit(resp, req)
	case "PUT":
		passengerUpdate(resp, req)
	case "GET":
		passengerRetrieve(resp, req)
	case "DELETE":
		passengerDelete(resp, req)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		resp.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	// instantiate passengers
	hailgrab_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hailgrab_db")

	db = hailgrab_db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/passengers", passenger).Methods("GET", "POST", "PUT", "DELETE")
	fmt.Println("Passenger --> Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))

	defer db.Close()
}
