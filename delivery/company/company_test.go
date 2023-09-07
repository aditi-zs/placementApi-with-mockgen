package company

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/aditi-zs/Placement-API/entities"
	"github.com/aditi-zs/Placement-API/errors"
	"github.com/aditi-zs/Placement-API/service"
)

func initializeTest(t *testing.T) *service.MockCompanySvc {
	ctrl := gomock.NewController(t)
	mockCompany := service.NewMockCompanySvc(ctrl)

	return mockCompany
}

func TestGet(t *testing.T) {
	mockCompany := initializeTest(t)
	validID := uuid.New()
	tests := []struct {
		description string
		res         []entities.Company
		err         error
		statusCode  int
	}{
		{"Success case: All entries are present", []entities.Company{{ID: validID, Name: "Wipro", Category: "MASS"}},
			nil, 200},
		{"Error case: server error", nil, errors.DB{Reason: "server error"}, 400},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("GET", "/companies", http.NoBody)
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()
		h := New(mockCompany)
		mockCompany.EXPECT().Get(gomock.Any()).Return(tc.res, tc.err)
		h.Get(resRec, req)

		var val []entities.Company
		_ = json.Unmarshal(resRec.Body.Bytes(), &val)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%d)", i, tc.description)
		assert.Equal(t, tc.res, val, "Test[%v] failed\n(%v)", i, tc.description)
	}
}
func TestGetByID(t *testing.T) {
	mockCompany := initializeTest(t)
	validID := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		mockTimes   int
		mockRes     entities.Company
		mockErr     error
		expRes      entities.Company
		statusCode  int
	}{
		{"Success case: for valid id", validID, 1, entities.Company{ID: validID, Name: "Wipro", Category: "MASS"},
			nil, entities.Company{ID: validID, Name: "Wipro", Category: "MASS"}, 200},
		{"Error case: db error", validID, 1, entities.Company{}, errors.DB{Reason: "server error"},
			entities.Company{}, 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("GET", "/companies/{id}", http.NoBody)
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})
		h := New(mockCompany)
		mockCompany.EXPECT().GetByID(gomock.Any(), tc.inputID).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)

		h.GetByID(resRec, req)

		var val entities.Company
		_ = json.Unmarshal(resRec.Body.Bytes(), &val) // json to go

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%d] failed\n(%v)", i, tc.description)
	}
}

func TestCreate(t *testing.T) {
	mockCompany := initializeTest(t)
	validID, err := uuid.Parse("1fa46d13-6a50-11ed-90d1-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		description string
		input       string
		mockTimes   int
		mockInput   entities.Company
		mockRes     entities.Company
		mockErr     error
		expRes      string
		statusCode  int
	}{
		{"Success case: All entries are present", `{"name":"Wipro","category":"MASS"}`, 1, entities.Company{Name: "Wipro", Category: "MASS"},
			entities.Company{ID: validID, Name: "Wipro", Category: "MASS"}, nil,
			`{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4","name":"Wipro","category":"MASS"}`, 201,
		},
		{"Error case: unmarshal err", ``, 0, entities.Company{}, entities.Company{},
			nil, "invalid body", 400,
		},
		{"Failure case: db error", `{"name":"Wipro","category":"MASS"}`, 1, entities.Company{Name: "Wipro", Category: "MASS"},
			entities.Company{}, errors.DB{Reason: "server error"}, "DB Error: server error", 400,
		},
		{"Failure case: missing parameters", `{}`, 0, entities.Company{},
			entities.Company{}, nil, "Missing Parameter: name,category", 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("POST", "/companies", strings.NewReader(tc.input))
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()

		h := New(mockCompany)
		mockCompany.EXPECT().Create(gomock.Any(), tc.mockInput).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		h.Create(resRec, req)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.expRes, resRec.Body.String(), "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	mockCompany := initializeTest(t)
	validID, err := uuid.Parse("1fa46d13-6a50-11ed-90d1-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		description string
		inputID     uuid.UUID
		input       string
		mockTimes   int
		mockInput   entities.Company
		mockRes     entities.Company
		mockErr     error
		expRes      string
		statusCode  int
	}{
		{"Success case: All entries are present", validID, `{"name":"Google","category":"DREAM IT"}`, 1,
			entities.Company{Name: "Google", Category: "DREAM IT"}, entities.Company{ID: validID, Name: "Google", Category: "DREAM IT"},
			nil, `{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4","name":"Google","category":"DREAM IT"}`, 201,
		},
		{"Error case: unmarshal err", validID, `{`, 0, entities.Company{}, entities.Company{},
			nil, "invalid body", 400,
		},
		{"Error case: db error", validID, `{"name":"Wipro","category":"MASS"}`, 1,
			entities.Company{Name: "Wipro", Category: "MASS"}, entities.Company{},
			errors.DB{Reason: "server error"}, "DB Error: server error", 400,
		},
		{
			"Error case: missing parameters", validID, `{}`, 0, entities.Company{}, entities.Company{},
			nil, "Missing Parameter: name,category", 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("PUT", "/companies/{id}", strings.NewReader(tc.input))
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})
		h := New(mockCompany)
		mockCompany.EXPECT().Update(gomock.Any(), tc.inputID, tc.mockInput).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		h.Update(resRec, req)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.expRes, resRec.Body.String(), "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	id := uuid.New()
	mockCompany := initializeTest(t)
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         error
		statusCode  int
	}{
		{"Success case: for valid id", id, nil, 204},
		{"Error case: when id is valid but id is not present in db", id, errors.DB{Reason: "server error"}, 400},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("DELETE", "/companies/{id}", http.NoBody)
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})
		h := New(mockCompany)
		mockCompany.EXPECT().Delete(gomock.Any(), tc.inputID).Return(tc.res)
		h.Delete(resRec, req)

		var actRes entities.Company

		_ = json.Unmarshal(resRec.Body.Bytes(), &actRes)
		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
	}
}
