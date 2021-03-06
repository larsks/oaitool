package api

import "time"

// Most of these structs were generated using https://mholt.github.io/json-to-go/
type (
	JsonObject interface {
		ToJSON() ([]byte, error)
	}

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

	ClusterList []Cluster

	ClusterNetworkPatch struct {
		ApiVip            string `json:"api_vip"`
		IngressVip        string `json:"ingress_vip"`
		VipDhcpAllocation bool   `json:"vip_dhcp_allocation"`
	}

	Cluster struct {
		AmsSubscriptionID          string               `json:"ams_subscription_id"`
		ApiVip                     string               `json:"api_vip"`
		BaseDNSDomain              string               `json:"base_dns_domain"`
		ClusterNetworkCidr         string               `json:"cluster_network_cidr"`
		ClusterNetworkHostPrefix   int                  `json:"cluster_network_host_prefix"`
		ConnectivityMajorityGroups string               `json:"connectivity_majority_groups"`
		ControllerLogsCollectedAt  time.Time            `json:"controller_logs_collected_at"`
		ControllerLogsStartedAt    time.Time            `json:"controller_logs_started_at"`
		CreatedAt                  time.Time            `json:"created_at"`
		EmailDomain                string               `json:"email_domain"`
		EnabledHostCount           int                  `json:"enabled_host_count"`
		FeatureUsage               string               `json:"feature_usage"`
		HighAvailabilityMode       string               `json:"high_availability_mode"`
		HostNetworks               []HostNetworks       `json:"host_networks"`
		Hosts                      []Host               `json:"hosts"`
		Href                       string               `json:"href"`
		Hyperthreading             string               `json:"hyperthreading"`
		ID                         string               `json:"id"`
		ImageInfo                  ImageInfo            `json:"image_info"`
		IngressVip                 string               `json:"ingress_vip"`
		InstallCompletedAt         time.Time            `json:"install_completed_at"`
		InstallStartedAt           time.Time            `json:"install_started_at"`
		Kind                       string               `json:"kind"`
		MachineNetworkCidr         string               `json:"machine_network_cidr"`
		MonitoredOperators         []MonitoredOperators `json:"monitored_operators"`
		Name                       string               `json:"name"`
		NetworkType                string               `json:"network_type"`
		OcpReleaseImage            string               `json:"ocp_release_image"`
		OpenshiftVersion           string               `json:"openshift_version"`
		OrgID                      string               `json:"org_id"`
		Platform                   Platform             `json:"platform"`
		Progress                   Progress             `json:"progress"`
		PullSecretSet              bool                 `json:"pull_secret_set"`
		SchedulableMasters         bool                 `json:"schedulable_masters"`
		ServiceNetworkCidr         string               `json:"service_network_cidr"`
		SshPublicKey               string               `json:"ssh_public_key"`
		Status                     string               `json:"status"`
		StatusInfo                 string               `json:"status_info"`
		StatusUpdatedAt            time.Time            `json:"status_updated_at"`
		TotalHostCount             int                  `json:"total_host_count"`
		UpdatedAt                  time.Time            `json:"updated_at"`
		UserManagedNetworking      bool                 `json:"user_managed_networking"`
		UserName                   string               `json:"user_name"`
		ValidationsInfo            string               `json:"validations_info"`
		VipDhcpAllocation          bool                 `json:"vip_dhcp_allocation"`
	}
	HostNetworks struct {
		Cidr    string   `json:"cidr"`
		HostIds []string `json:"host_ids"`
	}
	Progress struct {
		CurrentStage   string    `json:"current_stage"`
		StageStartedAt time.Time `json:"stage_started_at"`
		StageUpdatedAt time.Time `json:"stage_updated_at"`
	}
	Host struct {
		CheckedInAt           time.Time `json:"checked_in_at"`
		ClusterID             string    `json:"cluster_id"`
		Connectivity          string    `json:"connectivity"`
		CreatedAt             time.Time `json:"created_at"`
		DiscoveryAgentVersion string    `json:"discovery_agent_version"`
		DisksInfo             string    `json:"disks_info"`
		Href                  string    `json:"href"`
		ID                    string    `json:"id"`
		ImagesStatus          string    `json:"images_status"`
		InfraEnvID            string    `json:"infra_env_id"`
		InstallationDiskID    string    `json:"installation_disk_id"`
		InstallationDiskPath  string    `json:"installation_disk_path"`
		InstallerVersion      string    `json:"installer_version"`
		Inventory             string    `json:"inventory"`
		Kind                  string    `json:"kind"`
		LogsCollectedAt       time.Time `json:"logs_collected_at"`
		LogsInfo              string    `json:"logs_info"`
		LogsStartedAt         time.Time `json:"logs_started_at"`
		NtpSources            string    `json:"ntp_sources"`
		HostProgress          Progress  `json:"progress"`
		HostProgressStages    []string  `json:"progress_stages"`
		RequestedHostname     string    `json:"requested_hostname"`
		Role                  string    `json:"role"`
		StageStartedAt        time.Time `json:"stage_started_at"`
		StageUpdatedAt        time.Time `json:"stage_updated_at"`
		Status                string    `json:"status"`
		StatusInfo            string    `json:"status_info"`
		StatusUpdatedAt       time.Time `json:"status_updated_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		UserName              string    `json:"user_name"`
		ValidationsInfo       string    `json:"validations_info"`
		Bootstrap             bool      `json:"bootstrap,omitempty"`
	}
	MonitoredOperators struct {
		ClusterID       string    `json:"cluster_id"`
		Name            string    `json:"name"`
		OperatorType    string    `json:"operator_type"`
		StatusUpdatedAt time.Time `json:"status_updated_at"`
		TimeoutSeconds  int       `json:"timeout_seconds"`
	}
	Vsphere struct {
	}
	Platform struct {
		Type    string  `json:"type"`
		Vsphere Vsphere `json:"vsphere"`
	}
	HostProgress struct {
	}

	HostInventory struct {
		BmcAddress   string       `json:"bmc_address"`
		BmcV6Address string       `json:"bmc_v6address"`
		Boot         Boot         `json:"boot"`
		CPU          CPU          `json:"cpu"`
		Disks        []Disks      `json:"disks"`
		Gpus         []Gpus       `json:"gpus"`
		Hostname     string       `json:"hostname"`
		Interfaces   []Interfaces `json:"interfaces"`
		Memory       Memory       `json:"memory"`
		Routes       []Routes     `json:"routes"`
		SystemVendor SystemVendor `json:"system_vendor"`
		Timestamp    int          `json:"timestamp"`
	}
	Boot struct {
		CurrentBootMode string `json:"current_boot_mode"`
	}
	CPU struct {
		Architecture string   `json:"architecture"`
		Count        int      `json:"count"`
		Flags        []string `json:"flags"`
		Frequency    float32  `json:"frequency"`
		ModelName    string   `json:"model_name"`
	}
	Disks struct {
		Bootable            bool   `json:"bootable,omitempty"`
		ByID                string `json:"by_id,omitempty"`
		ByPath              string `json:"by_path"`
		DriveType           string `json:"drive_type"`
		Hctl                string `json:"hctl"`
		ID                  string `json:"id"`
		Model               string `json:"model"`
		Name                string `json:"name"`
		Path                string `json:"path"`
		Serial              string `json:"serial"`
		SizeBytes           int64  `json:"size_bytes,omitempty"`
		Smart               string `json:"smart"`
		Vendor              string `json:"vendor"`
		Wwn                 string `json:"wwn,omitempty"`
		IsInstallationMedia bool   `json:"is_installation_media,omitempty"`
	}
	Gpus struct {
		Address  string `json:"address"`
		DeviceID string `json:"device_id"`
		Name     string `json:"name"`
		Vendor   string `json:"vendor"`
		VendorID string `json:"vendor_id"`
	}
	Interfaces struct {
		Biosdevname   string        `json:"biosdevname"`
		Flags         []string      `json:"flags"`
		HasCarrier    bool          `json:"has_carrier,omitempty"`
		Ipv4Addresses []string      `json:"ipv4_addresses"`
		Ipv6Addresses []interface{} `json:"ipv6_addresses"`
		MacAddress    string        `json:"mac_address"`
		Mtu           int           `json:"mtu"`
		Name          string        `json:"name"`
		Product       string        `json:"product"`
		SpeedMbps     int           `json:"speed_mbps,omitempty"`
		Vendor        string        `json:"vendor"`
	}
	Memory struct {
		PhysicalBytes int64 `json:"physical_bytes"`
		UsableBytes   int64 `json:"usable_bytes"`
	}
	Routes struct {
		Destination string `json:"destination"`
		Family      int    `json:"family"`
		Gateway     string `json:"gateway,omitempty"`
		Interface   string `json:"interface"`
	}
	SystemVendor struct {
		Manufacturer string `json:"manufacturer"`
		ProductName  string `json:"product_name"`
		SerialNumber string `json:"serial_number"`
	}

	HostNameList struct {
		HostNames []HostName `json:"hosts_names"`
	}

	HostName struct {
		ID       string `json:"id"`
		HostName string `json:"hostname"`
	}

	ClusterInstallParams struct {
		Name                 string `json:"name"`
		OpenshiftVersion     string `json:"openshift_version"`
		PullSecret           string `json:"pull_secret"`
		HighAvailabilityMode string `json:"high_availability_mode"`
		OCPReleaseImage      string `json:"ocp_release_image"`
	}

	ClusterCreateParams struct {
		// Required
		Name             string `json:"name"`
		OpenshiftVersion string `json:"openshift_version"`
		PullSecret       string `json:"pull_secret"`

		// Optional
		HighAvailabilityMode     string `json:"high_availability_mode,omitempty"`
		OcpReleaseImage          string `json:"ocp_release_image,omitempty"`
		BaseDnsDomain            string `json:"base_dns_domain,omitempty"`
		ClusterNetworkCidr       string `json:"cluster_network_cidr,omitempty"`
		ClusterNetworkHostPrefix int    `json:"cluster_network_host_prefix,omitempty"`
		ServiceNetworkCidr       string `json:"service_network_cidr,omitempty"`
		IngressVip               string `json:"ingress_vip,omitempty"`
		SshPublicKey             string `json:"ssh_public_key,omitempty"`
		VipDhcpAllocation        bool   `json:"vip_dhcp_allocation,omitempty"`
		HttpProxy                string `json:"http_proxy,omitempty"`
		HttpsProxy               string `json:"https_proxy,omitempty"`
		NoProxy                  string `json:"no_proxy,omitempty"`
		UserManagedNetworking    bool   `json:"user_managed_networking,omitempty"`
		AdditionalNtpSource      string `json:"additional_ntp_source,omitempty"`
		Hyperthreading           string `json:"hyperthreading,omitempty"`
		NetworkType              string `json:"network_type,omitempty"`
		SchedulableMasters       bool   `json:"schedulable_masters,omitempty"`
	}

	PullSecret struct {
		Auths map[string]PullSecretCredential `json:"auths"`
	}

	PullSecretCredential struct {
		Auth  string `json:"auth"`
		Email string `json:"email"`
	}

	ImageCreateParams struct {
		ImageType    string `json:"image_type"`
		SshPublicKey string `json:"ssh_public_key"`
	}

	ImageInfo struct {
		SshPublicKey        string `json:"ssh_public_key"`
		SizeBytes           int    `json:"size_bytes"`
		DownloadUrl         string `json:"download_url"`
		GeneratorVersion    string `json:"generator_version"`
		CreatedAt           string `json:"created_at"`
		ExpiresAt           string `json:"expires_at"`
		StaticNetworkConfig string `json:"static_network_config"`
		Type                string `json:"type"`
	}
)
