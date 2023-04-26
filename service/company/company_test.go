package company

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

func initializeTest(t *testing.T) *store.MockCompanyStore {
	ctrl := gomock.NewController(t)
	mockCompany := store.NewMockCompanyStore(ctrl)

	return mockCompany
}

func TestGet(t *testing.T) {
	mockCompany := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		res         []entities.Company
		err         error
	}{
		{"Success case: All entries are present",
			[]entities.Company{{ID: id, Name: "Wipro", Category: "MASS"}}, nil,
		},
		{"Error case: error in scanning",
			[]entities.Company{}, errors.DB{Reason: "scan error"},
		},
	}

	for i, tc := range tests {
		c := New(mockCompany)

		mockCompany.EXPECT().Get(context.Background()).Return(tc.res, tc.err)
		output, _ := c.Get(context.Background())

		assert.Equal(t, tc.res, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestGetByID(t *testing.T) {
	mockCompany := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         entities.Company
		err         error
	}{
		{"Success case: for valid id",
			id, entities.Company{ID: id, Name: "Wipro", Category: "MASS"}, nil,
		},
		{"Error case: when id is not present in db",
			id, entities.Company{}, errors.DB{Reason: "id not found"},
		},
	}

	for i, tc := range tests {
		c := New(mockCompany)
		mockCompany.EXPECT().GetByID(context.Background(), tc.inputID).Return(tc.res, tc.err)
		output, err := c.GetByID(context.Background(), tc.inputID)
		assert.Equal(t, tc.err, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.res, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestCreate(t *testing.T) {
	mockCompany := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		input       entities.Company
		mockTimes   int
		mockRes     entities.Company
		mockErr     error
		expRes      entities.Company
		expErr      error
	}{
		{"Success case: All entries are present",
			entities.Company{Name: "Wipro", Category: "MASS"}, 1,
			entities.Company{ID: id, Name: "Wipro", Category: "MASS"}, nil,
			entities.Company{ID: id, Name: "Wipro", Category: "MASS"}, nil,
		},
		{"Error case: when company category is different",
			entities.Company{Name: "Google", Category: "A"}, 0, entities.Company{}, nil,
			entities.Company{}, errors.InvalidParam{Param: "invalid category"},
		},
		{"Error case: server error",
			entities.Company{Name: "Infosys", Category: "MASS"}, 1, entities.Company{},
			errors.DB{Reason: "server error"}, entities.Company{},
			errors.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		c := New(mockCompany)

		mockCompany.EXPECT().Create(context.Background(), gomock.Any()).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		output, err := c.Create(context.Background(), tc.input)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	mockCompany := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		input       entities.Company
		mockTimes   int
		mockRes     entities.Company
		mockErr     error
		expRes      entities.Company
		expErr      error
	}{
		{"Success case: for valid id",
			id, entities.Company{Name: "Google", Category: "DREAM IT"}, 1,
			entities.Company{ID: id, Name: "Google", Category: "DREAM IT"}, nil,
			entities.Company{ID: id, Name: "Google", Category: "DREAM IT"}, nil,
		},
		{"Error case: when company category is different",
			id, entities.Company{Name: "Google", Category: "A"}, 0, entities.Company{},
			nil, entities.Company{}, errors.InvalidParam{Param: "invalid category"},
		},
		{"Error case: when id is not present in db",
			id, entities.Company{Name: "Google", Category: "DREAM IT"}, 1,
			entities.Company{}, errors.DB{Reason: "id not found"}, entities.Company{},
			errors.DB{Reason: "id not found"},
		},
	}

	for i, tc := range tests {
		c := New(mockCompany)

		mockCompany.EXPECT().Update(context.Background(), tc.inputID, tc.input).Return(tc.mockRes, tc.mockErr).Times(tc.mockTimes)
		output, err := c.Update(context.Background(), tc.inputID, tc.input)
		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	mockCompany := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		description string
		inputID     uuid.UUID
		res         error
	}{
		{"Success case: for valid id",
			id,
			nil,
		},
		{"Error case: when id is not present in db",
			id,
			errors.DB{Reason: "id not found"},
		},
	}

	for i, tc := range tests {
		c := New(mockCompany)

		mockCompany.EXPECT().Delete(context.Background(), tc.inputID).Return(tc.res)
		err := c.Delete(context.Background(), tc.inputID)

		assert.Equal(t, tc.res, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}
