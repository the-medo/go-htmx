package main

import (
	"go-htmx/util"
	"html/template"
	"log"
	"net/http"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/fragments/results.html"))
		data := map[string][]Stock{
			"Results": SearchTicker(r.URL.Query().Get("key"), config.PolygonApiKey),
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			return
		}
	})

	http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ticker := r.PostFormValue("ticker")
			stk := SearchTicker(ticker, config.PolygonApiKey)[0]
			val := GetDailyValues(ticker, config.PolygonApiKey)
			tmpl := template.Must(template.ParseFiles("./templates/index.html"))
			err := tmpl.ExecuteTemplate(w, "stock-element",
				Stock{Ticker: stk.Ticker, Name: stk.Name, Price: val.Open})
			if err != nil {
				return
			}
		}
	})

	log.Println("App running on 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
