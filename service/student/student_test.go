package student

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Zopsmart-Training/student-recruitment-system/entities"
	"github.com/Zopsmart-Training/student-recruitment-system/errors"
	"github.com/Zopsmart-Training/student-recruitment-system/store"
)

func initializeTest(t *testing.T) *store.MockStudentStore {
	ctrl := gomock.NewController(t)
	mockStudent := store.NewMockStudentStore(ctrl)

	return mockStudent
}

func TestGet(t *testing.T) {
	mockStudent := initializeTest(t)
	id := uuid.New()
	cmpID := uuid.New()
	tests := []struct {
		description         string
		queryName           string
		queryBranch         string
		queryIncludeCompany string
		mockTimes           int
		mockOP              []entities.Student
		mockErr             error
		expRes              []entities.Student
		expErr              error
	}{
		{"Success case: getting all records with company details", "", "ECE", "true", 1,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil,
		},
		{"Success case: when query params are valid", "Monika", "ECE", "true", 1,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil,
		},
		{"Success case: get data without company details", "Monika", "ECE", "false", 1,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{}, Status: "ACCEPTED"}}, nil,
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{}, Status: "ACCEPTED"}}, nil,
		},
		{"Error case: when query name is less than three characters", "ab", "", "", 0,
			nil, nil, []entities.Student{}, errors.InvalidParam{Param: "name should be minimum of three characters long"},
		},
		{"Error case: when query branch is invalid", "", "ABC", "", 0,
			nil, nil, []entities.Student{}, errors.InvalidParam{Param: "this branch is not allowed"},
		},
		{"Error case: when query includeCompany is invalid",
			"", "ECE", "ABC", 0, nil, nil,
			[]entities.Student{}, errors.InvalidParam{Param: "this value is not allowed"},
		},
		{"Error case: server error", "Utkarsh", "ECE", "true", 1,
			[]entities.Student{}, errors.DB{Reason: "server error"}, []entities.Student{}, errors.DB{Reason: "server error"},
		},
		{"Error case: server error", "Utkarsh", "ECE", "false", 1, []entities.Student{},
			errors.DB{Reason: "server error"}, []entities.Student{}, errors.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		s := New(mockStudent)

		if tc.queryIncludeCompany == "true" {
			mockStudent.EXPECT().GetWithCompany(context.Background(), tc.queryName, tc.queryBranch).Return(tc.mockOP, tc.mockErr).Times(tc.mockTimes)
		} else {
			mockStudent.EXPECT().Get(context.Background(), tc.queryName, tc.queryBranch).Return(tc.mockOP, tc.mockErr).Times(tc.mockTimes)
		}

		output, err := s.Get(context.Background(), tc.queryName, tc.queryBranch, tc.queryIncludeCompany)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestGetByID(t *testing.T) {
	mockStudent := initializeTest(t)
	id := uuid.New()
	cmpID := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         entities.Student
		err         error
	}{
		{"Success case: for valid id", id, entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil},
		{"Error case: when id is valid but id is not present in db", id, entities.Student{}, errors.DB{Reason: "id not found"}},
	}

	for i, tc := range tests {
		s := New(mockStudent)

		mockStudent.EXPECT().GetByID(context.Background(), tc.inputID).Return(tc.res, tc.err)
		output, err := s.GetByID(context.Background(), tc.inputID)

		assert.Equal(t, tc.res, output, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.err, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestCreate(t *testing.T) {
	mockStudent := initializeTest(t)
	cmpID := uuid.New()
	tests := []struct {
		description           string
		input                 entities.Student
		mockGetCompByIDTimes  int
		mockPostData          int
		mockGetCompanyByIDRes entities.Company
		mockGetCompanyByIDErr error
		mockPostDataRes       entities.Student
		mockPostDataErr       error
		expRes                entities.Student
		expErr                error
	}{
		{"Success case: All entries are present", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 1, 1,
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil,
			entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil,
		},
		{"Error case: When branch is different from given branches", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ABC", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"},
			0, 0, entities.Company{}, nil, entities.Student{},
			nil, entities.Student{}, errors.InvalidParam{Param: "this branch is not allowed"},
		},
		{"Error case: When name has less than 3 characters", entities.Student{Name: "Mn", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 0,
			0, entities.Company{}, nil, entities.Student{}, nil,
			entities.Student{}, errors.InvalidParam{Param: "name should be minimum of three characters long"},
		},
		{"Error case: When phone number has less than 10 numbers", entities.Student{Name: "Monika Jaiswal", Phone: "638876811",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"},
			0, 0, entities.Company{}, nil, entities.Student{},
			nil, entities.Student{}, errors.InvalidParam{Param: "phone number must be 10-12 digit long"},
		},
		{"Error case: When age is less than 22", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2010",
			Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 0, 0,
			entities.Company{}, nil, entities.Student{}, nil,
			entities.Student{}, errors.InvalidParam{Param: "age should be greater than 22"},
		},
		{"Error case: invalid status", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
			Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ABC"}, 0, 0, entities.Company{},
			nil, entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "invalid status"},
		},
		{"Error case: db error", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
			Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 1, 1,
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil, entities.Student{},
			errors.DB{Reason: "server error"}, entities.Student{}, errors.DB{Reason: "server error"},
		},
		{"Error case: when id not found", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "CSE",
			Comp: entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"}, Status: "ACCEPTED"}, 1, 0,
			entities.Company{}, errors.DB{Reason: "id not found"}, entities.Student{}, nil,
			entities.Student{}, errors.DB{Reason: "id not found"},
		},
		{"Error case: When branch is different from given branches for a company category", entities.Student{Name: "Monika Jaiswal",
			Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"},
			Status: "ACCEPTED"}, 1, 0, entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"},
			nil, entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "invalid branch for this company category"},
		},
		{"Error case: When branch is different from given branches for a company category", entities.Student{Name: "Monika Jaiswal",
			Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "TCS", Category: "CORE"},
			Status: "ACCEPTED"}, 1, 0, entities.Company{ID: cmpID, Name: "TCS", Category: "CORE"}, nil,
			entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "invalid branch for this company category"},
		},
		{"Error case: When branch is different from given branches for a company category", entities.Student{Name: "Monika Jaiswal",
			Phone: "6388768118", DOB: "02/07/2000", Branch: "MECH", Comp: entities.Company{ID: cmpID, Name: "ZopSmart", Category: "OPEN DREAM"},
			Status: "ACCEPTED"}, 1, 0, entities.Company{ID: cmpID, Name: "ZopSmart", Category: "OPEN DREAM"}, nil,
			entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "invalid branch for this company category"},
		},
	}

	for i, tc := range tests {
		s := New(mockStudent)
		mockStudent.EXPECT().GetCompanyByID(context.Background(), tc.input.Comp.ID).
			Return(tc.mockGetCompanyByIDRes, tc.mockGetCompanyByIDErr).Times(tc.mockGetCompByIDTimes)
		mockStudent.EXPECT().Create(context.Background(), &tc.input).
			Return(tc.mockPostDataRes, tc.mockPostDataErr).Times(tc.mockPostData)

		output, err := s.Create(context.Background(), &tc.input)
		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	mockStudent := initializeTest(t)
	id := uuid.New()
	cmpID := uuid.New()
	tests := []struct {
		description          string
		inputID              uuid.UUID
		input                entities.Student
		mockGetCompByIDTimes int
		mockUpdateDataTimes  int
		mockGetCompByIDRes   entities.Company
		mockGetCompByIDErr   error
		mockUpdateDataRes    entities.Student
		mockUpdateDataErr    error
		expRes               entities.Student
		expErr               error
	}{
		{"Success case: for valid id", id, entities.Student{Name: "Aditi Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
			Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 1, 1,
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil, entities.Student{ID: id, Name: "Aditi Jaiswal",
				Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"},
				Status: "ACCEPTED"}, nil, entities.Student{ID: id, Name: "Aditi Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil,
		},
		{"Error case: When branch is different from given branches", id, entities.Student{Name: "Aditi Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ABC", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"},
			0, 0, entities.Company{}, nil, entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "this branch is not allowed"},
		},
		{"Error case: When name has less than 3 characters", id, entities.Student{Name: "Ad", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "ABC", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 0, 0, entities.Company{},
			nil, entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "name should be minimum of three characters long"},
		},
		{"Error case: When phone number has less than 10 numbers", id, entities.Student{Name: "Aditi", Phone: "638876811", DOB: "02/07/2000",
			Branch: "ABC", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 0, 0, entities.Company{},
			nil, entities.Student{}, nil, entities.Student{}, errors.InvalidParam{Param: "phone number must be 10-12 digit long"},
		},
		{"Error case: when id is valid but id is not present in db", id, entities.Student{Name: "Aditi Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, 1, 1,
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil, entities.Student{}, errors.DB{Reason: "id not found"},
			entities.Student{}, errors.DB{Reason: "id not found"},
		},
		{"Error case: When branch is different from given branches for a company category", id, entities.Student{Name: "Aditi Jaiswal",
			Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"},
			Status: "ACCEPTED"}, 1, 0, entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"}, nil, entities.Student{},
			nil, entities.Student{}, errors.InvalidParam{Param: "invalid branch for this company category"},
		},
		{"Error case: db error", id, entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "CSE", Comp: entities.Company{ID: cmpID, Name: "MicroSoft", Category: "DREAM IT"}, Status: "ACCEPTED"}, 1, 0, entities.Company{},
			errors.DB{Reason: "id not found"}, entities.Student{}, nil, entities.Student{}, errors.DB{Reason: "id not found"},
		},
	}

	for i, tc := range tests {
		s := New(mockStudent)
		mockStudent.EXPECT().GetCompanyByID(context.Background(), tc.input.Comp.ID).
			Return(tc.mockGetCompByIDRes, tc.mockGetCompByIDErr).Times(tc.mockGetCompByIDTimes)
		mockStudent.EXPECT().Update(context.Background(), tc.inputID, &tc.input).
			Return(tc.mockUpdateDataRes, tc.mockUpdateDataErr).Times(tc.mockUpdateDataTimes)

		output, err := s.Update(context.Background(), tc.inputID, &tc.input)

		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	mockStudent := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         error
	}{
		{"for valid id", id, nil},
		{"when id is valid but id is not present in db", id, errors.DB{Reason: "id not found"}},
	}

	for i, tc := range tests {
		s := New(mockStudent)
		mockStudent.EXPECT().Delete(context.Background(), tc.inputID).Return(tc.res)
		err := s.Delete(context.Background(), tc.inputID)

		assert.Equal(t, tc.res, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}
