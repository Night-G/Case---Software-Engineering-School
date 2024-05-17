package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"net/http"
	"os"
)

// exchange rate
const nbuAPI = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json"

type Entity struct {
	email string
}

type ExchangeRate struct {
	Rate float32 `json:"rate"`
}

// get the rate from nbuAPI
func getExchangeRate() (float32, error) {
	response, err := http.Get(nbuAPI)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {

		var rate []ExchangeRate
		if err := json.NewDecoder(response.Body).Decode(&rate); err != nil {
			return 0, err
		}

		return rate[0].Rate, nil
	}
	return 0, fmt.Errorf("status code: %d", response.StatusCode)

}

// handle /rate
func rateHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := getExchangeRate()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(w, "1 USD = %.2f UAH", rate)
}

// handle  /subscribe-> /subscribe?email= ' '
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	email := q.Get("email")
	fmt.Fprintf(w, email)

	db, err := sql.Open("mssql", "server=DESKTOP-LMCD5AI ;user id=root;password=root")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`INSERT INTO emails(email) VALUES ($1)`, email)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// handle /subscribe for file  -> /subscribe/file?email= <email>
func subscribeInFileHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	email := q.Get("email")
	fmt.Fprintf(w, email)

	var f *os.File
	var err error
	if _, err := os.Stat("./test.txt"); err != nil {
		f, err = os.Create("test.txt")

	}

	_, err = f.WriteString(email)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/rate", rateHandler)

	http.HandleFunc("/subscribe", subscribeHandler)

	http.HandleFunc("/subscribe/file", subscribeInFileHandler)

	//run
	log.Fatal(http.ListenAndServe(":8080", nil))
}
