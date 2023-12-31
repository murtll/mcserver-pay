package repository

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/render"
	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/util"
)

type ItemRepository struct {
	ApiUrl *url.URL
}

func NewItemRepository(apiUrl string) (*ItemRepository, error) {
	parsed, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	return &ItemRepository{
		ApiUrl: parsed,
	}, nil
}

func (ir *ItemRepository) GetPromo(promo string) (*entities.Promo, error) {
	requestUrl := ir.ApiUrl.JoinPath("check-promo")
	util.SetQueryParam(requestUrl, "promo", promo)
	res, err := http.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("can't get promo info, status '%s' is not acceptable", res.Status)
	}
	fullPromo := &entities.Promo{}
	err = render.DecodeJSON(res.Body, fullPromo)
	if err != nil {
		return nil, err
	}
	fullPromo.Promo = promo
	return fullPromo, nil
}

func (ir *ItemRepository) GetItem(id int) (*entities.Item, error) {
	requestUrl := ir.ApiUrl.JoinPath("item")
	util.SetQueryParam(requestUrl, "id", strconv.Itoa(id))
	res, err := http.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("can't get item, status '%s' is not acceptable", res.Status)
	}
	item := &entities.Item{}
	err = render.DecodeJSON(res.Body, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
