package company

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/aditi-zs/Placement-API/entities"
	"github.com/aditi-zs/Placement-API/errors"
	"github.com/aditi-zs/Placement-API/service"
)

type handler struct {
	service service.CompanySvc
}

//nolint:revive // it's a factory function
func New(s service.CompanySvc) handler {
	return handler{service: s}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.service.Get(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte("error in marshaling"))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cmpID := mux.Vars(r)["id"]

	id, err := uuid.Parse(cmpID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: cmpID}.Error()))

		return
	}

	resp, err := h.service.GetByID(ctx, id)
	if err != nil {
		if _, ok := err.(errors.EntityNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte("error in marshaling"))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cmp entities.Company

	req, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = json.Unmarshal(req, &cmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid body"))

		return
	}

	err = validateBody(cmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	resp, err := h.service.Create(ctx, cmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte("error in marshaling"))
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respBody)
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cmpID := mux.Vars(r)["id"]

	id, err := uuid.Parse(cmpID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: cmpID}.Error()))

		return
	}

	req, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	var cmp entities.Company
	if err = json.Unmarshal(req, &cmp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid body"))

		return
	}

	err = validateBody(cmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	resp, err := h.service.Update(ctx, id, cmp)
	if err != nil {
		if _, ok := err.(errors.EntityNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte("error in marshaling"))
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respBody)
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cmpID := mux.Vars(r)["id"]

	id, err := uuid.Parse(cmpID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: cmpID}.Error()))

		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		if _, ok := err.(errors.EntityNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("Data deleted Successfully"))
}

func validateBody(c entities.Company) error {
	var missingParams []string
	if c.Name == "" {
		missingParams = append(missingParams, "name")
	}

	if c.Category == "" {
		missingParams = append(missingParams, "category")
	}

	if len(missingParams) != 0 {
		return errors.MissingParam{Param: missingParams}
	}

	return nil
}
