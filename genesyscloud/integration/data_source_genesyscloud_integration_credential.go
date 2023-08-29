package integration

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

func dataSourceIntegrationCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdkConfig := m.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	credName := d.Get("name").(string)

	return gcloud.WithRetries(ctx, 15*time.Second, func() *resource.RetryError {
		for pageNum := 1; ; pageNum++ {
			const pageSize = 100
			integrationCredentials, _, getErr := integrationAPI.GetIntegrationsCredentials(pageNum, pageSize)

			if getErr != nil {
				return resource.NonRetryableError(fmt.Errorf("failed to get page of integration credentials: %s", getErr))
			}

			if integrationCredentials.Entities == nil || len(*integrationCredentials.Entities) == 0 {
				return resource.RetryableError(fmt.Errorf("no integration credentials found with name: %s", credName))
			}

			for _, credential := range *integrationCredentials.Entities {
				if credential.Name != nil && *credential.Name == credName {
					d.SetId(*credential.Id)
					return nil
				}
			}

		}
	})
}
