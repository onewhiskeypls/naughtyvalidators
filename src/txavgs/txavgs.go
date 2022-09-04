package txavgs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DTOSummary struct {
	BlockHeight int                `json:"block_height"`
	BlockCount  int                `json:"block_count"`
	TotalTxs    int                `json:"total_txs"`
	TotalTxAvg  float64            `json:"total_tx_avg"`
	Data        []DTOAllValidators `json:"data"`
}

type DTOAllValidators struct {
	Rank            int     `json:"rank"`
	Moniker         string  `json:"moniker"`
	AccountAddress  string  `json:"account_address"`
	OperatorAddress string  `json:"operator_address"`
	Status          int     `json:"status"`
	BlockCount      int     `json:"block_count"`
	TxAvg           float64 `json:"tx_avg"`
}

type DTOValidator struct {
	NumTxs int `json:"num_txs"`
}

var BLOCKCOUNT int = 0
var TOTALTXS int = 0

func GetAllValidators() []DTOAllValidators {
	var dtoAllValidators []DTOAllValidators
	var filteredDtoAllValidators []DTOAllValidators

	url := fmt.Sprintf("https://api-juno.cosmostation.io/v1/staking/validators")

	var allValsData = retrieveDataFromAPI(url)

	if allValsData != nil {
		json.Unmarshal(allValsData, &dtoAllValidators)
	}

	for _, v := range dtoAllValidators {
		if v.Status == 3 {
			filteredDtoAllValidators = append(filteredDtoAllValidators, v)
		}
	}

	return filteredDtoAllValidators
}

func CheckAverageTxs() DTOSummary {
	var summary DTOSummary

	summary.BlockHeight = CheckBlockHeight()

	var allValidators = GetAllValidators()

	for k, v := range allValidators {
		dto := checkAvgTx(v)
		allValidators[k] = dto
	}

	summary.Data = allValidators
	summary.BlockCount = BLOCKCOUNT
	summary.TotalTxs = TOTALTXS
	summary.TotalTxAvg = (float64(TOTALTXS)) / (float64(BLOCKCOUNT))

	return summary
}

func CheckBlockHeight() int {
	var objMap interface{}
	url := "https://api-juno.cosmostation.io/v1/status"
	var data = retrieveDataFromAPI(url)
	if err := json.Unmarshal([]byte(data), &objMap); err != nil {
		log.Fatal(err)
	}

	dataMap := objMap.(map[string]interface{})

	bh := int(dataMap["block_height"].(float64))

	return bh
}

func checkAvgTx(dto DTOAllValidators) DTOAllValidators {
	var blockCnt int = 0
	var numOfTxs int = 0
	var blockData []DTOValidator
	url := fmt.Sprintf("https://api-juno.cosmostation.io/v1/blocks/%s?limit=45&from=0", dto.AccountAddress)

	// TODO: expand to get more than last 45 based on BH
	// https://api-juno.cosmostation.io/v1/blocks/junovaloper1dru5985k4n5q369rxeqfdsjl8ezutch8mc6nx9?limit=45&from=4639284
	var data = retrieveDataFromAPI(url)

	if data != nil {
		json.Unmarshal(data, &blockData)
	}

	for _, v := range blockData {
		blockCnt++
		BLOCKCOUNT++
		numOfTxs += v.NumTxs
		TOTALTXS += v.NumTxs
	}

	dto.TxAvg = (float64(numOfTxs)) / (float64(blockCnt))
	dto.BlockCount = blockCnt

	return dto
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
