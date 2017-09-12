package main

//Dig into this https://github.com/heroku/go-getting-started

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	//"log"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/foo", foo)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "FooBot Indoor Air Quality Metrics:\n")
	
	//Open database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sensor_data")
//	checkErr(err)

	fmt.Fprintln(res, "device_uuid | seconds_since_epoch | datetime | particles_ugm3 | temp_c | humidity_perct | co2_ppm | voc_ppb | pollution_score ")

	for rows.Next() {
		var device_uuid string
		var seconds_since_epoch int64
		var datetime time.Time
		var particles_ugm3 float64
		var temp_c float64
		var humidity_perct float64
		var co2_ppm float64
		var voc_ppb float64
		var pollution_score float64
		err = rows.Scan(&device_uuid, &seconds_since_epoch, &datetime, &particles_ugm3, &temp_c, &humidity_perct, &co2_ppm, &voc_ppb, &pollution_score)
		//checkErr(err)
		//fmt.Fprintln(res, "%3v | %8v | %6v | %6v\n", device_uuid, seconds_since_epoch, datetime, particles_ugm3)
		fmt.Fprintln(res, device_uuid, seconds_since_epoch, datetime, particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score)
	}

}

func foo(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "This is a test")
}