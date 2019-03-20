package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/link", GetLinkData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

type StatusError struct {
	Code int
	Err  error
}

type Link struct {
	Link string `json:"link"`
	Host string `json:"host"`
	Scheme string `json:"scheme"`
	Path string `json:"path"`
	ResultData map[string]interface{} `json:"result"`
	BranchKey string `json:"branchKey"`
	error error `json:"error"`

}

func GetLinkData(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)

	var link Link

	_ = json.NewDecoder(r.Body).Decode(&link)
	log.Print(link.Link)
	updateLink(&link)
	getDomain(&link)

	json.NewEncoder(w).Encode(link)
}

func updateLink(link *Link){
	// update link with data from parse
	u, err := url.Parse(link.Link)

	if err != nil {
		link.error = err
		link.ResultData = nil
		return
	}

	link.Path = u.Path
	link.Scheme = u.Scheme
	link.Host = u.Host

	link.BranchKey = branchKeys[u.Host]
}

func getDomain(link *Link) {

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	requestUrl := fmt.Sprintf("https://api2.branch.io/v1/url?url=%s&branch_key=%s", link.Link, link.BranchKey)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil )

	if err != nil {
		link.error = err
		link.ResultData = nil
		return
	}

	res, getErr := client.Do(req)

	if getErr != nil {
		link.error = err
		link.ResultData = nil
		return
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		link.error = err
		link.ResultData = nil
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		link.error = err
		link.ResultData = nil
		return
	}
	link.ResultData = data


}

var branchKeys = map[string]string{
	"jz90.app.link" : "key_live_njOiciGe88AFf5y54I1kklbostfXSDB1",
}