# Monitor Azure Load Balancers Backend Pool using Azure Functions

## Introduction

This repository contains config and code to deploy and configure Azure Functions, which could be used to monitor Azure Load Balancers Backend Pools. Functions code is written on [Go](https://github.com/Azure/azure-sdk-for-go#azure-sdk-for-go) and uses [managed identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) to obtain an information about backend pools state.

## Prerequisites

For initial deployment you'll need:
* Azure subscription
* Linux environment ([WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) also works) with installed [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) and [Azure Functions Core Tools](https://docs.microsoft.com/en-us/azure/azure-functions/functions-run-local?tabs=linux%2Ccsharp%2Cbash), or alternatively, [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/overview)
* [Microsoft Teams](https://docs.microsoft.com/en-us/microsoftteams/teams-overview) (will be used for a notification sending)

## Practical Part

Before publishing the function you'll need to [create a new incoming webhook in Teams](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook) and store its URL.

To start the deployment clone this repository and execute "function/deploy.sh" script with following arguments:
1) Resource group name (which will be created/used for other resource deployment)
2) Load Balancers resource ID (one or many)
3) Teams Webhooks URL

![](/pictures/deploy.png)


## Result

If everything went according to plan, then a functions resource should be deployed and configured:

![](/pictures/resources.png)

Function is triggered every 5 minutes and alerts if a pool is empty:
![](/pictures/alert_example.png)

## Related resources

* https://github.com/MicrosoftDocs/azure-dev-docs/blob/master/articles/go/azure-sdk-authorization.md
* https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/internal/iam/authorizers.go
* https://github.com/Azure-Samples/functions-custom-handlers
* https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-standard-diagnostics
* https://azure.microsoft.com/en-us/blog/introducing-azure-load-balancer-insights-using-azure-monitor-for-networks/