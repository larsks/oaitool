package api

import "time"

type (
	TokenResponse struct {
		AccessToken       string `json:"access_token"`
		ExpiresIn         int    `json:"expires_in"`
		RefreshExpiresIn  int    `json:"refresh_expires_in"`
		RefreshToken      string `json:"refresh_token"`
		TokenType         string `json:"token_type"`
		IdToken           string `json:"id_token"`
		Not_before_policy int    `json:"not-before-policy"`
		SessionState      string `json:"session_state"`
		Scope             string `json:"scope"`
	}

	Cluster struct {
		APIVip             string    `json:"api_vip"`
		BaseDNSDomain      string    `json:"base_dns_domain"`
		EnabledHostCount   int       `json:"enabled_host_count"`
		ID                 string    `json:"id"`
		IngressVip         string    `json:"ingress_vip"`
		InstallCompletedAt time.Time `json:"install_completed_at"`
		InstallStartedAt   time.Time `json:"install_started_at"`
		Name               string    `json:"name"`
		OcpReleaseImage    string    `json:"ocp_release_image"`
		OpenshiftClusterID string    `json:"openshift_cluster_id"`
		OpenshiftVersion   string    `json:"openshift_version"`
		SSHPublicKey       string    `json:"ssh_public_key"`
		Status             string    `json:"status"`
		StatusInfo         string    `json:"status_info"`
		StatusUpdatedAt    time.Time `json:"status_updated_at"`
		TotalHostCount     int       `json:"total_host_count"`
		UpdatedAt          time.Time `json:"updated_at"`
		VipDhcpAllocation  bool      `json:"vip_dhcp_allocation"`
	}

	ClusterList []Cluster
)
