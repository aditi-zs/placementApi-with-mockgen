package student

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/Zopsmart-Training/student-recruitment-system/entities"
	"github.com/Zopsmart-Training/student-recruitment-system/errors"
)

type store struct {
	db *sql.DB
}

//nolint:revive // it's a factory function
func New(d *sql.DB) store {
	return store{db: d}
}

func (s store) GetByID(ctx context.Context, id uuid.UUID) (entities.Student, error) {
	var student entities.Student

	row := s.db.QueryRowContext(ctx, getByIDQuery, id)
	err := row.Scan(&student.ID, &student.Name, &student.Phone, &student.DOB, &student.Branch,
		&student.Comp.ID, &student.Comp.Name, &student.Comp.Category, &student.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Student{}, errors.EntityNotFound{Reason: "id not found"}
		}

		return entities.Student{}, errors.DB{Reason: "server error"}
	}

	return student, nil
}

func (s store) GetWithCompany(ctx context.Context, name, branch string) ([]entities.Student, error) {
	query := getDataWithCompQuery
	query += queryBuilder(name, branch)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []entities.Student{}, errors.DB{Reason: "server error"}
	}

	defer rows.Close()

	var students []entities.Student

	for rows.Next() {
		var student entities.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Phone, &student.DOB, &student.Branch,
			&student.Comp.ID, &student.Comp.Name, &student.Comp.Category, &student.Status)

		if err != nil {
			return []entities.Student{}, errors.DB{Reason: "scan error"}
		}

		students = append(students, student)
	}

	if rows.Err() != nil {
		return []entities.Student{}, errors.DB{Reason: "server error"}
	}

	if len(students) == 0 {
		return []entities.Student{}, errors.DB{Reason: "no rows found"}
	}

	return students, nil
}
func (s store) Get(ctx context.Context, name, branch string) ([]entities.Student, error) {
	query := getDataQuery
	query += queryBuilder(name, branch)
	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return []entities.Student{}, errors.DB{Reason: "server error"}
	}
	defer rows.Close()

	var students []entities.Student

	for rows.Next() {
		var student entities.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Phone, &student.DOB, &student.Branch, &student.Status)

		if err != nil {
			return []entities.Student{}, errors.DB{Reason: "scan error"}
		}

		students = append(students, student)
	}

	if rows.Err() != nil {
		return []entities.Student{}, errors.DB{Reason: "server error"}
	}

	if len(students) == 0 {
		return []entities.Student{}, errors.DB{Reason: "no rows found"}
	}

	return students, nil
}
func (s store) Create(ctx context.Context, st *entities.Student) (entities.Student, error) {
	st.ID = uuid.New()

	_, err := s.db.ExecContext(ctx, postQuery, st.ID,
		st.Name, st.Phone, st.DOB, st.Branch, st.Comp.ID, st.Status)
	if err != nil {
		return entities.Student{}, errors.DB{Reason: "server error"}
	}

	return *st, nil
}

func (s store) Update(ctx context.Context, id uuid.UUID, st *entities.Student) (entities.Student, error) {
	res, err := s.db.ExecContext(ctx, updateQuery,
		st.Name, st.Phone, st.DOB, st.Branch, st.Comp.ID, st.Status, id)
	if err != nil {
		return entities.Student{}, errors.DB{Reason: err.Error()}
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return entities.Student{}, errors.EntityNotFound{Reason: "id not found"}
	}

	st.ID = id

	return *st, nil
}

func (s store) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, deleteQuery, id)

	if err != nil {
		return errors.DB{Reason: err.Error()}
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.EntityNotFound{Reason: "id not found"}
	}

	return nil
}

func (s store) GetCompanyByID(ctx context.Context, id uuid.UUID) (entities.Company, error) {
	var company entities.Company

	row := s.db.QueryRowContext(ctx, getCompanyQuery, id)

	err := row.Scan(&company.ID, &company.Name, &company.Category)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Company{}, errors.EntityNotFound{Reason: "id not found"}
		}

		return entities.Company{}, errors.DB{Reason: "server error"}
	}

	return company, nil
}

func queryBuilder(name, branch string) string {
	var query string

	switch {
	case name != "" && branch != "":
		query = query + " " + fmt.Sprintf("where s.student_name='%v' AND s.branch='%v'", name, branch)
	case name != "" && branch == "":
		query = query + " " + fmt.Sprintf("where s.student_name='%v'", name)
	case name == "" && branch != "":
		query = query + " " + fmt.Sprintf("where s.branch='%v'", branch)
	}

	return query
}
