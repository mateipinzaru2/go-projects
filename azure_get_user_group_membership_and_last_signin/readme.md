# Azure Graph API Client
This is a multi-threaded client for the Microsoft Graph API using the `go` [Microsoft Graph SDK](https://github.com/microsoftgraph/msgraph-sdk-go). It fetches all users from the API and for each user, it fetches the last sign-in activity and the group membership. It then stores the data in a CSV file.

## Pre-requisites
- Azure App Registration with:
  - `User.Read.All` permission
  - `Group.Read.All` permission
  - `AuditLog.Read.All` permission
- `config.yaml` file with the following:
  ```
  azure:
    tenantID: "<tenantID>"
    clientID: "<clientID>"
    clientSecret: "<clientSecret>"
  ```

## Setup
- `brew install go`
- `cd` to where **main.go** is
- `go run .` to run the program

## Performance
The main loop of the program fetches batches of ***999*** users from the API, for each user, it expands the group membership and fetches the last sign-in activity. It then stores the data in a CSV file.

### Performance Stats
For ***3005*** users **->** ***N/A*** execution time