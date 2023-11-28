package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"

)

type Project struct {
	Name string
	Apikey string
}

type Config struct {
	Port string
	Projects []Project
}

type UserInfo struct {
	Credit string
	Sender string
	Priority string
}

type metrics struct {
	credit  *prometheus.GaugeVec
}


func main() {

	config := getConfig()
	fmt.Printf("INFO: loaded config file %s\n", os.Args[1])

	reg := prometheus.NewRegistry()

	for _, project := range config.Projects {
		data := getData(project.Apikey)
		credit, err := strconv.ParseFloat(data.Credit, 64)
		if err != nil {
			log.Fatal(err)
		}

		promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "smsmanager_credit",
			Help: "Amount of credit left on the account",
		}, []string{"project"}).WithLabelValues(project.Name).Add(credit)

	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":" + config.Port, nil))
}


func getData(apikey string) UserInfo {
	resp, err := http.Get("https://http-api.smsmanager.cz/GetUserInfo?apikey=" + apikey)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(string(body), "|")

	userInfo := UserInfo{data[0], data[1], data[2]}

	return userInfo
}


func getConfig() Config {

	if len(os.Args) < 2 {
		log.Fatal("ERROR: missing config file")
	}

	file, _ := ioutil.ReadFile(os.Args[1])

	configData := Config{}

	err := json.Unmarshal([]byte(file), &configData)
	if err != nil {
		log.Fatal(err)
	}

	return configData
}
