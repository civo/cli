package instance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

type ActiveVNCResponse struct {
	URI        string `json:"uri"`
	Expiration string `json:"expiration"`
}

type InactiveVNCResponse struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

func GetVNCSession(baseURL, region, token, instanceinstanceID string) {
	url := fmt.Sprintf("%s/v2/instances/%s/vnc?region=%s", baseURL, instanceinstanceID, region)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utility.Error("Request creation error:%s", err)
		return
	}

	req.Header.Set("Authorization", "bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		utility.Error("Request failed:%s", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.Error("Read body error:%s", err)
		return
	}

	var vncResp ActiveVNCResponse
	if err := json.Unmarshal(body, &vncResp); err == nil && vncResp.URI != "" {
		utility.Info("🟢 VNC Session Active")
		return

	}
	var errResp InactiveVNCResponse

	if err := json.Unmarshal(body, &errResp); err == nil {
		utility.Info("🔴 VNC Session Not Active")
		return
	}

	fmt.Println("Unexpected response:", string(body))
	utility.Error("Unexpected response:%s", string(body))

}

func fetchConfigStatus(instanceID, regionFlag string) {
	config.ReadConfig()
	baseURL := config.Current.Meta.URL
	region := config.Current.Meta.DefaultRegion
	token := config.DefaultAPIKey()
	if regionFlag != "" {
		region = regionFlag
	}
	GetVNCSession(baseURL, region, token, instanceID)
}

// Status command
var statusCmd = &cobra.Command{
	Use:   "status [param]",
	Short: "Check if a VNC session is active",
	Long:  `Use this command to check if a VNC session is currently running for a given instance. It takes one argument, instance id.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		param := args[0]
		fetchConfigStatus(param, regionFlag)

	},
}

func init() {
	statusCmd.Flags().StringVarP(&regionFlag, "region", "r", "", "Region to use for the request")
}
