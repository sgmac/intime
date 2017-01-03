package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ResponseTime holds the response from https://timezonedb.com
type ResponseTime struct {
	Status      string
	Message     string
	CountryCode string
	CountryName string
	ZoneName    string
	DateTime    string `json:"formatted"`
}

// A Config struct keeps the credentials information
// stored as a JSON document on `configFile`.
type Config struct {
	ApiKey string
}

// Define the URLs and format, this could be defined in a external
// configuration file.
const URLTimezone = "http://api.timezonedb.com/v2/get-time-zone?key="
const TimezoneFormat = "&format=json&by=zone&zone="

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".intime.cfg")

func getConfig() (string, error) {
	c := Config{}
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}
	return c.ApiKey, nil

}

// GetTime provides a formated response for the request time in a given timezone.
func GetTime(zone string) (string, error) {
	apikey, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	var RequestedURL string = URLTimezone + apikey + TimezoneFormat + zone
	respTime := ResponseTime{}

	resp, err := http.Get(RequestedURL)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&respTime)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %s\n", respTime.ZoneName, respTime.DateTime), nil

}

func main() {
	var country string
	flag.StringVar(&country, "c", "", "Country to get time information.")
	flag.Parse()

	if flag.NFlag() < 1 {
		fmt.Fprintf(os.Stderr, "Not country was provided.\n")
		os.Exit(1)
	}

	time, err := GetTime(country)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", time)
}
