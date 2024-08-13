package main

import (
	"fmt"
	"net/http"
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

	// homeTpl.Execute(w, struct{Name, Price string}{Price: "413", Name:"Nike AirForce1"})
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
