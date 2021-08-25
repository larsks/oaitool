package api

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func PullSecretFromFile(path string) (*PullSecret, error) {
	log.Debugf("reading pull secret from %s", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pullSecret PullSecret
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
