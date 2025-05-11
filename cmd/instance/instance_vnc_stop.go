package instance

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

type DeleteVNCResponse struct {
	Result string `json:"result"`
}

type VNCAlreadyDeletedResponse struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

func DeleteVNCSession(baseURL, region, instanceinstanceID, token string) {
	url := fmt.Sprintf("%s/v2/instances/%s/vnc?region=%s", baseURL, instanceinstanceID, region)
	req, err := http.NewRequest("DELETE", url, nil)
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

	if resp.StatusCode == http.StatusOK {
		var successResp DeleteVNCResponse
		if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
			utility.Error("Error decoding success response:%s", err)
			return
		}
		if successResp.Result == "ok" {
			utility.Info("✅ VNC session successfully deleted.")
		} else {
			utility.Error("Unexpected success response:%s", err)
		}
	} else {
		var errResp VNCAlreadyDeletedResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			utility.Error("Error decoding error response:%s", err)
			return
		}
		if errResp.Code == "database_operation_find" {
			utility.Info("⚠️ No active session found. Already deleted.")
		} else {
			utility.Error("Error: %s - %s\n", errResp.Code, errResp.Reason)
		}

	}

}

func fetchConfigStop(instanceID, regionFlag string) {
	config.ReadConfig()
	baseURL := config.Current.Meta.URL
	region := config.Current.Meta.DefaultRegion
	token := config.DefaultAPIKey()
	if regionFlag != "" {
		region = regionFlag
	}

	DeleteVNCSession(baseURL, region, instanceID, token)
}

var regionFlag string
var stopCmd = &cobra.Command{
	Use:   "stop [param]",
	Short: "Stops an active VNC session before it expires",
	Long:  ` Use this command to stop a VNC session before it ends on its own. It takes one argument, instance id. `,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		param := args[0]
		fetchConfigStop(param, regionFlag)

	},
}

func init() {
	stopCmd.Flags().StringVarP(&regionFlag, "region", "r", "", "Region to use for the request")
}
