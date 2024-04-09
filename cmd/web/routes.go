package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/user_banner", app.getBanner)
	mux.HandleFunc("/banner", app.getBanners)
	mux.HandleFunc("/banner/{id}", app.updateBanner)

	return mux
}
