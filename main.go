package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var teamsURL string = os.Getenv("TEAMS_CHANNEL_URL")

func parseResID(id *string) (map[string]string, error) {
	splittedID := strings.Split(*id, "/")
	result := make(map[string]string, len(splittedID))
	if len(splittedID) > 8 {
		for k, v := range splittedID {
			if k%2 == 0 && k != 0 {
				result[splittedID[k-1]] = v
			}
		}
	} else {
		return nil, errors.New("Incorrect Resource ID")
	}
	return result, nil
}

func checkBackPool(lbID *string, minLvl int) {
	backPools := make([]string, 0)
	parsedLB, err := parseResID(lbID)
	if err != nil {
		postTeams("Failed to get information about LB with ID: "+*lbID+"\nError: "+err.Error(), teamsURL)
		return
	}
	client := network.NewLoadBalancersClient(parsedLB["subscriptions"])
	client.Authorizer, _ = auth.NewAuthorizerFromCLI()
	lbPtr, err := client.Get(context.Background(), parsedLB["resourceGroups"], parsedLB["loadBalancers"], "")
	if err != nil {
		postTeams("Failed to get information about LB with ID: "+*lbID+"\nError: "+err.Error(), teamsURL)
	} else {
		poolSize := 0
		for _, v := range *lbPtr.BackendAddressPools {
			backPools = append(backPools, *v.ID)
			if v.BackendIPConfigurations != nil {
				poolSize = len(*v.BackendIPConfigurations)
			} else {
				poolSize = 0
			}
			if poolSize == 0 || poolSize < minLvl {
				postTeams("ID: "+*v.ID+"\nNodes: "+strconv.Itoa(poolSize), teamsURL)
			}
		}
	}
}

type TeamsMessage struct {
	Text string
}

func postTeams(msg, url string) {
	t := TeamsMessage{Text: msg}
	jsonMsg, err := json.Marshal(t)
	if err != nil {
		jsonMsg = []byte(`{"text": "Couldn't marshal the message"}`)
	}
	http.Post(url, "application/json", bytes.NewBuffer(jsonMsg))
}

func main() {
	var thresholdValue int
	if len(os.Args) < 3 {
		postTeams("Incorrect number of arguments", teamsURL)
		return
	}
	if os.Args[2] != "" {
		thresholdValue, _ = strconv.Atoi(os.Args[2])
	} else {
		thresholdValue = 1
	}
	for _, v := range strings.Split(os.Args[1], ",") {
		checkBackPool(&v, thresholdValue)
	}
}
