package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	api "test-zero-agency/internal/api/dto"
	routerMiddleware "test-zero-agency/internal/api/middleware"
	"test-zero-agency/internal/entity"
	"time"
)

type peopleHandler struct {
	PeopleService entity.IPeopleService
	Logger        *zap.Logger
}

func RegisterPeopleHandlers(r *chi.Mux, service entity.IPeopleService, logger *zap.Logger, routerMiddleware routerMiddleware.IMiddleware) {
	PeopleHandler := new(peopleHandler)
	PeopleHandler.PeopleService = service
	PeopleHandler.Logger = logger

	r.Route("/people", func(r chi.Router) {
		r.Use(routerMiddleware.PanicRecovery)
		r.Use(middleware.Timeout(time.Second * 10))
		r.Use(middleware.RequestID)
		r.Use(routerMiddleware.ContentTypeJSON)

		r.Get("/", routerMiddleware.DebugLogger(PeopleHandler.GetPeopleList))
		r.Delete("/{id}", routerMiddleware.DebugLogger(PeopleHandler.DeletePeople))
		r.Put("/", routerMiddleware.DebugLogger(PeopleHandler.PutPeople))
		r.Post("/", routerMiddleware.DebugLogger(PeopleHandler.PostPeople))
	})
}

func (h peopleHandler) GetPeopleList(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	peopleListRequest, err := api.NewGetPeopleListRequest(r.URL.Query())

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.NewLogicError(err, "bad request", http.StatusBadRequest))
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	respPeople, err := h.PeopleService.GatPeopleList(r.Context(), peopleListRequest)

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	resp, _ := json.Marshal(respPeople)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return resp, http.StatusOK
}

func (h peopleHandler) DeletePeople(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.NewLogicError(err, "bad request", http.StatusBadRequest))
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	err = h.PeopleService.DeletePeople(r.Context(), id)

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	w.WriteHeader(http.StatusOK)
	return nil, http.StatusOK
}

func (h peopleHandler) PutPeople(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	var peopleRequest api.PutPeopleRequest
	err := json.NewDecoder(r.Body).Decode(&peopleRequest)

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.NewLogicError(err, "bad request", http.StatusBadRequest))
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	peopleResponse, err := h.PeopleService.UpdatePeople(r.Context(), &peopleRequest)

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	resp, _ := json.Marshal(peopleResponse)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return nil, http.StatusOK
}

func (h peopleHandler) PostPeople(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	peopleRequest, err := api.NewPostPeopleRequest(json.NewDecoder(r.Body))

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.NewLogicError(err, "bad request", http.StatusBadRequest))
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	peopleResponse, err := h.PeopleService.CreatePeople(r.Context(), peopleRequest)

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	resp, _ := json.Marshal(peopleResponse)

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	return nil, http.StatusCreated
}
