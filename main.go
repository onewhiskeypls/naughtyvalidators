package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "time"
  "log"
  "io/ioutil"
  "os"
  "strings"
)

var DTOALLVOTES []DTOVote
var SOCIALHANDLES = make(map[string]string)
var NAUGHTYVALIDATORS []string
var MESSAGES []string

func main() {
  proposalID := os.Args[1]
  getMintscanData(proposalID)
  getDataFromFile()
  buildNaughtyList()
  buildMessages()

  for i := 0; i < len(MESSAGES); i++ {
    fmt.Println()
    fmt.Println()
    fmt.Println(strings.Trim(MESSAGES[i], " "))
    fmt.Println()
    fmt.Println()
  }
}

func buildMessages() {
  var str string = ""

  for i := 0; i < len(NAUGHTYVALIDATORS); i++ {
    handle := *&NAUGHTYVALIDATORS[i]

    if (len(str)+len(handle)+1) > 280 {
      MESSAGES = append(MESSAGES, str)
      str = ""
    }

    str = str + handle + " ";
  }

  MESSAGES = append(MESSAGES, str)
}

func buildNaughtyList() {
  for i := 0; i < len(DTOALLVOTES); i++ {
    vote := &DTOALLVOTES[i]

    if (vote.Answer == "did not vote") {
      if handle, valid := SOCIALHANDLES[vote.Voter]; valid {
        if len(handle) > 0 {
          NAUGHTYVALIDATORS = append(NAUGHTYVALIDATORS, handle)
        }
      }
    }
  }
}

func getMintscanData(proposalID string) {
  urlPath := fmt.Sprintf("https://api.mintscan.io/v1/juno/proposals/%s/votes/validators", proposalID)
  var allValsData = retrieveDataFromAPI(urlPath);

  if (allValsData != nil) {
    json.Unmarshal(allValsData, &DTOALLVOTES)
  }
}

func retrieveDataFromAPI(url string) []byte {
  httpClient := http.Client{
		Timeout: time.Second * 10, // TO after 10s
	}

  req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
    return nil
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
    return nil
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

  if err != nil {
		log.Fatal(err)
    return nil
	}

  return body;
}

func getDataFromFile() {
  jsonFile, err := os.Open("socials.json")

  if err != nil {
      fmt.Println(err)
  }

  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), &SOCIALHANDLES)
}

type DTOVote struct {
  Voter string `json:"voter"`
  Answer string `json:"answer"`
  Voted bool `json:"voted"`
}
