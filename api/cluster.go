package api

import (
	"encoding/json"
	"fmt"
	"io"
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
func (client *ApiClient) FindCluster(name string) (*ClusterDetail, error) {
	detail, err := client.GetCluster(name)
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
		if cluster.Name == name {
			selected = &cluster
			break
		}
	}

	// We didn't find anything
	if selected == nil {
		return nil, nil
	}

	// We found a cluster, let's try to get the cluster detail
	detail, err = client.GetCluster(selected.ID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (client *ApiClient) GetCluster(id string) (*ClusterDetail, error) {
	var clusterDetail ClusterDetail

	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters/%s", client.ApiUrl, id),
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
	if err := json.Unmarshal(body, &clusterDetail); err != nil {
		return nil, err
	}

	return &clusterDetail, nil
}
