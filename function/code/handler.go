package main

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network"
    "github.com/Azure/go-autorest/autorest/azure/auth"
)

type TeamsMessage struct {
    Text string
}

var teamsURL string = os.Getenv("TEAMS_CHANNEL_URL")
var loadBalancersID string = os.Getenv("LB_ID_LIST")

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
    if os.Getenv("FUNCTIONS_EXTENSION_VERSION") != "" {
        msiConfig := auth.NewMSIConfig()
        client.Authorizer, err = msiConfig.Authorizer()
        if err != nil {
            postTeams("Failed to authorize with MSI", teamsURL)
        }
    } else {
        client.Authorizer, err = auth.NewAuthorizerFromCLI()
        if err != nil {
            postTeams("Failed to authorize with CLI", teamsURL)
        }
    }
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
                postTeams("ID: "+*v.ID+"   \nNodes: "+strconv.Itoa(poolSize), teamsURL)
            }
        }
    }
}

func timerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    thresholdValue, _ := strconv.Atoi(os.Getenv("POOL_ALERT_LVL"))
    for _, v := range strings.Split(loadBalancersID, ",") {
        checkBackPool(&v, thresholdValue)
    }
    w.Write([]byte("{}"))
}

func postTeams(msg, url string) {
    t := TeamsMessage{Text: msg}
    jsonMsg, err := json.Marshal(t)
    if err != nil {
        jsonMsg = []byte(`{"text": "Couldn't marshal the message"}`)
    }
    http.Post(teamsURL, "application/json", bytes.NewBuffer(jsonMsg))
}

func main() {
    customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
    if !exists {
        customHandlerPort = "8080"
    }
    mux := http.NewServeMux()
    mux.HandleFunc("/timer", timerHandler)
    log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
