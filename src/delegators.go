package delegators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DTOAllDelegators struct {
	Height     int32          `json:"height"`
	CreatedAt  string         `json:"created_at"`
	TotalCount int32          `json:"total_count"`
	Delegators []DTODelegator `json:"delegators"`
}

type DTODelegator struct {
	DelegatorAddress string `json:"delegator_address"`
	Amount           string `json:"amount"`
}

func GetAllDelegatorsOfValidator(valoperAddr string) {
	var dto *DTOAllDelegators
	var offset int32 = 0
	var totalCount int32 = 0
	var cont bool = true
	var delegatorArr []DTODelegator = []DTODelegator{}
	var height int32 = 0

	for ok := true; ok; ok = ((offset < totalCount || totalCount == 0) && cont) {
		dto = getMintscanData(valoperAddr, offset)

		if dto != nil {
			if dto.TotalCount > 0 && totalCount == 0 {
				totalCount = dto.TotalCount
				height = dto.Height
			}

			delegatorArr = buildDelegatorList(*dto, delegatorArr)
		}

		offset += 60
		if totalCount == 0 {
			cont = false
		}
	}

	filename := fmt.Sprintf("%s_%d.log", valoperAddr, height)

	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for i := 0; i < len(delegatorArr); i++ {
		newStr := fmt.Sprintf("%s,%s\n", delegatorArr[i].DelegatorAddress, delegatorArr[i].Amount)
		if _, err := f.WriteString(newStr); err != nil {
			log.Println(err)
		}
	}
}

func getMintscanData(valoperAddr string, offset int32) *DTOAllDelegators {
	var dtoAllDelegators *DTOAllDelegators
	urlPath := fmt.Sprintf("https://api.mintscan.io/v1/juno/validators/%s/delegators?limit=60&offset=%d", valoperAddr, offset)
	var valData = retrieveDataFromAPI(urlPath)

	if valData != nil {
		json.Unmarshal(valData, &dtoAllDelegators)
	}

	return dtoAllDelegators
}

func buildDelegatorList(dtoAllDelegators DTOAllDelegators, delegatorArr []DTODelegator) []DTODelegator {
	fmt.Println(len(dtoAllDelegators.Delegators))
	for i := 0; i < len(dtoAllDelegators.Delegators); i++ {
		delegatorArr = append(delegatorArr, dtoAllDelegators.Delegators[i])
	}

	return delegatorArr
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
