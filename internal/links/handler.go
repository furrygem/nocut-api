package links

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/furrygem/nocut-api/internal/handlers"
	"github.com/furrygem/nocut-api/pkg/logging"
	"github.com/gorilla/mux"
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

func (h *handler) Register(router *mux.Router, urlPrefix string) {
	Links := "/links"
	LinkByID := "/links/{id}"
	CheckLink := "/url/{b64url}/check"

	if urlPrefix != "" {
		Links = AppendPrefixToURL(urlPrefix, "/links")
		LinkByID = AppendPrefixToURL(urlPrefix, "/links/{id}")
		CheckLink = AppendPrefixToURL(urlPrefix, "/url/{b64url}/check")
	}

	// router.HandleFunc("/links", h.GetList).Methods("GET")
	router.HandleFunc(Links, h.CreateHandler).Methods("POST")
	router.HandleFunc(LinkByID, h.GetLinkByIDHandler).Methods("GET")
	router.HandleFunc("/{slug}", h.GetLinkSlugHandler).Methods("GET")
	router.HandleFunc(CheckLink, h.CheckLinkHandler).Methods("GET")
}

func (h *handler) CheckLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := base64.RawURLEncoding.DecodeString(vars["b64url"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	link := string(l)
	uce := h.service.CheckLinkService(link)
	resp, err := uce.JSON()
	if err != nil {
		h.logger.Errorf("Error marshaling URLCheckException '%v' JSON. %v", resp, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}

func (h *handler) GetLinkSlugHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	id, err := URLToID(slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request. %s", err.Error()), http.StatusBadRequest)
		return
	}
	l, err := h.service.GetLinkByIDService(context.Background(), id)
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
	l, err := h.service.CreateService(context.Background(), dto)
	if err != nil {
		e, ok := err.(*URLCheckException)
		if ok {
			// h.service.SendURLCheckResults(e, w)
			resp, err := e.JSON()
			if err != nil {
				h.logger.Errorf("Error marshaling URLCheckException '%v' JSON. %v", e, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := json.Marshal(l)
	location := fmt.Sprintf("/links/%s", l.ID)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	return
}

func (h *handler) GetLinkByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	link, err := h.service.GetLinkByIDService(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp, _ := json.Marshal(link)
	w.Write(resp)
	return
}
