package links

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/furrygem/nocut-api/internal/handlers"
	"github.com/furrygem/nocut-api/pkg/logging"
	"github.com/gorilla/mux"
)

const (
	linksURL = "/links"
	linkURL  = "/:slug"
)

type handler struct {
	logger  logging.Logger
	service Service
}

// NewHandler Return pointer to links handler
func NewHandler(logger logging.Logger, storage Storage, linkTTL time.Duration) handlers.Handler {
	service := Service{
		logger:  &logger,
		storage: storage,
		linkTTL: linkTTL,
	}

	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/links", h.GetList).Methods("GET")
	router.HandleFunc("/links", h.CreateHandler).Methods("POST")
	router.HandleFunc("/links/{id}", h.GetLinkByIdHandler).Methods("GET")
	router.HandleFunc("/{slug}", h.GetLinkSlugHandler).Methods("GET")
}

func (h *handler) GetLinkSlugHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	id, err := UrlToId(slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request. %s", err.Error()), http.StatusBadRequest)
		return
	}
	l, err := h.service.GetLinkById(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get link by id. %s", err.Error()), http.StatusInternalServerError)
	}
	http.Redirect(w, r, l.Source, http.StatusTemporaryRedirect)
}

func (h *handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	dto := CreateLinkDTO{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	l, err := h.service.Create(context.Background(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := json.Marshal(l)
	location := fmt.Sprintf("/links/%s", l.ID)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *handler) GetLinkByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	link, err := h.service.GetLinkById(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp, _ := json.Marshal(link)
	w.Write(resp)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the list of active link"))
}
