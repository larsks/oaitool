package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (host *Host) GetInventory() (*HostInventory, error) {
	var inventory HostInventory

	if err := json.Unmarshal([]byte(host.Inventory), &inventory); err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (client *ApiClient) FindHost(clusterid string, hostid string) (*Host, error) {
	cluster, err := client.FindCluster(clusterid)
	if err != nil {
		return nil, err
	}

	host, err := client.GetHost(cluster.ID, hostid)
	if err == nil {
		return host, nil
	}

	for _, host := range cluster.Hosts {
		if host.RequestedHostname == hostid {
			return &host, nil
		}
	}

	return nil, fmt.Errorf("no host matching %s", hostid)
}

func (client *ApiClient) SetHostnames(clusterid string, hostnames []HostName) error {
	var hnl HostNameList
	for _, hostname := range hostnames {
		hnl.HostNames = append(hnl.HostNames, hostname)
	}

	hnljson, err := json.Marshal(hnl)
	if err != nil {
		return err
	}

	fmt.Println(string(hnljson))

	req, err := client.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/clusters/%s", client.ApiUrl, clusterid),
		bytes.NewReader(hnljson),
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}

	// XXX: We should factor response error handling into a common
	// function.
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return fmt.Errorf(
			"failed to acquire token: %s [%d]: %s",
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	return nil
}

func (client *ApiClient) GetHost(clusterid, hostid string) (*Host, error) {
	var host Host

	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters/%s/hosts/%s", client.ApiUrl, clusterid, hostid),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return nil, fmt.Errorf(
			"failed to get host %s: %s [%d]: %s",
			hostid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &host); err != nil {
		return nil, err
	}

	return &host, nil
}
