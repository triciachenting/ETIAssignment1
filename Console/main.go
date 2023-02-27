package ETIAssignment1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const DRIVER_URL = "http://localhost:5000/drivers"
const PASSENGER_URL = "http://localhost:5001/passengers"
const TRIP_URL = "http://localhost:5002/trips"

type Driver struct {
	DriverID     string
	FirstName    string
	LastName     string
	MobileNumber string
	EmailAddress string
	CarPlateNo   string
	isAvailable  bool
}

type Passenger struct {
	PassengerID  string
	FirstName    string
	LastName     string
	MobileNumber string
	EmailAddress string
}

type Trip struct {
	TripID       string
	StatusOfTrip string
	PassengerID  string
	DriverID     string
	PickUp       string
	DropOff      string
}

// HTTP Functions
func postData(url string, body interface{}) (*http.Response, error) {
	reqBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	return resp, err
}

func updateData(url string, body interface{}) (*http.Response, error) {
	reqBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

// Start page for the main menu
func start() {
	mainMenu()
}

// Main menu when user launch HailGrab's Console Application
func mainMenu() {
	for {
		fmt.Println("Welcome to our HailGrab!")
		fmt.Println("[1] Login as a passenger user")
		fmt.Println("[2] Login as a driver user")
		fmt.Println("[3] Register as a passenger user")
		fmt.Println("[4] Register as driver user")
		fmt.Println("[0] Quit")

		fmt.Print("\nEnter your option: ")
		var option string
		fmt.Scanln(&option)

		if option == "1" {
			passengerLogin()
		} else if option == "2" {
			driverLogin()
		} else if option == "3" {
			passengerRegister()
		} else if option == "4" {
			driverRegister()
		} else if option == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			break
		}
	}
}

// Main menu features for user accessing HailGrab as Driver
func driverLogin() {
	for {
		fmt.Print("Enter your mobile number")
		var mobilenumber string
		fmt.Scanln(&mobilenumber)

		driver := getDriverMobileNumber(mobilenumber)
		if mobilenumber != driver.MobileNumber {
			fmt.Println("\nInvalid mobile number. Login failed.")
			break
		} else {
			fmt.Printf("\nWelcome to HailGrab Driver, %s %s!\n", driver.FirstName, driver.LastName)
			driverMainMenu(driver)
			break
		}
	}
}

func driverRegister() {
	//Prompt user to enter all the required information to register as a driver
	fmt.Print("Driver ID: ")
	var driverid string
	fmt.Scanln(&driverid)

	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobilenumber string
	fmt.Scanln(&mobilenumber)

	fmt.Print("Email Address: ")
	var emailaddress string
	fmt.Scanln(&emailaddress)

	fmt.Print("Car Plate Number: ")
	var carplateno string
	fmt.Scanln(&carplateno)

	newDriver := Driver{
		DriverID:     driverid,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNumber: mobilenumber,
		EmailAddress: emailaddress,
		CarPlateNo:   carplateno,
	}

	newDriver.isAvailable = true //Driver is available for trip booking by passenger
	err := InitNewDriver(newDriver)
	if err != nil {
		fmt.Println("Driver's account creation fail")
	} else {
		fmt.Println("Successfully created driver's account")
	}
}

//	Driver's menu features

func driverMainMenu(driver Driver) {
	for {
		fmt.Println("[1] Start trip")
		fmt.Println("[2] End trip")
		fmt.Println("[3] Update personal information")
		fmt.Println("[0] Logout")

		fmt.Print("\nEnter your option: ")
		var option string
		fmt.Scanln(&option)

		if option == "1" {

		} else if option == "2" {

		} else if option == "3" {
			driverUpdate(driver)
		} else if option == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			driverMainMenu(driver)
		}
	}
}

func driverUpdate(driver Driver) {
	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobilenumber string
	fmt.Scanln(&mobilenumber)

	fmt.Print("Email Address: ")
	var emailaddress string
	fmt.Scanln(&emailaddress)

	fmt.Print("Car Plate Number: ")
	var carplateno string
	fmt.Scanln(&carplateno)

	newDriverInfo := Driver{
		DriverID:     driver.DriverID,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNumber: mobilenumber,
		EmailAddress: emailaddress,
		CarPlateNo:   carplateno,
	}
	err := UpdateDriver(newDriverInfo)
	if err != nil {
		fmt.Println("Driver information update failed")
	} else {
		fmt.Println("Successfully updated driver's information")
	}
}

//	Main menu features for user accessing HailGrab as Passenger

func passengerLogin() {
	fmt.Print("Please enter your mobile number: ") //Passenger login using their mobile number
	var mobilenumber string
	fmt.Scanln(&mobilenumber)

	//Look for registered driver with the given mobile number
	passenger := getPassengerMobileNumber(mobilenumber)
	if mobilenumber != passenger.MobileNumber {
		fmt.Println("\nInvalid mobile number. Login failed.") //If mobile number does not exist, login fail.
	} else {
		fmt.Printf("\nWelcome to HailGrab Passenger, %s %s!\n", passenger.FirstName, passenger.LastName)
		passengerMainMenu(passenger)
	}
}

