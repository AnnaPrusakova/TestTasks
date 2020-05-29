package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


type Employees struct {
	emp_no   int
	birth_date string
	first_name string
	last_name string
	gender string
	hire_date string
}

type Titles struct {
	emp_no int
	titile string
	from_date string
	to_date string
}

type Departments struct {
	dept_no string
	dept_name string
}

type Salaries struct {
	emp_no int
	salary int
	from_date string
	to_date string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "091040tardis"
	dbName := "employees"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func allManagers(){
	db := dbConn()

	// Execute the query
	results, err := db.Query("SELECT t.title, e.first_name, e.last_name, s.salary" +
		"\nFROM titles as t" +
		"\nINNER JOIN employees as e ON" +
		"\nt.emp_no = e.emp_no" +
		"\nINNER JOIN salaries as s ON" +
		"\ne.emp_no = s.emp_no" +
		"\nWHERE t.title=\"Manager\";")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var employess Employees
		var titles Titles
		var salary Salaries

		err = results.Scan(&titles.titile, &employess.first_name,&employess.last_name,&salary.salary)
		if err != nil {
			panic(err.Error())
		}

		log.Println(titles.titile,employess.first_name,employess.last_name,salary.salary)
	}
	db.Close()
}

func allEmployees(){
	db := dbConn()

	results, err := db.Query("SELECT d.dept_name, t.title, e.first_name, e.last_name, e.hire_date, year(DATE_SUB(current_date(), interval year(e.hire_date) YEAR)) as many_years_work" +
		"\nFROM departments as d" +
		"\nINNER JOIN dept_emp as de ON" +
		"\nd.dept_no = de.dept_no" +
		"\nINNER JOIN employees as e ON" +
		"\ne.emp_no = de.emp_no" +
		"\nINNER JOIN titles as t ON " +
		"\nt.emp_no = e.emp_no" +
		"\nWHERE month(e.hire_date) = \"05\";")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var employess Employees
		var depart Departments
		var title Titles
		var years int

		err = results.Scan(&depart.dept_name,&title.titile,&employess.first_name,&employess.last_name,&employess.hire_date,&years)
		if err != nil {
			panic(err.Error())
		}

		log.Println(depart.dept_name,employess.first_name,employess.last_name,employess.hire_date,years)
	}
	db.Close()
}


func allDepartments(){
	db:= dbConn()

	results, err := db.Query("SELECT d.dept_name, count(e.emp_no) as employees_count, sum(s.salary) as salary_sum" +
		"\nFROM departments as d\n" +
		"INNER JOIN dept_emp as de ON \n" +
		"  d.dept_no = de.dept_no\n" +
		"INNER JOIN employees as e ON \n" +
		"de.emp_no = e.emp_no \n" +
		"INNER JOIN salaries as s ON" +
		"\ne.emp_no = s.emp_no\n" +
		" GROUP BY d.dept_no;")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var depart Departments
		var employees_count int
		var salary_sum int

		err = results.Scan(&depart.dept_name,&employees_count,&salary_sum)
		if err != nil {
			panic(err.Error())
		}

		log.Println(depart.dept_name,employees_count,salary_sum)
	}
	db.Close()
}

func main() {

	allManagers()
	allEmployees()
	allDepartments()


}

