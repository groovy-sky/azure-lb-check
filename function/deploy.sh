arm_out=$(az deployment group create --resource-group $func_grp --template-file template/azuredeploy.json --parameters teamsURL=$teams_url loadBalancersID=$lb_ids alertLevel='1' --query properties.outputs)
spn_id=$(jq -r .principalId.value <<< $arm_out)
func_name=$(jq -r .functionName.value <<< $arm_out)
az role assignment create --assignee $spn_id --role "Network Contributor" --scope $lb_ids
az role assignment create --assignee $spn_id --role "Reader" --scope $lb_ids

cd code
az functionapp restart --name $func_name --resource-group $func_grp
func azure functionapp publish $func_name