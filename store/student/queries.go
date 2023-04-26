package student

const (
	getByIDQuery = "select s.student_id,s.student_name,s.student_phone,s.dob,s.branch," +
		"c.company_id,c.company_name,c.category,s.status from students s join companies c on s. company_id=c. company_id " +
		"where s.student_id=?"
	getDataWithCompQuery = "SELECT s.student_id,s.student_name,s.student_phone,s.dob,s.branch,c.company_id,c.company_name," +
		"c.category,s.status from students s join companies c on s. company_id=c. company_id"
	getDataQuery    = "SELECT student_id,student_name,student_phone,dob,branch,status From students s"
	postQuery       = "INSERT INTO students values (?,?,?,?,?,?,?)"
	updateQuery     = "UPDATE students SET student_name=?,student_phone=?,dob=?,branch=?,company_id=?,status=? WHERE student_id=?"
	deleteQuery     = "DELETE FROM students WHERE student_id=?"
	getCompanyQuery = "SELECT c.company_id,c.company_name,c.category from companies c where c.company_id=?"
)
