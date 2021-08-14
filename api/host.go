package api

import "encoding/json"

func (host *Host) GetInventory() (*HostInventory, error) {
	var inventory HostInventory

	if err := json.Unmarshal([]byte(host.Inventory), &inventory); err != nil {
		return nil, err
	}
	return &inventory, nil
}
