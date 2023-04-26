package student

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/Zopsmart-Training/student-recruitment-system/entities"
	"github.com/Zopsmart-Training/student-recruitment-system/errors"
	"github.com/Zopsmart-Training/student-recruitment-system/service"
)

type handler struct {
	service service.StudentSvc
}

//nolint:revive // it's a factory function
func New(s service.StudentSvc) handler {
	return handler{service: s}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	branch := r.URL.Query().Get("branch")
	includeCompany := r.URL.Query().Get("includeCompany")

	resp, err := h.service.Get(ctx, name, branch, includeCompany)
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

	studentID := mux.Vars(r)["id"]

	id, err := uuid.Parse(studentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: studentID}.Error()))

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
	var stu entities.Student

	ctx := r.Context()

	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = json.Unmarshal(req, &stu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid body"))

		return
	}

	err = validateBody(&stu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	resp, err := h.service.Create(ctx, &stu)
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

	studentID := mux.Vars(r)["id"]

	id, err := uuid.Parse(studentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: studentID}.Error()))

		return
	}

	var stu entities.Student

	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = json.Unmarshal(req, &stu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid body"))

		return
	}

	err = validateBody(&stu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	resp, err := h.service.Update(ctx, id, &stu)
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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respBody)
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	studentID := mux.Vars(r)["id"]

	id, err := uuid.Parse(studentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.InvalidParam{Param: studentID}.Error()))

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

func validateBody(s *entities.Student) error {
	var missingParams []string
	if s.Name == "" {
		missingParams = append(missingParams, "name")
	}

	if s.Phone == "" {
		missingParams = append(missingParams, "phone")
	}

	if s.DOB == "" {
		missingParams = append(missingParams, "dob")
	}

	if s.Branch == "" {
		missingParams = append(missingParams, "branch")
	}

	if s.Comp.ID == uuid.Nil {
		missingParams = append(missingParams, "company id")
	}

	if s.Status == "" {
		missingParams = append(missingParams, "status")
	}

	if len(missingParams) != 0 {
		return errors.MissingParam{Param: missingParams}
	}

	return nil
}
