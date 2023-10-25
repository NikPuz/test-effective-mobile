package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	api "test-zero-agency/internal/api/dto"
	routerMiddleware "test-zero-agency/internal/api/middleware"
	"test-zero-agency/internal/entity"
	"time"
)

type peopleHandler struct {
	PeopleService entity.IPeopleService
}

func RegisterPeopleHandlers(r *chi.Mux, service entity.IPeopleService, routerMiddleware routerMiddleware.IMiddleware) {
	PeopleHandler := new(peopleHandler)
	PeopleHandler.PeopleService = service

	r.Route("/people", func(r chi.Router) {
		r.Use(routerMiddleware.PanicRecovery)
		r.Use(middleware.Timeout(time.Second * 10))
		r.Use(middleware.RequestID)
		r.Use(routerMiddleware.ContentTypeJSON)

		r.Get("/", routerMiddleware.RequestLogger(PeopleHandler.GetPeopleList))
		//r.Delete("/{id}", routerMiddleware.RequestLogger(PeopleHandler.DeletePeople))
		//r.Put("/", routerMiddleware.RequestLogger(PeopleHandler.UpdatePeople))
		//r.Post("/", routerMiddleware.RequestLogger(PeopleHandler.CreatePeople))
	})
}

func (h peopleHandler) GetPeopleList(w http.ResponseWriter, r *http.Request) ([]byte, int, *entity.LogicError) {

	getPeopleListRequest, err := api.NewGetPeopleListRequest(r.URL.Query())

	if err != nil {
		logicError := entity.NewLogicError(err, http.StatusBadRequest)
		resp := logicError.JsonMarshal()
		w.WriteHeader(logicError.Code)
		w.Write(resp)
		return resp, logicError.Code, logicError
	}

	respPeople, err := h.PeopleService.ReadPeopleList(r.Context(), getPeopleListRequest)

	if err != nil {
		logicError := entity.ResponseLogicError(err)
		resp := logicError.JsonMarshal()
		w.WriteHeader(logicError.Code)
		w.Write(resp)
		return resp, logicError.Code, logicError
	}

	resp, _ := json.Marshal(respPeople)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return resp, http.StatusOK, nil
}
