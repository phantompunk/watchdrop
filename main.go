package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/gocolly/colly"
)

type ProductDetail struct {
	Site  string
	Name  string
	Price string
}

func main() {
	http.HandleFunc("/", homeFunc)
	http.HandleFunc("/search", searchFunc)
	http.HandleFunc("/watchlist", watchFunc)

	fmt.Println("Started http server on port 8080")
	http.ListenAndServe(":8080", nil)
}

func homeFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Init home")
	homeTpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	homeTpl.Execute(w, nil)
}

const details = `<h1>{{.Name}}</h1>`

func searchFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Searching...")

	url := r.FormValue("product")
	fmt.Println("SearchQuery:" + url)

	// homeTpl := template.New("details")
	// homeTpl, err := homeTpl.Parse(details)
	homeTpl, err := template.ParseFiles("details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productInfo, err := scrapeInfo(url)
	if err != nil {
	}

	homeTpl.Execute(w, productInfo)
}

func scrapeInfo(url string) (ProductDetail, error) {
	fmt.Printf("Scraping info for product: %s\n", url)

	var price, name, site string
	c := colly.NewCollector()

	c.OnHTML("body", func(el *colly.HTMLElement) {
		prices := strings.Split(el.ChildText("#price-container > span"), "$")
		fmt.Println(prices)
		price = prices[1]
		name = el.ChildText("#title-container > h1")
		site = "nike.com"

		fmt.Println("Site: ", site)
		fmt.Println("Name: ", name)
		fmt.Println("Price: ", price)
	})
	c.Visit(url)
	return ProductDetail{Name: name, Price: price, Site: site}, nil
}

const (
	EMAILER_SYSTEM    = "GMAIL_SENDER_ADDRESS"
	EMAILER_APP_TOKEN = "GMAIL_SENDER_TOKEN"
)

func watchFunc(w http.ResponseWriter, r *http.Request) {
	watcher := r.PostFormValue("email")
	fmt.Printf("Emailing %s", watcher)

	from := os.Getenv(EMAILER_SYSTEM)
	password := os.Getenv(EMAILER_APP_TOKEN)
	if from == "" || password == "" {
		fmt.Println("Missing system email address or password")
	}
	to := []string{watcher}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("This is a simple email body.")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
	}
}
