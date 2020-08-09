// main.go
package main

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

// Site - Our struct for all sites
type Site struct {
	ID        int     `json:"id"`
	NAME      string  `json:"name"`
	LATITUDE  float32 `json:"latitude"`
	LONGITUDE float32 `json:"longitude"`
	STATUS    float32 `json:"status"`
}

type Response struct {
	MESSAGE string `json:"message"`
}

var db *sql.DB
var err error

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go API !")
	log.Println("Endpoint Hit: homePage")
}

func getSites(w http.ResponseWriter, r *http.Request) {
	var sites []Site
	log.Println("Getting All Sites from DB")
	results, err := db.Query("SELECT ID,NAME,LATITUDE,LONGITUDE,STATUS FROM SITE")
	// if there is an error getting, handle it
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var site Site
		err = results.Scan(&site.ID, &site.NAME, &site.LATITUDE, &site.LONGITUDE, &site.STATUS)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		sites = append(sites, site)
	}
	//x, err := xml.MarshalIndent(sites, "", "  ")
	//w.Write(x)
	json.NewEncoder(w).Encode(sites)
}

func getSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["id"]
	results, err := db.Query("SELECT ID,NAME,LATITUDE,LONGITUDE,STATUS FROM SITE WHERE ID = ?", siteId)
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		var site Site
		err = results.Scan(&site.ID, &site.NAME, &site.LATITUDE, &site.LONGITUDE, &site.STATUS)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		json.NewEncoder(w).Encode(site)
	}
}

func createSite(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	var response Response
	reqBody, _ := ioutil.ReadAll(r.Body)
	var site Site
	json.Unmarshal(reqBody, &site)
	insertQuery, err := db.Prepare("INSERT INTO SITE (NAME,LATITUDE,LONGITUDE,STATUS) VALUES (?,?,?,?)")
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
		response.MESSAGE = "An Exception occoured while framing the site"
		json.NewEncoder(w).Encode(response)
	}
	result, err := insertQuery.Exec(site.NAME, site.LATITUDE, site.LONGITUDE, 0)
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
		response.MESSAGE = err.Error()
		json.NewEncoder(w).Encode(response)
	}
	log.Println(result.RowsAffected)
	response.MESSAGE = "1 Site Inserted Successfully"
	json.NewEncoder(w).Encode(response)
}

func deleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["id"]
	result, err := db.Exec("DELETE FROM SITE WHERE ID = ?", siteId)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(result.RowsAffected)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/sites", getSites).Methods("GET")
	myRouter.HandleFunc("/sites", createSite).Methods("POST")
	myRouter.HandleFunc("/sites/{id}", deleteSite).Methods("DELETE")
	myRouter.HandleFunc("/sites/{id}", getSite).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func insertSites(name string) {
	insertQuery, err := db.Prepare("INSERT INTO SITE (NAME,LATITUDE,LONGITUDE,STATUS) VALUES (?,?,?,?)")
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
	}
	result, err := insertQuery.Exec(name, 12.56, 12.56, 0)
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
	}
	log.Println(result.RowsAffected)
}

func main() {
	db, err = sql.Open("mysql", "hari:hari@tcp(127.0.0.1:3306)/Atlas")
	if err != nil {
		log.Fatal("Error occoured when getting a connection to database")
		panic(err.Error())
	}

	/*for i := 1000; i < 10000; i++ {
		var name = "TestSite-" + strconv.Itoa(i)
		insertSites(name)
	}*/

	defer db.Close()
	handleRequests()
}
