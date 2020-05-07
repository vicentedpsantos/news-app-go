package handlers

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vicentedpsantos/news-app-go/models"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

var apiKey *string

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	apiKey = flag.String("apikey", "", "Newsapi.org access key")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey must be set")
	}

	u, err := url.Parse(r.URL.String())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	search := &models.Search{}
	search.SearchKey = searchKey

	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}

	search.NextPage = next
	pageSize := 20

	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(search.SearchKey), pageSize, search.NextPage, *apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		newError := &models.NewsAPIError{}
		err := json.NewDecoder(resp.Body).Decode(newError)
		if err != nil {
			http.Error(w, "Unexpected server error", http.StatusInternalServerError)
			return
		}

		http.Error(w, newError.Message, http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&search.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / pageSize)))

	if ok := !search.IsLastPage(); ok {
		search.NextPage++
	}

	err = tpl.Execute(w, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
