package main

import (
	"io/ioutil"

	"github.com/go-chi/chi"
	"github.com/hex_url_shortner/shortner"
	"github.com/pkg/errors"

	"log"
	"net/http"

	js "github.com/hex_url_shortner/serializer/json"
	ms "github.com/hex_url_shortner/serializer/msgpack"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortner.RedirectService
}

func NewHandler(redirectService shortner.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *handler) serializer(contentType string) shortner.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.Redirect{}
	}
	return &js.Redirect{}
}

func (h handler) Get(responseWriter http.ResponseWriter, req *http.Request) {

	code := chi.URLParam(req, "code")
	redirect, err := h.redirectService.Find(code)

	if err != nil {
		if errors.Cause(err) == shortner.ErrRedirectNotFound {
			http.Error(responseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(responseWriter, req, redirect.URL, http.StatusMovedPermanently)
}

func (h handler) Post(responseWriter http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	contentType := req.Header.Get("Content-Type")
	serializer := h.serializer(contentType)
	redirect, err := serializer.Decode(data)

	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortner.ErrRedirectInvalid {
			http.Error(responseWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := serializer.Encode(redirect)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(responseWriter, contentType, responseBody, http.StatusCreated)
}
