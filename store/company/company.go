package company

import (
	"context"
	"database/sql"

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

func (c store) GetByID(ctx context.Context, id uuid.UUID) (entities.Company, error) {
	var company entities.Company

	row := c.db.QueryRowContext(ctx, getByIDQuery, id)

	err := row.Scan(&company.ID, &company.Name, &company.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Company{}, errors.EntityNotFound{Reason: "id not found: " + id.String()}
		}

		return entities.Company{}, errors.DB{Reason: "server error"}
	}

	return company, nil
}
func (c store) Get(ctx context.Context) ([]entities.Company, error) {
	rows, err := c.db.QueryContext(ctx, getQuery)
	if err != nil {
		return []entities.Company{}, errors.DB{Reason: "no rows found"}
	}

	defer rows.Close()

	var companies []entities.Company

	for rows.Next() {
		var company entities.Company

		err = rows.Scan(&company.ID, &company.Name, &company.Category)
		if err != nil {
			return []entities.Company{}, errors.DB{Reason: "scan error"}
		}

		companies = append(companies, company)
	}

	return companies, nil
}
func (c store) Create(ctx context.Context, cmp entities.Company) (entities.Company, error) {
	cmp.ID = uuid.New()

	_, err := c.db.ExecContext(ctx, postQuery, cmp.ID, cmp.Name, cmp.Category)
	if err != nil {
		return entities.Company{}, errors.DB{Reason: "server error"}
	}

	return cmp, nil
}

func (c store) Update(ctx context.Context, id uuid.UUID, cmp entities.Company) (entities.Company, error) {
	res, err := c.db.ExecContext(ctx, updateQuery, cmp.Name, cmp.Category, id)
	if err != nil {
		return entities.Company{}, errors.DB{Reason: err.Error()}
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return entities.Company{}, errors.EntityNotFound{Reason: "id not found: " + id.String()}
	}

	cmp.ID = id

	return cmp, nil
}

func (c store) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := c.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return errors.DB{Reason: err.Error()}
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.EntityNotFound{Reason: "id not found: " + id.String()}
	}

	return nil
}
