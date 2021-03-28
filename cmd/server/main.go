package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	pm "github.com/stanistan/present-me"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{owner}/{repo}/pull/{number}/{reviewID}", renderReview).
		Methods("GET").
		Name("review")

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("starting server at port %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func run(w http.ResponseWriter, r *http.Request, params *pm.ReviewParams) error {
	model, err := params.Model(
		pm.Context{
			Ctx:    r.Context(),
			Client: github.NewClient(nil),
		},
		r.URL.Query().Get("refresh") == "1",
	)
	if err != nil {
		return err
	}

	return model.AsMarkdown(w, pm.AsMarkdownOptions{
		AsHTML: true,
		InBody: true,
	})
}

func renderReview(w http.ResponseWriter, r *http.Request) {
	handle(w, func() error {
		params, err := pm.ReviewParamsFromMap(mux.Vars(r))
		if err != nil {
			return err
		}

		return run(w, r, params)
	})
}

func handle(w http.ResponseWriter, f func() error) {
	if err := f(); err != nil {
		log.Printf("Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
