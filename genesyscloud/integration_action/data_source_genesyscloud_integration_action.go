package integration_action

import (
	"context"
	"fmt"
	"time"

	gcloud "terraform-provider-genesyscloud/genesyscloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v105/platformclientv2"
)

func DataSourceIntegrationActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdkConfig := m.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	actionName := d.Get("name").(string)

	return gcloud.WithRetries(ctx, 15*time.Second, func() *resource.RetryError {
		for pageNum := 1; ; pageNum++ {
			const pageSize = 100
			integrationAction, _, getErr := integrationAPI.GetIntegrationsActions(pageSize, pageNum, "", "", "", "", "", actionName, "", "", "")

			if getErr != nil {
				return resource.NonRetryableError(fmt.Errorf("failed to get page of integration actions: %s", getErr))
			}

			if integrationAction.Entities == nil || len(*integrationAction.Entities) == 0 {
				return resource.RetryableError(fmt.Errorf("no integration actions found with name: %s", actionName))
			}

			for _, action := range *integrationAction.Entities {
				if action.Name != nil && *action.Name == actionName {
					d.SetId(*action.Id)
					return nil
				}
			}

		}
	})
}
