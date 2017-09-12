// FooBot API Client pulls data from the FooBot API and inserts to a database.
package main

/*  To Do
* Separate into functions?
* Error handling & logging
* NOT NEEDED FOR GIT DEPLOYMENT, BUT DESCRIBE: Document godep install (go get github.com/tools/godep) and 
* Terraform & document the device ID & API key parameter setup: FOOBOT_DEVICE_UUID, FOOBOT_API_KEY
 */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//JSON to Struct Converter: http://json2struct.mervine.net/
type FoobotJson struct {
	UUID       string      `json:"uuid"`
	Start      int         `json:"start"`
	End        int         `json:"end"`
	Sensors    []string    `json:"sensors"`
	Units      []string    `json:"units"`
	Datapoints [][]float64 `json:"datapoints"`
}

// getLastDBInsert checks database for prior inserted data, if found, returns the
// most recent seconds since epoch value, otherwise returns 0.
func getLastDBInsert() int64 {

	//Open database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Query DB for last imported dataset, based on epoch time.
	var sse int64 = 0
	row := db.QueryRow("select seconds_since_epoch as sse from sensor_data order by sse desc limit 1;")
	switch err := row.Scan(&sse); err {
	case sql.ErrNoRows:
		return 0
	case nil:
		return sse
	default:
		//panic(err)
		return 0
	}
}

func main() {

	//Get Foobot Device ID and API key from environment variables
	fooDeviceUUID := os.Getenv("FOOBOT_DEVICE_UUID")
	fooApiKey := os.Getenv("FOOBOT_API_KEY")
	//************************ADD LOGGING ENTRY HERE TO RECORD DEVICE ID?

	//Open database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Set startEpoch default to 41 days ago. Foobot API limits returned data to 42 days.
	startEpoch := time.Now().AddDate(0, 0, -41).Unix()
	//Use more recent value if data has previously been inserted into database, based on seconds_since_epoch eval above
	sse := getLastDBInsert()
	if sse > startEpoch {
		startEpoch = sse
	}
	currentEpoch := time.Now().Unix()
	//*******************ADD LOGGING ENTRIES FOR START & CURRENT EPOCHS USED, or just log the URL

	url := fmt.Sprintf("https://api.foobot.io/v2/device/%v/datapoint/%v/%v/3600/", fooDeviceUUID, startEpoch, currentEpoch)

	foobotClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "go-foobot-api")
	req.Header.Set("X-API-KEY-TOKEN", fooApiKey)

	res, getErr := foobotClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//Unmarshal the JSON response data
	foobot := FoobotJson{}
	jsonErr := json.Unmarshal(body, &foobot)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//*******************REMOVE BELOW DEBUGGING PRINTLN STATEMENTS
	
		log.Println("Called API URL: ", url)
		//log.Println("Current Epoch Time: ", currentEpoch)
		//log.Println("UUID: ", foobot.UUID)
		log.Println("Start time: ", time.Unix(int64(foobot.Start),0))
		log.Println("End time: ", time.Unix(int64(foobot.End),0))
		log.Println("Datapoints Array Length: ", len(foobot.Datapoints))
	

	//Set SQL statement to insert returned API data into database. On conflict (data already exists), do nothing
	sqlStatement := `  
	INSERT INTO sensor_data (device_uuid, seconds_since_epoch, datetime, particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score)  
	VALUES ($1, $2, to_timestamp($3), $4, $5, $6, $7, $8, $9)
	ON CONFLICT DO NOTHING`
	for i := 0; i < len(foobot.Datapoints); i++ {
		_, err = db.Exec(sqlStatement, foobot.UUID, foobot.Datapoints[i][0], int(foobot.Datapoints[i][0]), foobot.Datapoints[i][1], foobot.Datapoints[i][2],
			foobot.Datapoints[i][3], foobot.Datapoints[i][4], foobot.Datapoints[i][5], foobot.Datapoints[i][6])
	}

	if err != nil {
		panic(err)
	}

	//**********************CHANGE BELOW TO LOGGING, REMOVE PRINTLN
	log.Println("Attempted insert of", len(foobot.Datapoints), "records into database.")

}
