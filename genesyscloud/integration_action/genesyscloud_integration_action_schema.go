package integration_action

import (
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIntegrationAction() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud integration action. Select an integration action by name",
		ReadContext: gcloud.ReadWithPooledClient(DataSourceIntegrationActionRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the integration action",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func ResourceIntegrationAction() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Integration Actions. See this page for detailed information on configuring Actions: https://help.mypurecloud.com/articles/add-configuration-custom-actions-integrations/",

		CreateContext: gcloud.CreateWithPooledClient(createIntegrationAction),
		ReadContext:   gcloud.ReadWithPooledClient(readIntegrationAction),
		UpdateContext: gcloud.UpdateWithPooledClient(updateIntegrationAction),
		DeleteContext: gcloud.DeleteWithPooledClient(deleteIntegrationAction),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Name of the action. Can be up to 256 characters long",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"category": {
				Description:  "Category of action. Can be up to 256 characters long.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"integration_id": {
				Description: "The ID of the integration this action is associated with. Changing the integration_id attribute will cause the existing integration_action to be dropped and recreated with a new ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"secure": {
				Description: "Indication of whether or not the action is designed to accept sensitive data. Changing the secure attribute will cause the existing integration_action to be dropped and recreated with a new ID.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"config_timeout_seconds": {
				Description:  "Optional 1-60 second timeout enforced on the execution or test of this action. This setting is invalid for Custom Authentication Actions.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"contract_input": {
				Description:      "JSON Schema that defines the body of the request that the client (edge/architect/postman) is sending to the service, on the /execute path. Changing the contract_input attribute will cause the existing integration_action to be dropped and recreated with a new ID.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: gcloud.SuppressEquivalentJsonDiffs,
			},
			"contract_output": {
				Description:      "JSON schema that defines the transformed, successful result that will be sent back to the caller. Changing the contract_output attribute will cause the existing integration_action to be dropped and recreated with a new ID.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: gcloud.SuppressEquivalentJsonDiffs,
			},
			"config_request": {
				Description: "Configuration of outbound request.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        actionConfigRequest,
			},
			"config_response": {
				Description: "Configuration of response processing.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        actionConfigResponse,
			},
		},
	}
}

func IntegrationActionExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: gcloud.GetAllWithPooledClient(getAllIntegrationActions),
		RefAttrs: map[string]*resourceExporter.RefAttrSettings{
			"integration_id": {RefType: "genesyscloud_integration"},
		},
		JsonEncodeAttributes: []string{"contract_input", "contract_output"},
	}
}
