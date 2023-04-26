package student

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Zopsmart-Training/student-recruitment-system/entities"
	errors2 "github.com/Zopsmart-Training/student-recruitment-system/errors"
)

func TestGetWithCompany(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()
	cmpID := uuid.New()

	tests := []struct {
		description string
		inputName   string
		inputBranch string
		queryR      string
		rows        *sqlmock.Rows
		expRes      []entities.Student
		mockErr     error
		expErr      error
	}{
		{"Success case: All entries are present", "Monika", "ECE", getDataWithCompQuery + " where s.student_name='Monika' AND s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, "Monika Jaiswal", "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when only name is present as query param", "Monika", "", getDataWithCompQuery + " where s.student_name='Monika'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, "Monika Jaiswal", "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when only branch is present as query param", "", "ECE", getDataWithCompQuery + " where s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, "Monika Jaiswal", "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when no query param is present", "", "", getDataWithCompQuery,
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, "Monika Jaiswal", "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case", "Monika", "E", getDataWithCompQuery + " where s.student_name='Monika' AND s.branch='E'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}),
			[]entities.Student{}, nil, errors2.DB{Reason: "no rows found"},
		},
		{"Error case", "Monika", "ECE", getDataWithCompQuery + " where s.student_name='Monika' AND s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, nil, "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			[]entities.Student{}, nil, errors2.DB{Reason: "scan error"},
		},
		{"failure case", "Monika", "ECE", getDataWithCompQuery + " where s.student_name='Monika' AND s.branch='E'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}),
			[]entities.Student{}, errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}
	for i, tc := range tests {
		mock.ExpectQuery(tc.queryR).WillReturnRows(tc.rows).WillReturnError(tc.mockErr).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.GetWithCompany(ctx, tc.inputName, tc.inputBranch)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()

	tests := []struct {
		description string
		inputName   string
		inputBranch string
		queryR      string
		rows        *sqlmock.Rows
		expRes      []entities.Student
		mockErr     error
		expErr      error
	}{
		{"Success case: All entries are present", "Aditi", "ECE", getDataQuery + " where s.student_name='Aditi' AND s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}).
				AddRow(id, "Aditi", "6388768119", "02/03/2000", "ECE", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Aditi", Phone: "6388768119", DOB: "02/03/2000", Branch: "ECE", Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when only name is present as query param", "Aditi", "", getDataQuery + " where s.student_name='Aditi'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}).
				AddRow(id, "Aditi", "6388768119", "02/03/2000", "ECE", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Aditi", Phone: "6388768119", DOB: "02/03/2000", Branch: "ECE", Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when only branch is present as query param", "", "ECE", getDataQuery + " where s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}).
				AddRow(id, "Aditi", "6388768119", "02/03/2000", "ECE", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Aditi", Phone: "6388768119", DOB: "02/03/2000", Branch: "ECE", Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case: when no query params present", "", "", getDataQuery,
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}).
				AddRow(id, "Aditi", "6388768119", "02/03/2000", "ECE", "ACCEPTED"),
			[]entities.Student{{ID: id, Name: "Aditi", Phone: "6388768119", DOB: "02/03/2000", Branch: "ECE", Status: "ACCEPTED"}}, nil, nil,
		},
		{"Error case", "Aditi", "E", getDataQuery + " where s.student_name='Aditi' AND s.branch='E'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}), []entities.Student{},
			nil, errors2.DB{Reason: "no rows found"},
		},
		{"Error case", "Monika", "ECE", getDataQuery + " where s.student_name='Monika' AND s.branch='ECE'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}).
				AddRow(id, nil, "6388768118", "02/07/2000", "ECE", "ACCEPTED"), []entities.Student{}, nil, errors2.DB{Reason: "scan error"},
		},
		{"Error case", "Aditi", "E", getDataQuery + " where s.student_name='Aditi' AND s.branch='E'",
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "status"}), []entities.Student{},
			errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		mock.ExpectQuery(tc.queryR).WillReturnRows(tc.rows).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.Get(ctx, tc.inputName, tc.inputBranch)

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

	id := uuid.New()
	cmpID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		rows        *sqlmock.Rows
		expRes      entities.Student
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id", id,
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}).
				AddRow(id, "Monika Jaiswal", "6388768118", "02/07/2000", "ECE", cmpID, "Wipro", "MASS", "ACCEPTED"),
			entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
				Comp: entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, Status: "ACCEPTED"}, nil, nil,
		},
		{"Error case : when id is not present in db", id,
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}),
			entities.Student{}, sql.ErrNoRows, errors2.EntityNotFound{Reason: "id not found"},
		},
		{"Error case: server error", id,
			sqlmock.NewRows([]string{"ID", "Name", "Phone", "dob", "branch", "ID", "Name", "category", "status"}),
			entities.Student{}, errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}
	for i, tc := range tests {
		mock.ExpectQuery(getByIDQuery).
			WithArgs(tc.inputID).WillReturnRows(tc.rows).WillReturnError(tc.mockErr)

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

	id := uuid.New()
	cmpID := uuid.New()

	tests := []struct {
		description string
		input       entities.Student
		res         driver.Result
		expRes      entities.Student
		mockErr     error
		expErr      error
	}{
		{"Success case: All entries are present", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
			Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, sqlmock.NewResult(1, 1),
			entities.Student{ID: id, Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, nil, nil,
		},
		{"Error case: server error", entities.Student{Name: "Monika Jaiswal", Phone: "6388768118", DOB: "02/07/2000", Branch: "ECE",
			Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, sqlmock.NewResult(0, 0),
			entities.Student{}, errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}
	for i, tc := range tests {
		mock.ExpectExec(postQuery).
			WithArgs(sqlmock.AnyArg(), tc.input.Name, tc.input.Phone, tc.input.DOB, tc.input.Branch,
				tc.input.Comp.ID, tc.input.Status).WillReturnResult(tc.res).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.Create(ctx, &tc.input)

		assert.Equal(t, tc.expRes.Name, output.Name, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Phone, output.Phone, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Branch, output.Branch, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.DOB, output.DOB, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Comp, output.Comp, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes.Status, output.Status, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()
	cmpID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		input       entities.Student
		res         driver.Result
		expRes      entities.Student
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id", id, entities.Student{ID: id, Name: "Aditi Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, sqlmock.NewResult(1, 1),
			entities.Student{ID: id, Name: "Aditi Jaiswal", Phone: "6388768118", DOB: "02/07/2000",
				Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, nil, nil,
		},
		{"Error case: when id is valid but it doesn't exist in db", id, entities.Student{Name: "Aditi Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"}, sqlmock.NewResult(0, 0),
			entities.Student{}, nil, errors2.EntityNotFound{Reason: "id not found"},
		},
		{"Error case: when company id is foreign key", id, entities.Student{Name: "Aditi Jaiswal", Phone: "6388768118",
			DOB: "02/07/2000", Branch: "ECE", Comp: entities.Company{ID: cmpID}, Status: "ACCEPTED"},
			sqlmock.NewResult(0, 0), entities.Student{}, errors.New("this id is used as a foreign key"),
			errors2.DB{Reason: "this id is used as a foreign key"},
		},
	}
	for i, tc := range tests {
		mock.ExpectExec(updateQuery).
			WithArgs(tc.input.Name, tc.input.Phone, tc.input.DOB, tc.input.Branch,
				tc.input.Comp.ID, tc.input.Status, tc.inputID).
			WillReturnResult(tc.res).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.Update(ctx, tc.inputID, &tc.input)

		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		res         driver.Result
		sqlErr      error
		expErr      error
	}{{"Success case: for valid id", id, sqlmock.NewResult(1, 1), nil, nil},
		{"Error case: when id is valid but it doesn't exist in db", id, sqlmock.NewResult(0, 0), nil,
			errors2.EntityNotFound{Reason: "id not found"},
		},
		{"Error case: when id is used as foreign key", id,
			sqlmock.NewResult(0, 0), errors.New("this id is used as a foreign key"),
			errors2.DB{Reason: "this id is used as a foreign key"},
		},
	}
	for i, tc := range tests {
		mock.ExpectExec(deleteQuery).WithArgs(tc.inputID).WillReturnResult(tc.res).WillReturnError(tc.sqlErr)

		store := New(db)
		ctx := context.TODO()
		err := store.Delete(ctx, tc.inputID)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
	}
}

func TestGetCompanyByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()
	cmpID := uuid.New()

	tests := []struct {
		description string
		inputID     uuid.UUID
		rows        *sqlmock.Rows
		expRes      entities.Company
		mockErr     error
		expErr      error
	}{
		{"Success case: for valid id", cmpID, sqlmock.NewRows([]string{"ID", "Name", "category"}).AddRow(cmpID, "Wipro", "MASS"),
			entities.Company{ID: cmpID, Name: "Wipro", Category: "MASS"}, nil, nil,
		},
		{"Error case: when id is not present in db", id, sqlmock.NewRows([]string{"ID", "Name", "category"}),
			entities.Company{}, nil, errors2.EntityNotFound{Reason: "id not found"},
		},
		{"Error case: server error", cmpID, sqlmock.NewRows([]string{"ID", "Name", "category"}).AddRow(cmpID, "Wipro", "MASS"),
			entities.Company{}, errors.New("server error"), errors2.DB{Reason: "server error"},
		},
	}

	for i, tc := range tests {
		mock.ExpectQuery(getCompanyQuery).WithArgs(tc.inputID).
			WillReturnRows(tc.rows).WillReturnError(tc.mockErr)

		store := New(db)
		ctx := context.TODO()
		output, err := store.GetCompanyByID(ctx, tc.inputID)

		assert.Equal(t, tc.expErr, err, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, output, "Test[%d] failed\n(%s)", i, tc.description)
	}
}
