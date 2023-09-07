package student

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/aditi-zs/Placement-API/entities"
	"github.com/aditi-zs/Placement-API/errors"
	"github.com/aditi-zs/Placement-API/store"
)

type handler struct {
	datastore store.StudentStore
}

//nolint:revive // it's a factory function
func New(student store.StudentStore) handler {
	return handler{datastore: student}
}

func (s handler) GetByID(ctx context.Context, id uuid.UUID) (entities.Student, error) {
	resp, err := s.datastore.GetByID(ctx, id)
	if err != nil {
		return entities.Student{}, err
	}

	return resp, nil
}

func (s handler) Get(ctx context.Context, name, branch, includeCompany string) ([]entities.Student, error) {
	if err := validateQuery(name, branch, includeCompany); err != nil {
		return []entities.Student{}, err
	}

	if includeCompany == trueVal {
		resp, err := s.datastore.GetWithCompany(ctx, name, branch)
		if err != nil {
			return []entities.Student{}, err
		}

		return resp, nil
	}

	resp, err := s.datastore.Get(ctx, name, branch)
	if err != nil {
		return []entities.Student{}, err
	}

	return resp, nil
}

func (s handler) Create(ctx context.Context, st *entities.Student) (entities.Student, error) {
	if err := validateStudent(st); err != nil {
		return entities.Student{}, err
	}

	companyID := st.Comp.ID

	company, err := s.datastore.GetCompanyByID(ctx, companyID)
	if err != nil {
		return entities.Student{}, err
	}

	err = validateBranch(company.Category, st.Branch)
	if err != nil {
		return entities.Student{}, err
	}

	resp, err := s.datastore.Create(ctx, st)
	if err != nil {
		return entities.Student{}, err
	}

	return resp, nil
}
func (s handler) Update(ctx context.Context, id uuid.UUID, st *entities.Student) (entities.Student, error) {
	if err := validateStudent(st); err != nil {
		return entities.Student{}, err
	}

	companyID := st.Comp.ID

	company, err := s.datastore.GetCompanyByID(ctx, companyID)
	if err != nil {
		return entities.Student{}, err
	}

	err = validateBranch(company.Category, st.Branch)
	if err != nil {
		return entities.Student{}, err
	}

	resp, err := s.datastore.Update(ctx, id, st)

	if err != nil {
		return entities.Student{}, err
	}

	return resp, nil
}

func (s handler) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.datastore.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

const trueVal, mech, ise, civil, cse, ece, eee = "true", "MECH", "ISE", "CIVIL", "CSE", "ECE", "EEE"

func getAge(dob string) int {
	date, _ := time.Parse("02/01/2006", dob)
	today := time.Now()

	return today.Year() - date.Year()
}
func validateQuery(name, branch, includeCompany string) error {
	switch {
	case name != "" && len(name) < 3:
		return errors.InvalidParam{Param: "name should be minimum of three characters long"}
	case !IsValidBranch(branch):
		return errors.InvalidParam{Param: "this branch is not allowed"}
	case includeCompany != "" && includeCompany != trueVal && includeCompany != "false":
		return errors.InvalidParam{Param: "this value is not allowed"}
	default:
		return nil
	}
}

const minAge, nameLen, minphnlen, maxPhnLen = 22, 3, 10, 12

func validateStudent(stu *entities.Student) error {
	switch {
	case len(stu.Name) < nameLen:
		return errors.InvalidParam{Param: "name should be minimum of three characters long"}
	case len(stu.Phone) < minphnlen || len(stu.Phone) > maxPhnLen:
		return errors.InvalidParam{Param: "phone number must be 10-12 digit long"}
	case !entities.IsValidBranch(stu.Branch):
		return errors.InvalidParam{Param: "this branch is not allowed"}
	case getAge(stu.DOB) < minAge:
		return errors.InvalidParam{Param: "age should be greater than 22"}
	case !entities.IsValidStatus(stu.Status):
		return errors.InvalidParam{Param: "invalid status"}
	default:
		return nil
	}
}

//nolint:gocognit,gocyclo    //this function has many conditions to check
func validateBranch(category entities.Category, branch entities.Branch) error {
	switch {
	case category == "CORE" && branch != civil && branch != mech:
		return errors.InvalidParam{Param: "invalid branch for this company category"}
	case category == "OPEN DREAM" && branch != cse && branch != ise && branch != ece && branch != eee:
		return errors.InvalidParam{Param: "invalid branch for this company category"}
	case category == "DREAM IT" && branch != cse && branch != ise:
		return errors.InvalidParam{Param: "invalid branch for this company category"}
	default:
		return nil
	}
}

func IsValidBranch(b string) bool {
	switch b {
	case cse, ise, mech, ece, eee, civil, "":
		return true
	default:
		return false
	}
}
