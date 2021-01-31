func_grp=$1
lb_ids=$2
teams_url=$3
az group create --location "westeurope" --name $func_grp
arm_out=$(az deployment group create --resource-group $func_grp --template-file template/azuredeploy.json --parameters teamsURL=$teams_url loadBalancersID=$lb_ids alertLevel='1' --query properties.outputs)
spn_id=$(jq -r .principalId.value <<< $arm_out)
func_name=$(jq -r .functionName.value <<< $arm_out)
IFS=',' read -ra ADDR <<< "$lb_ids"
for i in "${ADDR[@]}"; do
    az role assignment create --assignee $spn_id --role "Network Contributor" --scope $i
done
cd code
func azure functionapp publish $func_name --no-build --custom
az functionapp restart --name $func_name --resource-group $func_grp