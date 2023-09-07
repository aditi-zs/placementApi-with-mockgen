package company

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/aditi-zs/Placement-API/entities"
	errors2 "github.com/aditi-zs/Placement-API/errors"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}

	defer db.Close()

	cmpID := uuid.New()

	tests := []struct {
		description string
		rows        *sqlmock.Rows
		expRes      []entities.Company
		mockErr     error
		expErr      error
	}{
		{"Success case: All entries are present",
			sqlmock.NewRows([]string{"ID", "Name", "Category"}).AddRow(cmpID.String(), "Wipro", "MASS"),
			[]entities.Company{{ID: cmpID, Name: "Wipro", Category: "MASS"}},
			nil, nil,
		},
		{"Error case: server error",
			sqlmock.NewRows([]string{"ID", "Name", "Category"}).AddRow(cmpID, "Wipro", "MASS"),
			[]entities.Company{}, errors.New("no rows found"), errors2.DB{Reason: "no rows found"},
		},
		{"Error case: scan error",
			sqlmock.NewRows([]string{"ID", "Name", "Category"}).AddRow(nil, nil, nil),
			[]entities.Company{}, nil, errors2.DB{Reason: "scan error"},
		},
	}
	for i, tc := range tests {
		store := New(db)

		mock.ExpectQuery(getQuery).WillReturnRows(tc.rows).WillReturnError(tc.mockErr)

		ctx := context.TODO()
		output, err := store.Get(ctx)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmpID := uuid.New()
	valID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		rows        *sqlmock.Rows
		expRes      entities.Company
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id", cmpID,
			sqlmock.NewRows([]string{"ID", "Name", "category"}).AddRow(cmpID, "Wipro", "MASS"),
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil, nil,
		},
		{"Error case: when id is not present in db", valID,
			sqlmock.NewRows([]string{"ID", "Name", "category"}), entities.Company{},
			nil, errors2.EntityNotFound{Reason: "id not found: " + valID.String()},
		},
		{"error case: server error", cmpID,
			sqlmock.NewRows([]string{"ID", "Name", "category"}).AddRow(cmpID, "Wipro", "MASS"),
			entities.Company{}, errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		mock.ExpectQuery(getByIDQuery).WithArgs(tc.inputID).
			WillReturnRows(tc.rows).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.GetByID(ctx, tc.inputID)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmpID := uuid.New()

	tests := []struct {
		description string
		input       entities.Company
		res         driver.Result
		expRes      entities.Company
		err         error
	}{
		{"Success case: All entries are present", entities.Company{Name: "Wipro", Category: "MASS"}, sqlmock.NewResult(1, 1),
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil,
		},
		{"Error case :server error", entities.Company{Name: "Wipro", Category: "MASS"},
			sqlmock.NewResult(0, 0), entities.Company{}, errors2.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		store := New(db)

		mock.ExpectExec(postQuery).WithArgs(sqlmock.AnyArg(), tc.input.Name, tc.input.Category).
			WillReturnResult(tc.res).WillReturnError(tc.err)

		ctx := context.TODO()
		output, err := store.Create(ctx, tc.input)

		assert.Equal(t, tc.err, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Name, output.Name, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Category, output.Category, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmpID := uuid.New()
	valID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		input       entities.Company
		res         driver.Result
		expRes      entities.Company
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id", cmpID, entities.Company{ID: cmpID, Name: "Google", Category: "DREAM IT"},
			sqlmock.NewResult(1, 1), entities.Company{ID: cmpID, Name: "Google", Category: "DREAM IT"}, nil, nil,
		},
		{"Error case: when id is valid but it doesn't exist in db", valID, entities.Company{Name: "Google", Category: "DREAM IT"},
			sqlmock.NewResult(0, 0), entities.Company{}, nil, errors2.EntityNotFound{Reason: "id not found: " + valID.String()},
		},
		{"Error case: invalid data", cmpID, entities.Company{ID: cmpID, Name: "", Category: "DREAM IT"},
			sqlmock.NewResult(0, 0), entities.Company{}, errors.New("invalid data"), errors2.DB{Reason: "invalid data"},
		},
	}
	for i, tc := range tests {
		c := New(db)

		mock.ExpectExec(updateQuery).
			WithArgs(tc.input.Name, tc.input.Category, tc.inputID).
			WillReturnResult(tc.res).WillReturnError(tc.mockErr)

		ctx := context.TODO()
		output, err := c.Update(ctx, tc.inputID, tc.input)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmpID := uuid.New()
	valID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		res         driver.Result
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id",
			cmpID, sqlmock.NewResult(1, 1), nil, nil,
		},
		{"Error case: when id is valid but it doesn't exist in db", valID, sqlmock.NewResult(0, 0), nil,
			errors2.EntityNotFound{Reason: "id not found: " + valID.String()},
		},
		{"Error case: when id is used as foreign key", valID, sqlmock.NewResult(0, 0),
			errors.New("this id is used as a foreign key"), errors2.DB{Reason: "this id is used as a foreign key"},
		},
	}

	for i, tc := range tests {
		mock.ExpectExec(deleteQuery).WithArgs(tc.inputID).WillReturnResult(tc.res).WillReturnError(tc.mockErr)

		c := New(db)
		ctx := context.TODO()
		err := c.Delete(ctx, tc.inputID)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}
