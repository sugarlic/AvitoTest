package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/avitoTest/pkg/models"
)

func (app *application) getBanner(w http.ResponseWriter, r *http.Request) {
	// tag_id, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	// if err != nil {
	// 	app.notFound(w)
	// 	return
	// }
	// feature_id, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	// if err != nil {
	// 	app.notFound(w)
	// 	return
	// }
	// use_last_revision, err := strconv.ParseBool(r.URL.Query().Get("tag_id"))
	// if err != nil {
	// 	use_last_revision = false
	// }
	// token := r.Header.Get("Authorization")
}

// получение списка баннеров или создание баннера
func (app *application) Banners(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.clientError(w, http.StatusUnauthorized)
	}

	if r.Method == http.MethodPost {
		if token != "admin_token" {
			app.noAccess(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}

		var data models.Banner
		err = json.Unmarshal(body, &data)
		if err != nil {
			app.serverError(w, err)
			return
		}

		id, err := app.banners.Insert(data)
		if err != nil {
			app.clientError(w, 400)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(id)))
	}
	// if r.Method == http.MethodGet {

	// }
}

func (app *application) updateBanner(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.clientError(w, http.StatusUnauthorized)
	}
	if token != "admin_token" {
		app.noAccess(w)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		err = app.banners.Delete(id)
		if err != nil {
			app.notFound(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}

		var data models.Banner
		err = json.Unmarshal(body, &data)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.banners.Update(data)
		if err != nil {
			app.clientError(w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
