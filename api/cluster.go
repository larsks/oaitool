package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
func (client *ApiClient) FindCluster(clusterid string) (*ClusterDetail, error) {
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

func (client *ApiClient) GetCluster(clusterid string) (*ClusterDetail, error) {
	var clusterDetail ClusterDetail

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
