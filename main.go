package main

import (
	"github.com/larsks/oaitool/cli"
	"github.com/spf13/cobra"
)

/*
func main() {
	var response TokenResponse

	offlinetoken, err := ioutil.ReadFile("offlinetoken.txt")
	if err != nil {
		panic(err)
	}

	ssourl := "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"
	listclusterurl := "https://api.openshift.com/api/assisted-install/v1/clusters"

	params := url.Values{}
	params.Add("client_id", "cloud-services")
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", strings.TrimSuffix(string(offlinetoken), "\n"))

	client := &http.Client{}

	resp, err := client.PostForm(ssourl, params)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("failed to acquire token: %s",
			http.StatusText(resp.StatusCode)))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", listclusterurl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", response.Access_token))

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var clusters ClusterList

	fmt.Printf("status: %d\n", resp.StatusCode)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &clusters); err != nil {
		panic(err)
	}

	for _, cluster := range clusters {
		fmt.Printf("%+v\n", cluster)
	}
}
*/

func main() {
	root := cli.NewCmdRoot()
	cobra.CheckErr(root.Execute())
}
