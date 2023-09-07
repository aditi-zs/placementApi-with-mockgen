package student

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

const valID = "1fa46d13-6a50-11ed-90d1-64bc589051b4"

func initializeTest(t *testing.T) *service.MockStudentSvc {
	ctrl := gomock.NewController(t)
	mockStudent := service.NewMockStudentSvc(ctrl)

	return mockStudent
}

func TestGet(t *testing.T) {
	id := uuid.New()
	cmpID, err := uuid.Parse(valID)

	if err != nil {
		t.Errorf(err.Error())
	}

	mockStudent := initializeTest(t)

	tests := []struct {
		description    string
		name           string
		branch         string
		includeCompany string
		res            []entities.Student
		err            error
		statusCode     int
	}{
		{"Success case: All entries are present", "Monika", "ECE", "true",
			[]entities.Student{{ID: id, Name: "Monika", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil, 200,
		},
		{"Error case: server error", "Aditi", "CSE", "true",
			nil, errors.DB{Reason: "server error"}, 400,
		},
	}

	for i, tc := range tests {
		req := httptest.NewRequest("GET", "/students", http.NoBody)
		r := req.URL.Query()
		r.Set("name", tc.name)
		r.Set("branch", tc.branch)
		r.Set("includeCompany", tc.includeCompany)
		req.URL.RawQuery = r.Encode()
		resRec := httptest.NewRecorder()
		h := New(mockStudent)
		mockStudent.EXPECT().Get(gomock.Any(), tc.name, tc.branch, tc.includeCompany).Return(tc.res, tc.err)

		h.Get(resRec, req)

		var val []entities.Student
		_ = json.Unmarshal(resRec.Body.Bytes(), &val)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.res, val, "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestGetByID(t *testing.T) {
	id := uuid.New()
	cmpID, err := uuid.Parse(valID)

	if err != nil {
		t.Errorf(err.Error())
	}

	mockStudent := initializeTest(t)

	tests := []struct {
		description string
		inputID     uuid.UUID
		mockTimes   int
		mockRes     entities.Student
		mockErr     error
		expRes      entities.Student
		statusCode  int
	}{
		{"Success case: for valid id", id, 1, entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil,
			entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 200,
		},
		{"Error case: when id is valid but id is not present in db", id, 1, entities.Student{},
			errors.DB{Reason: "server error"}, entities.Student{}, 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("GET", "/students/{id}", http.NoBody)
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})
		h := New(mockStudent)
		mockStudent.EXPECT().GetByID(gomock.Any(), tc.inputID).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		h.GetByID(resRec, req)

		var val entities.Student
		_ = json.Unmarshal(resRec.Body.Bytes(), &val)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.expRes, val, "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestCreate(t *testing.T) {
	mockStudent := initializeTest(t)
	id, err := uuid.Parse("71bbdbb9-6bde-11ed-aaff-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	cmpID, err := uuid.Parse("1fa46d13-6a50-11ed-90d1-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		description string
		input       string
		mockTimes   int
		mockInput   entities.Student
		mockRes     entities.Student
		mockErr     error
		expRes      string
		statusCode  int
	}{
		{"Success case: All entries are present",
			`{"name":"Monika Jaiswal","phone":"6388768118","dob":"02/07/2000","branch":"ECE",` +
				`"comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`,
			1,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID},
				Status: "ACCEPTED"},
			entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, nil,
			`{"id":"71bbdbb9-6bde-11ed-aaff-64bc589051b4","name":"Monika Jaiswal","phone":"6388768118","dob":"02/07/2000","branch":"ECE",` +
				`"comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`, 201,
		},
		{"Error case: unmarshal error", ``, 0, entities.Student{},
			entities.Student{}, nil, "invalid body", 400,
		},
		{"Error case: Failure case: db error",
			`{"name":"Monika Jaiswal","phone":"6388768118","dob":"02/07/2000","branch":"ECE",
"comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`, 1,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, entities.Student{},
			errors.DB{Reason: "server error"}, "DB Error: server error", 400,
		},
		{"Error case: missing parameters", `{}`, 0, entities.Student{},
			entities.Student{}, nil, "Missing Parameter: name,phone,dob,branch,company id,status", 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("POST", "/students", strings.NewReader(tc.input))
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()
		h := New(mockStudent)
		mockStudent.EXPECT().Create(gomock.Any(), &tc.mockInput).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		h.Create(resRec, req)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.expRes, resRec.Body.String(), "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	mockStudent := initializeTest(t)
	id, err := uuid.Parse("71bbdbb9-6bde-11ed-aaff-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	cmpID, err := uuid.Parse("1fa46d13-6a50-11ed-90d1-64bc589051b4")

	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		description string
		inputID     uuid.UUID
		input       string
		mockTimes   int
		mockInput   entities.Student
		mockRes     entities.Student
		mockErr     error
		expRes      string
		statusCode  int
	}{
		{"Success case: for valid id", id,
			`{"name":"Monika Jaiswal","phone":"6388768118","dob":"02/07/2000","branch":"ECE",
"comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`, 1,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"},
			entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, nil,
			`{"id":"71bbdbb9-6bde-11ed-aaff-64bc589051b4","name":"Monika Jaiswal","phone":"6388768118",` +
				`"dob":"02/07/2000","branch":"ECE","comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`,
			201,
		},
		{"Error case: when id is valid but id is not present in db", id,
			`{"name":"Monika Jaiswal","phone":"6388768119","dob":"02/07/2000","branch":"ECE",
"comp":{"id":"1fa46d13-6a50-11ed-90d1-64bc589051b4"},"status":"ACCEPTED"}`, 1,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768119", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID},
				Status: "ACCEPTED"}, entities.Student{}, errors.DB{Reason: "server error"}, "DB Error: server error", 400,
		},
		{"Error case: missing parameters", id, `{}`, 0, entities.Student{},
			entities.Student{}, nil, "Missing Parameter: name,phone,dob,branch,company id,status", 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("PUT", "/students/{id}", strings.NewReader(tc.input))
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})

		h := New(mockStudent)
		mockStudent.EXPECT().Update(gomock.Any(), tc.inputID, &tc.mockInput).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)

		h.Update(resRec, req)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
		assert.Equal(t, tc.expRes, resRec.Body.String(), "Test[%v] failed\n(%v)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	mockStudent := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         error
		statusCode  int
	}{
		{"Success case: for valid id", id, nil, 204},
		{"Error case: when id is valid but id is not present in db", id,
			errors.DB{Reason: "server error"}, 400,
		},
	}

	for i, tc := range tests {
		req, err := http.NewRequest("DELETE", "/students/{id}", http.NoBody)
		if err != nil {
			t.Errorf(err.Error())
		}

		resRec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": tc.inputID.String()})

		h := New(mockStudent)
		mockStudent.EXPECT().Delete(gomock.Any(), tc.inputID).Return(tc.res)
		h.Delete(resRec, req)

		assert.Equal(t, tc.statusCode, resRec.Code, "Test[%v] failed\n(%v)", i, tc.description)
	}
}
