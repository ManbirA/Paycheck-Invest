package controllers

import(
	plaid "github.com/plaid/plaid-go/plaid"
	"context"
	"time"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"encoding/json"
)

type TokenController struct {
	client *plaid.APIClient
}

type public_token_struct struct {
	Public_token string 
}

func NewTokenController(client_id string, secret string) *TokenController {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", "")
	configuration.AddDefaultHeader("PLAID-SECRET", "")
	configuration.UseEnvironment(plaid.Sandbox)
	client := plaid.NewAPIClient(configuration)
	return &TokenController{client}
}

func (tc TokenController) Get_link_token(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: "USER_ID_FROM_YOUR_DB",
	}
	request := plaid.NewLinkTokenCreateRequest(
	  "CmpdIntr",
	  "en",
	  []plaid.CountryCode{plaid.COUNTRYCODE_US},
	  user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH, plaid.PRODUCTS_TRANSACTIONS})
	request.SetLinkCustomizationName("default")
	request.SetWebhook("https://webhook-uri.com")
	request.SetAccountFilters(plaid.LinkTokenAccountFilters{
	  Depository: &plaid.DepositoryFilter{
		AccountSubtypes: []plaid.AccountSubtype{plaid.ACCOUNTSUBTYPE_CHECKING, plaid.ACCOUNTSUBTYPE_SAVINGS},
	  },
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, _, err := tc.client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		log.Fatal(err)
	}

	token_json, err_token_json := json.Marshal(resp.GetLinkToken());
	if err_token_json != nil {
		log.Fatal(err_token_json)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", token_json)

}

func (tc TokenController) Process_access_token(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	public_token := public_token_struct{};
	json.NewDecoder(r.Body).Decode(&public_token);
	fmt.Println(public_token.Public_token);
	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(public_token.Public_token)
	exchangePublicTokenResp, _, err := tc.client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
	*exchangePublicTokenReq,
	).Execute()

	if err != nil {
		log.Fatal(err)
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	fmt.Println(accessToken)
	
	access_token_json, err_token_json := json.Marshal(accessToken)
	if err_token_json != nil {
		log.Fatal(err_token_json)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", access_token_json)
}

func (tc TokenController) Get_transactions (w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	access_token := p.ByName("access_token")
	// json.NewDecoder(r.Body).Decode(access_token)

	const iso8601TimeFormat = "2006-01-02"
	startDate := time.Now().Add(-365 * 24 * time.Hour).Format(iso8601TimeFormat)
	endDate := time.Now().Format(iso8601TimeFormat)

	request := plaid.NewTransactionsGetRequest(
	access_token,
	startDate,
	endDate,
	)

	options := plaid.TransactionsGetRequestOptions{
	Count:  plaid.PtrInt32(100),
	Offset: plaid.PtrInt32(0),
	}

	request.SetOptions(options)

	transactionsResp, _, err := tc.client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(*request).Execute()
	if err != nil {
		log.Fatal(err)
	}
	 
	test := transactionsResp.GetTransactions()
	for i := 0; i < len(test); i++ {
		fmt.Println(test[i].Name,": ",test[i].Amount)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "%s\n", access_token_json)
}




