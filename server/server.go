package main

import (

	"github.com/ManbirA/CmpdIntr/controllers"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"io/ioutil"
)

func main() {
	r := httprouter.New()
	tc := controllers.NewTokenController(get_client_id(), get_secret())
	r.GET("/linktoken", tc.Get_link_token)
	r.POST("/accesstoken", tc.Process_access_token)
	r.GET("/transactions/:access_token", tc.Get_transactions)

	http.ListenAndServe("localhost: 8080", r);
}

func get_client_id () string {
	client_id, err := ioutil.ReadFile("./client_id.txt")
	if err != nil {
        fmt.Print(err)
	}
	return string(client_id);
}

func get_secret() string {
	secret, err := ioutil.ReadFile("./secret.txt")
	if err != nil {
        fmt.Print(err)
	}
	return string(secret);
}