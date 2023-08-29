package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	gcloud "terraform-provider-genesyscloud/genesyscloud"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v105/platformclientv2"
)

var (
	integrationConfigResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Integration name.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"notes": {
				Description: "Integration notes.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"properties": {
				Description:      "Integration config properties (JSON string).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: gcloud.SuppressEquivalentJsonDiffs,
			},
			"advanced": {
				Description:      "Integration advanced config (JSON string).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: gcloud.SuppressEquivalentJsonDiffs,
			},
			"credentials": {
				Description: "Credentials required for the integration. The required keys are indicated in the credentials property of the Integration Type.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
)

func getAllIntegrations(_ context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(clientConfig)

	for pageNum := 1; ; pageNum++ {
		const pageSize = 100
		integrations, _, err := integrationAPI.GetIntegrations(pageSize, pageNum, "", nil, "", "")
		if err != nil {
			return nil, diag.Errorf("Failed to get page of integrations: %v", err)
		}

		if integrations.Entities == nil || len(*integrations.Entities) == 0 {
			break
		}

		for _, integration := range *integrations.Entities {
			resources[*integration.Id] = &resourceExporter.ResourceMeta{Name: *integration.Name}
		}
	}

	return resources, nil
}

func createIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	intendedState := d.Get("intended_state").(string)
	integrationType := d.Get("integration_type").(string)

	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	createIntegration := platformclientv2.Createintegrationrequest{
		IntegrationType: &platformclientv2.Integrationtype{
			Id: &integrationType,
		},
	}

	integration, _, err := integrationAPI.PostIntegrations(createIntegration)

	if err != nil {
		return diag.Errorf("Failed to create integration : %s", err)
	}

	d.SetId(*integration.Id)

	//Update integration config separately
	diagErr, name := updateIntegrationConfig(d, integrationAPI)
	if diagErr != nil {
		return diagErr
	}

	// Set attributes that can only be modified in a patch
	if d.HasChange(
		"intended_state") {
		log.Printf("Updating additional attributes for integration %s", name)
		const pageSize = 25
		const pageNum = 1
		_, _, patchErr := integrationAPI.PatchIntegration(d.Id(), pageSize, pageNum, "", nil, "", "", platformclientv2.Integration{
			IntendedState: &intendedState,
		})

		if patchErr != nil {
			return diag.Errorf("Failed to update integration %s: %v", name, patchErr)
		}
	}

	log.Printf("Created integration %s %s", name, *integration.Id)
	return readIntegration(ctx, d, meta)
}

func readIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	log.Printf("Reading integration %s", d.Id())

	return gcloud.WithRetriesForRead(ctx, d, func() *resource.RetryError {
		const pageSize = 100
		const pageNum = 1
		currentIntegration, resp, getErr := integrationAPI.GetIntegration(d.Id(), pageSize, pageNum, "", nil, "", "")
		if getErr != nil {
			if gcloud.IsStatus404(resp) {
				return resource.RetryableError(fmt.Errorf("Failed to read integration %s: %s", d.Id(), getErr))
			}
			return resource.NonRetryableError(fmt.Errorf("Failed to read integration %s: %s", d.Id(), getErr))
		}

		d.Set("integration_type", *currentIntegration.IntegrationType.Id)
		if currentIntegration.IntendedState != nil {
			d.Set("intended_state", *currentIntegration.IntendedState)
		} else {
			d.Set("intended_state", nil)
		}

		// Use returned ID to get current config, which contains complete configuration
		integrationConfig, _, err := integrationAPI.GetIntegrationConfigCurrent(*currentIntegration.Id)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Failed to read config of integration %s: %s", d.Id(), getErr))
		}

		d.Set("config", flattenIntegrationConfig(integrationConfig))

		log.Printf("Read integration %s %s", d.Id(), *currentIntegration.Name)

		return nil
	})
}

func updateIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	intendedState := d.Get("intended_state").(string)

	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	diagErr, name := updateIntegrationConfig(d, integrationAPI)
	if diagErr != nil {
		return diagErr
	}

	if d.HasChange("intended_state") {

		log.Printf("Updating integration %s", name)
		const pageSize = 25
		const pageNum = 1
		_, _, patchErr := integrationAPI.PatchIntegration(d.Id(), pageSize, pageNum, "", nil, "", "", platformclientv2.Integration{
			IntendedState: &intendedState,
		})
		if patchErr != nil {
			return diag.Errorf("Failed to update integration %s: %s", name, patchErr)
		}
	}

	log.Printf("Updated integration %s %s", name, d.Id())
	return readIntegration(ctx, d, meta)
}

func deleteIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	_, _, err := integrationAPI.DeleteIntegration(d.Id())
	if err != nil {
		return diag.Errorf("Failed to delete the integration %s: %s", d.Id(), err)
	}

	return gcloud.WithRetries(ctx, 30*time.Second, func() *resource.RetryError {
		const pageSize = 100
		const pageNum = 1
		_, resp, err := integrationAPI.GetIntegration(d.Id(), pageSize, pageNum, "", nil, "", "")
		if err != nil {
			if gcloud.IsStatus404(resp) {
				// Integration deleted
				log.Printf("Deleted Integration %s", d.Id())
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting integration %s: %s", d.Id(), err))
		}
		return resource.RetryableError(fmt.Errorf("Integration %s still exists", d.Id()))
	})
}
