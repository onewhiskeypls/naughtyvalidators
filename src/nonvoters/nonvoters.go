package nonvoters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DTOVote struct {
	Voter  string `json:"voter"`
	Answer string `json:"answer"`
	Voted  bool   `json:"voted"`
}

func GetMintscanData(proposalID string) []DTOVote {
	var dtoAllVotes []DTOVote
	urlPath := fmt.Sprintf("https://api.mintscan.io/v1/juno/proposals/%s/votes/validators", proposalID)
	var allValsData = retrieveDataFromAPI(urlPath)

	if allValsData != nil {
		json.Unmarshal(allValsData, &dtoAllVotes)
	}

	return dtoAllVotes
}

func GetDataFromFile() map[string]string {
	var socialHandles = make(map[string]string)
	jsonFile, err := os.Open("socials.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &socialHandles)

	return socialHandles
}

func BuildNaughtyList(allVotes []DTOVote, socialHandles map[string]string) []string {
	var naughtyValidators []string

	for i := 0; i < len(allVotes); i++ {
		vote := &allVotes[i]

		if vote.Answer == "did not vote" {
			if handle, valid := socialHandles[vote.Voter]; valid {
				if len(handle) > 0 {
					naughtyValidators = append(naughtyValidators, handle)
				}
			}
		}
	}

	return naughtyValidators
}

func BuildMessages(naughtyValidators []string) []string {
	var messages []string
	var str string = ""

	for i := 0; i < len(naughtyValidators); i++ {
		handle := *&naughtyValidators[i]

		if (len(str) + len(handle) + 1) > 280 {
			messages = append(messages, str)
			str = ""
		}

		str = str + handle + " "
	}

	return append(messages, str)
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

	req.Header.Set("User-Agent", "asdf")

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

	return body
}
