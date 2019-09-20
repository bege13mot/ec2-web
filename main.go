package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type resultStruct struct {
	Latest struct {
		BuildDate  string `json:"build_date,string"`
		AmiID      string `json:"ami_id,string"`
		CommitHash string `json:"commit_hash,string"`
	} `json:"latest"`
}

type build struct {
	RuntimeSec string `json:"runtime_seconds"`
	BuildDate  string `json:"build_date"`
	Result     string `json:"result"`
	Output     string `json:"output"`
}

type builds struct {
	Root struct {
		Builds []build
	} `json:"Build base AMI"`
}

func fixJSON(body []byte) string {

	str := string(body)

	r1 := regexp.MustCompile(`,\W*}`)
	r2 := regexp.MustCompile(`{\W*\w{1,4}\s*{`)
	r3 := regexp.MustCompile(`}([^a-zA-z{}])+}$`)

	correctString := r1.ReplaceAllString(str, "}")

	if r2.MatchString(correctString) {
		correctString = r3.ReplaceAllString(r2.ReplaceAllString(correctString, "{"), "}")
	}

	return correctString
}

func processResult(body []byte) (resultStruct, error) {

	storedBuilds := &builds{}
	result := resultStruct{}

	correctString := fixJSON(body)

	err := json.Unmarshal([]byte(correctString), storedBuilds)
	if err != nil {
		log.Println("Can't Unmarshal json: ", correctString)
		return result, err
	}

	max := int64(0)

	for _, v := range storedBuilds.Root.Builds {

		n, err := strconv.ParseInt(v.BuildDate, 10, 64)
		if err != nil {
			log.Printf("Cant cast %s to Int", v.BuildDate)
			return result, err
		}

		if n > max {
			max = n

			s := strings.Split(v.Output, " ")

			if len(s) != 4 {
				log.Println("Output is not correct: ", v.Output)
				return result, fmt.Errorf("Output is not correct")
			}

			result.Latest.BuildDate = v.BuildDate
			result.Latest.AmiID = s[2]
			result.Latest.CommitHash = s[3]

		}

	}

	return result, err

}

func getBuilds(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Can't read body: ", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Can't read body: %v", err)
			return
		}

		result, err := processResult(body)
		if err != nil {
			log.Println("Can't process body: ", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Can't process body: %v", err)
			return
		}

		answer, err := json.Marshal(result)
		if err != nil {
			log.Println("Can't Marshal answer: ", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Can't Marshal answer: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(answer)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only POST is allowed")
	}
}

func main() {
	fmt.Println("Starting server ")

	mux := http.NewServeMux()
	mux.HandleFunc("/builds", getBuilds)

	server := http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
