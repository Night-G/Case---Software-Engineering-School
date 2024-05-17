package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// exchange rate
const nbuAPI = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json"

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

// handle /subscribe?email= ' '
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	email := q.Get("email")
}

func main() {
	http.HandleFunc("/rate", rateHandler)

	http.HandleFunc("/subscribe", subscribeHandler)

	//run
	log.Fatal(http.ListenAndServe(":8080", nil))
}
