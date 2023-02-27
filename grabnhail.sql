/* STEP 1. Create username and password */
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

-- STEP 2. Create Database
CREATE DATABASE hailgrab;

-- STEP 3. Use the Database to Create Table
USE hailgrab;

CREATE TABLE Passengers (PassengerID VARCHAR(5) NOT NULL PRIMARY KEY, FirstName VARCHAR(50), LastName VARCHAR(50), MobileNumber VARCHAR(10), EmailAddress VARCHAR(100)); 
CREATE TABLE Drivers (DriverID VARCHAR(5) NOT NULL PRIMARY KEY, FirstName VARCHAR(50), LastName VARCHAR(50), MobileNumber VARCHAR(10), EmailAddress VARCHAR(100), CarPlateNo VARCHAR(100), isAvailable BOOL);  
CREATE TABLE Trips (TripID VARCHAR(5) NOT NULL PRIMARY KEY, StatusOfTrip VARCHAR(20), PassengerID VARCHAR(5), DriverID VARCHAR(5), PickUp VARCHAR(100), DropOff VARCHAR(100)); 

-- STEP 4. With the tables that are created, insert some data to Passengers, Drivers and Trips 
--    to see if you can test the GET, PUT Method that you have defined in microservices.
--    You are required to set up Database on SMSS, then test POST, PUT and GET via POSTMAN, then look at your database to see if any of the value are up there. 
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES ("0001", "Sulli", "Lim", "81234567", "sullilim@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES ("0002", "Lim", "Jiu", "87654321", "limjiu@gmail.com");

INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, CarPlateNo) VALUES ("0001", "Mee Su", "To", "91234567", "meesuto@gmail.com","S1234567A");
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, CarPlateNo) VALUES ("0002", "Lai", "Ong Tek", "97654321", "laiongtek@gmail.com","S7938281B");

INSERT INTO Trips (TripID, StatusOfTrip, PassengerID, DriverID, PickUp, DropOff) VALUES ("0001", "Ongoing", "0001", "0002", "Bukit Timah","Bukit Merah");
INSERT INTO Trips (TripID, StatusOfTrip, PassengerID, DriverID, PickUp, DropOff) VALUES ("0002", "Ended", "0001", "0001", "Holland Village","Changi Airport");