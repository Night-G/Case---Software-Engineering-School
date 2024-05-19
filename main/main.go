package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
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

	db, err := sql.Open("sqlite3", "dbFile.db")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	// check if in the db already such email if none - write`
	err = db.QueryRow(`SELECT email FROM emails WHERE email = $1`, email).Scan()
	if err == sql.ErrNoRows {
		_, err = db.Exec(`INSERT INTO emails(email) VALUES ($1)`, email)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

// check file txt for the same row
func checkEmail(fileName string, email string) (isPresent bool, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return true, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	found := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), email) {
			found = true
			break
		}
	}

	return found, err
}

// handle /subscribe for file  -> /subscribe/file?email= <email>
func subscribeInFileHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	email := q.Get("email")
	fmt.Fprintf(w, email)

	fileName := "emails.txt"

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	// if there is no such email already - writes
	isPresent, err := checkEmail(fileName, email)
	if err != nil {
		log.Fatal(err)
		return
	}
	if !isPresent {
		email += "\n"
		_, err = f.WriteString(email)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func sendMail(sendTo []string) {
	apw := ""          // ---- put created app password from gmail from which sending emails
	sendingEmail := "" // ---- put gmail from which sending emails
	auth := smtp.PlainAuth("", sendingEmail, apw, "smtp.gmail.com")
	to := "To: " + sendTo[0] + "\r\n"

	rateFloat, err := getExchangeRate()
	if err != nil {
		log.Fatal(err)
		return
	}
	// float32 to string
	rate := strconv.FormatFloat(float64(rateFloat), 'f', -1, 32)
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := []byte(to + "Subject: Rate\r\n" + rate)

	//sending
	err = smtp.SendMail("smtp.gmail.com:587", auth, sendingEmail, sendTo, msg)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// sends rate on subscribed emails
func sendEmailsHandler() {

	//getting emails from db
	db, err := sql.Open("sqlite3", "dbFile.db")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	// get all subscribed emails
	res, err := db.Query(`SELECT email FROM emails`)
	if err != nil {
		log.Fatal(err)
		return
	}

	for res.Next() {
		var sendTo string
		err = res.Scan(&sendTo)
		if err != nil {
			log.Fatal(err)
		}

		sendMail([]string{sendTo})
	}
}

func main() {
	http.HandleFunc("/rate", rateHandler)

	http.HandleFunc("/subscribe", subscribeHandler)

	http.HandleFunc("/subscribe/file", subscribeInFileHandler)

	// http.HandleFunc("/sendEmails", sendEmailsHandler) - used for testing purposes

	// daily sending usd:UAH rate
	c := cron.New()
	_, err := c.AddFunc("@daily", sendEmailsHandler)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	//run
	log.Fatal(http.ListenAndServe(":8080", nil))
}
