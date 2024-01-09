package main

import (
	"context"
	"fmt"
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"github.com/spf13/viper"
)

type User struct {
	UserPrincipalName string `json:"userPrincipalName"`
}

type SignIn struct {
	Value []struct {
		CreatedDateTime string `json:"createdDateTime"`
	} `json:"value"`
}

type Group struct {
	DisplayName string `json:"displayName"`
}

// func processUsers(users []User, client *http.Client, token string, csvWriter *csv.Writer, wg *sync.WaitGroup) {
// 	for _, user := range users {
// 		wg.Add(1)
// 		go func(user User) {
// 			defer wg.Done()

// 			encodedUserPrincipalName := url.QueryEscape(user.UserPrincipalName)

// 			baseURL := "https://graph.microsoft.com/v1.0/auditLogs/signIns"
// 			u, _ := url.Parse(baseURL)
// 			q := u.Query()
// 			q.Set("$filter", fmt.Sprintf("userPrincipalName eq '%s'", encodedUserPrincipalName))
// 			q.Set("$orderby", "createdDateTime desc")
// 			q.Set("$top", "1")
// 			u.RawQuery = q.Encode() // Set the query parameters to the URL

// 			signInResponse, err := getResponseWithRetry(client, u.String(), token) // Pass the string representation of the URL
// 			if err != nil {
// 				fmt.Println("Error getting sign in response:", err)
// 				return
// 			}

// 			var signInData SignIn
// 			if err := json.Unmarshal(signInResponse, &signInData); err != nil {
// 				fmt.Println("Error unmarshalling sign in data:", err)
// 				return
// 			}

// 			lastSignInDate := ""
// 			if len(signInData.Value) > 0 {
// 				lastSignInDate = signInData.Value[0].CreatedDateTime[:10]
// 			}

// 			groupMembershipEndpoint := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/memberOf", encodedUserPrincipalName)
// 			groupResponse, err := getResponseWithRetry(client, groupMembershipEndpoint, token)
// 			if err != nil {
// 				fmt.Println("Error getting group response:", err)
// 				return
// 			}

// 			var groupMembership []Group
// 			if err := json.Unmarshal(groupResponse, &groupMembership); err != nil {
// 				fmt.Println("Error unmarshalling group membership:", err)
// 				return
// 			}

// 			groupNames := make([]string, len(groupMembership))
// 			for i, group := range groupMembership {
// 				groupNames[i] = group.DisplayName
// 			}

// 			if err := csvWriter.Write([]string{user.UserPrincipalName, lastSignInDate, fmt.Sprint(groupNames)}); err != nil {
// 				fmt.Println("Error writing to CSV:", err)
// 				return
// 			}
// 		}(user)
// 	}
// }

const envVarPrefix = "matei"

var userReadScopeList = []string{"https://graph.microsoft.com/.default"}

type config struct {
	Azure struct {
		TenantID     string
		ClientID     string
		ClientSecret string
	}
}

type userInformation struct {
	PrincipalName  string
	LastSignInDate string
	Groups         []string
}

type AzureUserQuerier interface {
	ListAll() ([]User, error)
}

type ODataAzureUserQuerier struct {
}

func (q *ODataAzureUserQuerier) ListAll() ([]User, error) {
	return nil, nil
}

func processUser(u models.Userable) userInformation {
	result := userInformation{}
	result.PrincipalName = *u.GetDisplayName()

	sia := u.GetSignInActivity()
	if sia != nil {
		lsidt := sia.GetLastSignInDateTime()
		if lsidt != nil {
			result.LastSignInDate = lsidt.String()
		} else {
			result.LastSignInDate = "N/A"
		}

	}

	mofs := u.GetMemberOf()
	result.Groups = make([]string, len(mofs))
	for i, m := range mofs {
		dni, _ := m.GetBackingStore().Get("displayName")
		dn := dni.(*string)
		result.Groups[i] = *dn
	}

	return result
}

func main() {
	ctx := context.TODO()
	cfg := config{}

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix(envVarPrefix)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}

	v.Unmarshal(&cfg)

	log.Printf("config=%v", cfg)

	cred, err := azidentity.NewClientSecretCredential(cfg.Azure.TenantID, cfg.Azure.ClientID, cfg.Azure.ClientSecret, nil)
	if err != nil {
		log.Fatalf("can not get clientsecretcredential provider %s", err)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, userReadScopeList)
	if err != nil {
		log.Fatalf("can not construct graph client %s", err)
	}

	result, err := client.Users().Get(ctx, &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
			Select: []string{"signInActivity"},
			Expand: []string{"memberOf"},
		},
	})
	if err != nil {
		log.Fatalf("can not get users %s", err)
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.Userable](result,
		client.GetAdapter(),
		models.CreateUserCollectionResponseFromDiscriminatorValue,
	)

	err = pageIterator.Iterate(context.Background(), func(user models.Userable) bool {
		go func(user models.Userable) {
			res := processUser(user)
			log.Println(res)
		}(user)

		return true
	})

	// // Correctly encode the URL
	// filter := url.QueryEscape("accountEnabled eq true")
	// apiEndpoint := fmt.Sprintf("https://graph.microsoft.com/v1.0/users?$filter=%s&$top=999", filter)

	// csvFile, err := os.Create("active_users.csv")
	// if err != nil {
	// 	fmt.Println("Error creating CSV file:", err)
	// 	return
	// }
	// defer csvFile.Close()

	// csvWriter := csv.NewWriter(csvFile)
	// defer csvWriter.Flush()

	// if err := csvWriter.Write([]string{"UserPrincipalName", "LastSignInDate", "GroupMembership"}); err != nil {
	// 	fmt.Println("Error writing header to CSV:", err)
	// 	return
	// }

	// var wg sync.WaitGroup

	// for {
	// 	response, err := getResponseWithRetry(client, apiEndpoint, token)
	// 	if err != nil {
	// 		fmt.Println("Error getting response:", err)
	// 		return
	// 	}

	// 	var data struct {
	// 		Value         []User `json:"value"`
	// 		OdataNextLink string `json:"@odata.nextLink"`
	// 	}
	// 	if err := json.Unmarshal(response, &data); err != nil {
	// 		fmt.Println("Error unmarshalling response data:", err)
	// 		return
	// 	}

	// 	processUsers(data.Value, client, token, csvWriter, &wg)

	// 	if data.OdataNextLink == "" {
	// 		break
	// 	}
	// 	apiEndpoint = data.OdataNextLink
	// }

	// wg.Wait()

	// csvWriter.Flush()
	// csvFile.Close()
}
