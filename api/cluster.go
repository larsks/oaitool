package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var supportedFiles = []string{
	"bootstrap.ign",
	"master.ign",
	"metadata.json",
	"worker.ign",
	"kubeadmin-password",
	"kubeconfig",
	"kubeconfig-noingress",
	"install-config.yaml",
	"discovery.ign",
	"custom_manifests.json",
	"custom_manifests.yaml",
}

var supportedNetworkTypes = []string{
	"OpenShiftSDN",
	"OVNKubernetes",
}

var supportedImageTypes = []string{
	"minimal-iso",
	"full-iso",
}

var supportedClusterStatus = []string{
	"insufficient",
	"ready",
	"error",
	"preparing-for-installation",
	"pending-for-input",
	"installing",
	"finalizing",
	"installed",
	"adding-hosts",
	"cancelled",
	"installing-pending-user-action",
}

var supportedHostStatus = []string{
	"discovering",
	"known",
	"disconnected",
	"insufficient",
	"disabled",
	"preparing-for-installation",
	"preparing-successful",
	"pending-for-input",
	"installing",
	"installing-in-progress",
	"installing-pending-user-action",
	"resetting-pending-user-action",
	"installed",
	"error",
	"resetting",
	"added-to-existing-cluster",
	"cancelled",
	"binding",
	"unbinding",
	"known-unbound",
	"disconnected-unbound",
	"insufficient-unbound",
	"disabled-unbound",
	"discovering-unbound",
}

func valInList(value string, allowed_values []string) bool {
	for _, this := range allowed_values {
		if this == value {
			return true
		}
	}

	return false
}

func ValidateHostStatus(status string) bool {
	return valInList(status, supportedHostStatus)
}

func ValidateClusterStatus(status string) bool {
	return valInList(status, supportedClusterStatus)
}

func ValidateNetworkType(networkType string) bool {
	return valInList(networkType, supportedNetworkTypes)
}

func ValidateImageType(imageType string) bool {
	return valInList(imageType, supportedImageTypes)
}

func ValidateDownloadFile(filename string) bool {
	return valInList(filename, supportedFiles)
}

func (client *ApiClient) ListClusters() (ClusterList, error) {
	var clusters ClusterList

	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters", client.ApiUrl),
		nil,
	)
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
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

func (client *ApiClient) GetKubeconfig(clusterid string) ([]byte, error) {
	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters/%s/downloads/kubeconfig", client.ApiUrl, clusterid),
		nil,
	)
	if err != nil {
		return nil, err
	}
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
			"failed to fetch kubeconfig: %s [%d]: %s",
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *ApiClient) GetFile(clusterid, filename string) ([]byte, error) {
	req, err := client.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters/%s/downloads/files?file_name=%s",
			client.ApiUrl, clusterid, filename),
		nil,
	)
	if err != nil {
		return nil, err
	}
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
			"failed to fetch %s: %s [%d]: %s",
			filename,
			http.StatusText(resp.StatusCode), resp.StatusCode, body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *ApiClient) GetPullSecret() (*PullSecret, error) {
	var pullSecret PullSecret
	var accessTokenUrl string = "https://api.openshift.com/api/accounts_mgmt/v1/access_token"

	req, err := client.NewRequest(
		"POST",
		accessTokenUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
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

func (client *ApiClient) PatchCluster(clusterid string, patch JsonObject) (*Cluster, error) {
	patchJson, err := patch.ToJSON()
	log.Debugf("patching cluster %s with: %s", clusterid, string(patchJson))
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/clusters/%s", client.ApiUrl, clusterid),
		bytes.NewBuffer(patchJson),
	)
	if err != nil {
		return nil, err
	}
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
			"failed to patch cluster: %s [%d]: %s",
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

func (cluster *Cluster) ToJSON() ([]byte, error) {
	clusterJson, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}

	return clusterJson, nil
}

func (createParams *ClusterCreateParams) ToJSON() ([]byte, error) {
	createParamsJson, err := json.Marshal(createParams)
	if err != nil {
		return nil, err
	}

	return createParamsJson, nil
}

func (patch *ClusterNetworkPatch) ToJSON() ([]byte, error) {
	patchJson, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchJson, nil
}
