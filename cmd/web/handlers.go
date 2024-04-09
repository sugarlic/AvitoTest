package main

import "net/http"

func (app *application) getBanner(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Banner"))
}

func (app *application) getBanners(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Banner list"))
}

func (app *application) updateBanner(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update"))
}
