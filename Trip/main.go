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

// Trip's Global variable
var trips Trip
var db *sql.DB

type Trip struct {
	TripID       string
	StatusOfTrip string
	PassengerID  string
	DriverID     string
	PickUp       string
	DropOff      string
}

// MYSQL functions
func InitNewTrip(db *sql.DB, t Trip) {
	query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%s')",
		t.TripID, t.StatusOfTrip, t.PassengerID, t.DriverID, t.PickUp, t.DropOff)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func GetTripInfo(db *sql.DB, TripID string) []Trip {
	query := fmt.Sprintf("SELECT * FROM Trips where TripID = '%s'", TripID)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() {
		var newTrip Trip
		err = results.Scan(&newTrip.TripID, &newTrip.StatusOfTrip, &newTrip.PassengerID, &newTrip.DriverID, &newTrip.PickUp, &newTrip.DropOff)
		if err != nil {

			panic(err.Error())
		}
		trips = append(trips, newTrip) //Store them in a list -> var trips []Trip
	}
	return trips
}

func UpdateTrip(db *sql.DB, t Trip) {
	query := fmt.Sprintf("UPDATE Trips SET StatusOfTrip='%s', PassengerID='%s', DriverID='%s', PickUp='%s', DropOff='%s' WHERE TripID='%s'",
		t.StatusOfTrip, t.PassengerID, t.DriverID, t.PickUp, t.DropOff, t.TripID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// HTTP functions
func tripInit(resp http.ResponseWriter, req *http.Request) {
	var newTrip Trip

	if req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Please supply trip information in JSON format"))
			return
		}

		err = json.Unmarshal(reqBody, &newTrip)
		if err != nil {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Invalid trip information provided"))
			return
		}

		if newTrip.TripID == "" || newTrip.PassengerID == "" || newTrip.DriverID == "" || newTrip.PickUp == "" || newTrip.DropOff == "" {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Please supply all required trip information"))
			return
		}

		newTrip.StatusOfTrip = "Processing"
		InitNewTrip(db, newTrip)

		resp.WriteHeader(http.StatusCreated)
		resp.Write([]byte("201 - Successfully created trip"))
	} else {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported media type"))
	}
}

func tripUpdate(resp http.ResponseWriter, req *http.Request) {
	var trips []Trip
	queryParams := req.URL.Query()
	TripID := queryParams.Get("TripID")

	if req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Please supply trip information in JSON format"))
			return
		}

		err = json.Unmarshal(reqBody, &trips)
		if err != nil {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Invalid trip information provided"))
			return
		}

		if trips[0].TripID != TripID {
			resp.WriteHeader(http.StatusUnprocessableEntity)
			resp.Write([]byte("422 - Please supply a valid TripID and updated trip information"))
			return
		}

		UpdateTrip(db, trips[0])

		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("200 - Successfully updated trip"))
	} else {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		resp.Write([]byte("415 - Unsupported media type"))
	}
}

func tripRetrieve(resp http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	TripID := queryParams.Get("TripID")

	if TripID != "" {
		trips := GetTripInfo(db, TripID)
		if len(trips) == 0 {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte("404 - Trip not found"))
			return
		}

		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(trips)
	} else {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		resp.Write([]byte("422 - Please supply a valid TripID"))
	}
}

func tripDelete(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusForbidden)
	resp.Write([]byte("403 - For audit purposes, trip cannot be deleted."))
}

func trip(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		tripInit(resp, req)
	case "PUT":
		tripUpdate(resp, req)
	case "GET":
		tripRetrieve(resp, req)
	case "DELETE":
		tripDelete(resp, req)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		resp.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	// instantiate trips
	fmt.Println("Grab & Hail!")
	hailgrab_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hailgrab_db")

	db = hailgrab_db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/trips", trip).Methods(
		"GET", "PUT", "POST", "DELETE")
	fmt.Println("Trip --> Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

	defer db.Close()
}
