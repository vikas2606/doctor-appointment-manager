package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var fname []string
var lname []string
var mobile []string
var OPD []string
var issue []string
var date []string
var Time []string
var email []string
var sex []string

func main() {
	http.HandleFunc("/", appoint)

	http.HandleFunc("/final", final)

	http.HandleFunc("/appointOrtho", appointOrtho)

	http.HandleFunc("/appointDiabet", appointDiabet)

	http.HandleFunc("/appointPedia", appointPedia)
	conn, err := sql.Open("mysql", "root:root@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}

	query := "CREATE TABLE IF NOT EXISTS patients(PID int NOT NULL AUTO_INCREMENT,fname varchar(200),lname varchar(200),sex varchar(200),Mobile varchar(200),email varchar(200),OPD varchar(200),Issue varchar(200),Date varchar(200),Time varchar(200),PRIMARY KEY(PID))"
	create, err := conn.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(create)
	defer conn.Close()
	for conn.Ping() != nil {
		fmt.Println("attempting connection to db")
		time.Sleep(5 * time.Second)

	}

	http.ListenAndServe(":8080", nil)

}
func appoint(w http.ResponseWriter, r *http.Request) {
	tmplt := template.Must(template.ParseFiles("templates/appoint.html"))
	tmplt.Execute(w, nil)
}
func final(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fname = r.Form["fname"]
	lname = r.Form["lname"]
	mobile = r.Form["mobile"]
	email = r.Form["email"]
	OPD = r.Form["OPD"]
	sex = r.Form["sex"]
	issue = r.Form["issue"]
	fmt.Println(fname, lname, mobile, OPD[0], issue)

	if check_user(mobile[0], email[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/already.html"))
		tmplt.Execute(w, nil)
	} else {

		if OPD[0] == "Orthopedic" {
			var tmplt = template.Must(template.ParseFiles("templates/appointOrtho.html"))
			tmplt.Execute(w, nil)

		} else if OPD[0] == "Diabetes" {
			var tmplt = template.Must(template.ParseFiles("templates/appointDiabet.html"))
			tmplt.Execute(w, nil)

		} else {
			var tmplt = template.Must(template.ParseFiles("templates/appointPedia.html"))
			tmplt.Execute(w, nil)

		}

	}

}
func appointDiabet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	fmt.Println(Time, date)

	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		mail(email[0], fname[0], lname[0], sex[0], "Naga Vikas", OPD[0], Time[0], date[0])
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Naga Vikas",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}
func appointOrtho(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		mail(email[0], fname[0], lname[0], sex[0], "Praveen Kumar", OPD[0], Time[0], date[0])
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Praveen Kumar",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}
func appointPedia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		mail(email[0], fname[0], lname[0], sex[0], "Khadar Basha", OPD[0], Time[0], date[0])
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Khadar Basha",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}

func check_user(mobile string, email string) bool {
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT fname FROM patients WHERE Mobile='%s' AND email='%s')", (mobile), (email))
	row := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println(row)
	defer db.Close()
	return exists
}
func add_user(fname string, lname string, sex string, mobile string, email string, OPD string, issue string, time string, date string) bool {
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}
	add, err := db.Query("INSERT INTO patients(fname,lname,sex,Mobile,email,OPD,Issue,Date,Time) VALUES (?,?,?,?,?,?,?,?,?)", (fname), (lname), (sex), (mobile), (email), (OPD), (issue), (date), (time))
	if err != nil {
		panic(err)
	}
	fmt.Println(add)
	defer db.Close()
	return true
}
func check_appoint(time string, date string) bool {
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT OPD FROM patients WHERE Time='%s' AND Date='%s')", (time), (date))
	row := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println(row)
	defer db.Close()
	return exists

}
func mail(email string, fname string, lname string, sex string, doctor string, OPD string, time string, date string) bool {
	from := "v26kas@gmail.com"
	password := "Vikas@123"
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := "Your appointment is successful\n DETAILS:\nName:" + fname + " " + lname + "\nSex:" + sex + "\nDoctor:" + doctor + "\nSpecialist" + OPD + "\nslot:" + time + "\ndate:" + date + "\nAppointID:" + PID(email)
	message := []byte(msg)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Email Sent Successfully!")
	return true
}
func PID(email string) string {
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}
	var PID string
	err = db.QueryRow("select PID from patients where email = ?", email).Scan(&PID)
	if err != nil {
		log.Fatal(err)
	}
	return PID

}
