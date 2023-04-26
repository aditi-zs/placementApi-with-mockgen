package company

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Zopsmart-Training/student-recruitment-system/entities"
	errors2 "github.com/Zopsmart-Training/student-recruitment-system/errors"
	"github.com/Zopsmart-Training/student-recruitment-system/store"
)

type handler struct {
	datastore store.CompanyStore
}

//nolint:revive // it's a factory function
func New(company store.CompanyStore) handler {
	return handler{datastore: company}
}

func (c handler) GetByID(ctx context.Context, id uuid.UUID) (entities.Company, error) {
	resp, err := c.datastore.GetByID(ctx, id)

	if err != nil {
		return entities.Company{}, errors2.DB{Reason: "id not found"}
	}

	return resp, nil
}
func (c handler) Get(ctx context.Context) ([]entities.Company, error) {
	resp, err := c.datastore.Get(ctx)
	if err != nil {
		return []entities.Company{}, err
	}

	return resp, nil
}
func (c handler) Create(ctx context.Context, cmp entities.Company) (entities.Company, error) {
	if err := validateCompany(cmp); err != nil {
		return entities.Company{}, errors2.InvalidParam{Param: "invalid category"}
	}

	resp, err := c.datastore.Create(ctx, cmp)
	if err != nil {
		return entities.Company{}, err
	}

	return resp, nil
}

func (c handler) Update(ctx context.Context, id uuid.UUID, cmp entities.Company) (entities.Company, error) {
	if err := validateCompany(cmp); err != nil {
		return entities.Company{}, errors2.InvalidParam{Param: "invalid category"}
	}

	resp, err := c.datastore.Update(ctx, id, cmp)
	if err != nil {
		return entities.Company{}, err
	}

	return resp, nil
}

func (c handler) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.datastore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func validateCompany(cmp entities.Company) error {
	if !entities.IsValidCategory(cmp.Category) {
		return errors.New("this category is invalid")
	}

	return nil
}
