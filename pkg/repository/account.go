package repository

import (
	"encoding/json"
	"fmt"
	"go-workshop-example/pkg/entity"
	"io/ioutil"
	"net/http"
)

type IAccount interface {
	FetchAccount(id string) (*entity.Account, error)
}

type account struct {
}

func NewAccount() *account {
	return &account{}
}

func (account) FetchAccount(id string) (*entity.Account, error) {
	url := fmt.Sprintf("https://go-workshop-2zcpzmfnyq-de.a.run.app/account/%s", id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var account entity.Account
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
