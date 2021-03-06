package api

import (
	"encoding/json"
	"io/ioutil"
)

func PullSecretFromFile(path string) (*PullSecret, error) {
	var pullSecret PullSecret
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &pullSecret); err != nil {
		return nil, err
	}

	return &pullSecret, nil
}

func (pullSecret *PullSecret) ToJSON() ([]byte, error) {
	pullSecretJson, err := json.Marshal(pullSecret)
	if err != nil {
		return nil, err
	}

	return pullSecretJson, nil
}
