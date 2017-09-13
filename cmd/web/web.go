package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"html/template"
	//"log"
	"os"
	//"time"
)

type FooData struct {
	Datetime string
	Seconds_since_epoch int64
	Particles_ugm3 float64
	Temp_c float64
	Humidity_perct float64
	Co2_ppm float64
	Voc_ppb float64
	Pollution_score float64
}

func main() {
	http.HandleFunc("/", index)
	//http.HandleFunc("/foo", foo)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
//		http.ServeFile(w, r, r.URL.Path[1:])
//	})
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	//Open database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT datetime::timestamp at time zone 'UTC' at time zone 'America/New_York', seconds_since_epoch, particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score FROM sensor_data")

	data := []FooData{}

	foo := FooData{}

	for rows.Next() {
		//var device_uuid string
		var datetime string
		var seconds_since_epoch int64
		var particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score float64
		err = rows.Scan(&datetime, &seconds_since_epoch, &particles_ugm3, &temp_c, &humidity_perct, &co2_ppm, &voc_ppb, &pollution_score)

		foo.Datetime = datetime
		foo.Seconds_since_epoch = seconds_since_epoch
		foo.Particles_ugm3 = particles_ugm3
		foo.Temp_c = temp_c
		foo.Humidity_perct = humidity_perct
		foo.Co2_ppm = co2_ppm
		foo.Voc_ppb = voc_ppb
		foo.Pollution_score = pollution_score
		data = append(data, foo)
		//fmt.Println(datetime, seconds_since_epoch, particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/header.html"))

	tmpl.Execute(res, data)

}
/*
func foo(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(res, "This is a test")

	data := "hello world"

	tmpl := template.Must(template.ParseFiles("test.html", "header.html"))

	tmpl.Execute(res, data)
}
*/