func passengerRegister() {
	//Prompt user to enter all the required information to register as a passenger
	fmt.Print("Passenger ID: ")
	var passengerid string
	fmt.Scanln(&passengerid)

	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobilenumber string
	fmt.Scanln(&mobilenumber)

	fmt.Print("Email Address: ")
	var emailaddress string
	fmt.Scanln(&emailaddress)

	newPassenger := Passenger{
		PassengerID:  passengerid,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNumber: mobilenumber,
		EmailAddress: emailaddress,
	}

	err := InitNewPassenger(newPassenger)
	if err != nil {
		fmt.Println("Passenger's account creation fail")
	} else {
		fmt.Println("Successfully created passenger's account")
	}
}

// Passenger's menu features
func passengerMainMenu(passenger Passenger) {
	for {
		fmt.Println("[1] New trip")
		fmt.Println("[2] View past trip")
		fmt.Println("[3] Update my information")
		fmt.Println("[0] Logout")

		fmt.Print("\nEnter your option: ")
		var option string
		fmt.Scanln(&option)

		if option == "1" {

		} else if option == "2" {
			pastPassengerTrips(passenger)
		} else if option == "3" {
			passengerUpdate(passenger)
		} else if option == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			passengerMainMenu(passenger)
		}
	}
}

// View passenger trip
func pastPassengerTrips(passenger Passenger) {
	trips := viewPassengerTrips(passenger.PassengerID)

	for t := len(trips) - 1; t >= 0; t-- { //Look for
		trip := trips[t]
		fmt.Println()
		fmt.Println("Trip ID: ", trip.TripID)
		fmt.Println("Trip Status: ", trip.StatusOfTrip)
		fmt.Println("Driver ID: ", trip.DriverID)
		fmt.Println("Passenger ID: ", trip.PassengerID)
		fmt.Println("Pick Up Location: ", trip.PickUp)
		fmt.Println("Drop Off Location: ", trip.DropOff)
		fmt.Println()
	}
}

// Update passenger's information
func passengerUpdate(passenger Passenger) {
	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobilenumber string
	fmt.Scanln(&mobilenumber)

	fmt.Print("Email Address: ")
	var emailaddress string
	fmt.Scanln(&emailaddress)

	newPassengerInfo := Passenger{
		PassengerID:  passenger.PassengerID,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNumber: mobilenumber,
		EmailAddress: emailaddress,
	}
	err := UpdatePassenger(newPassengerInfo)
	if err != nil {
		fmt.Println("Passenger information update failed")
	} else {
		fmt.Println("Successfully updated passenger's information")
	}
}

//Functions for Passenger API

func InitNewPassenger(newPassenger Passenger) error {
	url := PASSENGER_URL
	_, err := postData(url, newPassenger)
	return err
}

func UpdatePassenger(newPassengerInfo Passenger) error {
	url := fmt.Sprintf("%s?PassengerID=%s", PASSENGER_URL, newPassengerInfo.PassengerID)

	_, err := updateData(url, newPassengerInfo)
	return err
}

func getPassengerMobileNumber(MobileNumber string) Passenger {
	var passenger Passenger

	url := fmt.Sprintf("%s?MobileNumber=%s", PASSENGER_URL, MobileNumber)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return passenger
	}

	json.NewDecoder(resp.Body).Decode(&passenger)
	return passenger
}

func viewPassengerTrips(PassengerID string) []Trip {
	var trips []Trip

	url := fmt.Sprintf("%s?PassengerID=%s", TRIP_URL, PassengerID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Passenger is not assigned to any trips")
		return trips
	}

	json.NewDecoder(resp.Body).Decode(&trips)
	return trips
}

// Functions for Driver API

func InitNewDriver(newDriver Driver) error {
	url := DRIVER_URL

	_, err := postData(url, newDriver)
	return err
}

func UpdateDriver(newDriverInfo Driver) error {
	url := fmt.Sprintf("%s?DriverID=%s", DRIVER_URL, newDriverInfo.DriverID)

	_, err := updateData(url, newDriverInfo)
	return err
}

func getDriverMobileNumber(MobileNumber string) Driver {
	var driver Driver

	url := fmt.Sprintf("%s?MobileNumber=%s", DRIVER_URL, MobileNumber)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return driver
	}

	json.NewDecoder(resp.Body).Decode(&driver)
	return driver
}

func getAvailableDriver() Driver {
	var driver Driver

	url := fmt.Sprintf("%s/available?isAvailable=%t", DRIVER_URL, true)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return driver
	}

	json.NewDecoder(resp.Body).Decode(&driver)
	return driver
}
