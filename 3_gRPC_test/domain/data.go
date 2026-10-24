package domain

import (
	_ "embed"
	"encoding/json"
)

//go:embed data.json
var DataJson []byte

func GetLocalData() (Products, error) {
	var products Products
	err := json.Unmarshal(DataJson, &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
