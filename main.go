package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	delegators "naughtyvalidators/src"
	nonvoters "naughtyvalidators/src/nonvoters"
	txavgs "naughtyvalidators/src/txavgs"
)

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "nonvoters":
		checkNonVoters()
	case "txavgs":
		checkTxAvgs()
	case "delegators":
		getDelegators()
	default:
		fmt.Println("Need cmd")
	}
}

func checkNonVoters() {
	var DTOALLVOTES []nonvoters.DTOVote
	var SOCIALHANDLES = make(map[string]string)
	var NAUGHTYVALIDATORS []string
	var MESSAGES []string

	proposalID := os.Args[2]
	DTOALLVOTES = nonvoters.GetMintscanData(proposalID)
	SOCIALHANDLES = nonvoters.GetDataFromFile()
	NAUGHTYVALIDATORS = nonvoters.BuildNaughtyList(DTOALLVOTES, SOCIALHANDLES)
	MESSAGES = nonvoters.BuildMessages(NAUGHTYVALIDATORS)

	fmt.Println(fmt.Sprintf("Number of nonvoters: %d", len(NAUGHTYVALIDATORS)))

	for i := 0; i < len(MESSAGES); i++ {
		fmt.Println()
		fmt.Println()
		fmt.Println(strings.Trim(MESSAGES[i], " "))
		fmt.Println()
		fmt.Println()
	}
}

func checkTxAvgs() {
	var VALIDATORTXS = txavgs.CheckAverageTxs()

	prettyJSON, err := json.MarshalIndent(VALIDATORTXS, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))

	/*
		for i := 0; i < len(VALIDATORTXS); i++ {
			fmt.Println()
			fmt.Println()
			fmt.Println(VALIDATORTXS[i])
			fmt.Println()
			fmt.Println()
		}
	*/
}

func getDelegators() {
	valoperAddr := os.Args[2]

	delegators.GetAllDelegatorsOfValidator(valoperAddr)
}
