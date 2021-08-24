package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ValidateNetworkType(networkType string) error {
	switch networkType {
	case "OpenShiftSDN":
		return nil
	case "OVNKubernetes":
		return nil
	}

	return fmt.Errorf("unknown network type; %s", networkType)
}

func ValidateImageType(imageType string) error {
	switch imageType {
	case "minimal-iso":
		return nil
	case "full-iso":
		return nil
	}

	return fmt.Errorf("unknown image type; %s", imageType)
}

func (client *ApiClient) ListClusters() (ClusterList, error) {
	var clusters ClusterList

	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters", client.ApiUrl),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &clusters); err != nil {
		return nil, err
	}

	return clusters, nil
}

// Find a cluster by id or name. First treat `name` as a cluster id
// and attempt to fetch it directly. If that fails, get a list of
// available clusters and look for the cluster name.
func (client *ApiClient) FindCluster(clusterid string) (*Cluster, error) {
	detail, err := client.GetCluster(clusterid)
	if err == nil {
		return detail, nil
	}

	// It was either a name or a bad cluster id; in any case,
	// we get a list of clusters and then search for matching
	// names.
	clusters, err := client.ListClusters()
	if err != nil {
		return nil, err
	}

	var selected *Cluster
	for _, cluster := range clusters {
		if cluster.Name == clusterid {
			selected = &cluster
			break
		}
	}

	// We didn't find anything
	if selected == nil {
		return nil, fmt.Errorf("no cluster matching %s", clusterid)
	}

	// We found a cluster, let's try to get the cluster detail
	detail, err = client.GetCluster(selected.ID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (client *ApiClient) GetCluster(clusterid string) (*Cluster, error) {
	var clusterDetail Cluster

	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters/%s", client.ApiUrl, clusterid),
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
			"failed to get cluster %s: %s [%d]: %s",
			clusterid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &clusterDetail); err != nil {
		return nil, err
	}

	return &clusterDetail, nil
}

func (client *ApiClient) InstallCluster(clusterid string) error {
	req, err := client.NewRequest(
		"POST",
		fmt.Sprintf("%s/clusters/%s/actions/install", client.ApiUrl, clusterid),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return fmt.Errorf(
			"failed to install cluster %s: %s [%d]: %s",
			clusterid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	return nil

}

func (client *ApiClient) CancelCluster(clusterid string) error {
	req, err := client.NewRequest(
		"POST",
		fmt.Sprintf("%s/clusters/%s/actions/cancel", client.ApiUrl, clusterid),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return fmt.Errorf(
			"failed to install cluster %s: %s [%d]: %s",
			clusterid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	return nil

}

func (client *ApiClient) ResetCluster(clusterid string) error {
	req, err := client.NewRequest(
		"POST",
		fmt.Sprintf("%s/clusters/%s/actions/reset", client.ApiUrl, clusterid),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return fmt.Errorf(
			"failed to install cluster %s: %s [%d]: %s",
			clusterid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	return nil

}

func (client *ApiClient) DeleteCluster(clusterid string) error {
	req, err := client.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/clusters/%s", client.ApiUrl, clusterid),
		nil,
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return fmt.Errorf(
			"failed to install cluster %s: %s [%d]: %s",
			clusterid,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	return nil

}

func (client *ApiClient) GetPullSecret() (*PullSecret, error) {
	var pullSecret PullSecret
	var accessTokenUrl string = "https://api.openshift.com/api/accounts_mgmt/v1/access_token"

	req, err := client.NewRequest(
		"POST",
		accessTokenUrl,
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
			"failed to get pull secret: %s [%d]: %s",
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &pullSecret); err != nil {
		return nil, err
	}

	return &pullSecret, nil
}

func (client *ApiClient) CreateDiscoveryImage(
	clusterid string, imageType string, sshPublicKey string) (*Cluster, error) {
	var cluster Cluster

	if err := ValidateImageType(imageType); err != nil {
		return nil, err
	}

	createParams := ImageCreateParams{
		ImageType:    imageType,
		SshPublicKey: sshPublicKey,
	}
	createParamsJson, err := json.Marshal(createParams)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/clusters/%s/downloads/image",
			client.ApiUrl,
			clusterid),
		bytes.NewBuffer(createParamsJson),
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return nil, fmt.Errorf(
			"failed to create cluster: %s [%d]: %s",
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &cluster); err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (client *ApiClient) CreateCluster(cluster *ClusterCreateParams) (*Cluster, error) {
	createParamsJson, err := cluster.ToJSON()
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		"POST",
		fmt.Sprintf("%s/clusters", client.ApiUrl),
		bytes.NewBuffer(createParamsJson),
	)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("unknown error")
		}
		return nil, fmt.Errorf(
			"failed to create cluster: %s [%d]: %s",
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var detail Cluster
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (createParams *ClusterCreateParams) ToJSON() ([]byte, error) {
	createParamsJson, err := json.Marshal(createParams)
	if err != nil {
		return nil, err
	}

	return createParamsJson, nil
}
