package company

const (
	getQuery     = "SELECT company_id,company_name,category from companies"
	getByIDQuery = "SELECT c.company_id,c.company_name,c.category from companies c where c.company_id=?"
	postQuery    = "INSERT INTO companies values (?,?,?)"
	updateQuery  = "UPDATE companies SET company_name=?,category=? WHERE company_id=?"
	deleteQuery  = "DELETE FROM companies  WHERE company_id=?"
)
