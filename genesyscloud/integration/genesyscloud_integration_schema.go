package integration

import (
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Integration",

		CreateContext: gcloud.CreateWithPooledClient(createIntegration),
		ReadContext:   gcloud.ReadWithPooledClient(readIntegration),
		UpdateContext: gcloud.UpdateWithPooledClient(updateIntegration),
		DeleteContext: gcloud.DeleteWithPooledClient(deleteIntegration),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"intended_state": {
				Description:  "Integration state (ENABLED | DISABLED | DELETED).",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DISABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED", "DELETED"}, false),
			},
			"integration_type": {
				Description: "Integration type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config": {
				Description: "Integration config. Each integration type has different schema, use [GET /api/v2/integrations/types/{typeId}/configschemas/{configType}](https://developer.mypurecloud.com/api/rest/v2/integrations/#get-api-v2-integrations-types--typeId--configschemas--configType-) to check schema, then use the correct attribute names for properties.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem:        integrationConfigResource,
			},
		},
	}
}

func IntegrationExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: gcloud.GetAllWithPooledClient(getAllIntegrations),
		RefAttrs: map[string]*resourceExporter.RefAttrSettings{
			"config.credentials.*": {RefType: "genesyscloud_integration_credential"},
		},
		JsonEncodeAttributes: []string{"config.properties", "config.advanced"},
		EncodedRefAttrs: map[*resourceExporter.JsonEncodeRefAttr]*resourceExporter.RefAttrSettings{
			{Attr: "config.properties", NestedAttr: "groups"}: {RefType: "genesyscloud_group"},
		},
	}
}

func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud integration. Select an integration by name",
		ReadContext: gcloud.ReadWithPooledClient(dataSourceIntegrationRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the integration",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func ResourceCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Credential",

		CreateContext: gcloud.CreateWithPooledClient(createCredential),
		ReadContext:   gcloud.ReadWithPooledClient(readCredential),
		UpdateContext: gcloud.UpdateWithPooledClient(updateCredential),
		DeleteContext: gcloud.DeleteWithPooledClient(deleteCredential),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Credential name.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"credential_type_name": {
				Description: "Credential type name. Use [GET /api/v2/integrations/credentials/types](https://developer.genesys.cloud/api/rest/v2/integrations/#get-api-v2-integrations-credentials-types) to see the list of available integration credential types. ",
				Type:        schema.TypeString,
				Required:    true,
			},
			"fields": {
				Description: "Credential fields. Different credential types require different fields. Missing any correct required fields will result API request failure. Use [GET /api/v2/integrations/credentials/types](https://developer.genesys.cloud/api/rest/v2/integrations/#get-api-v2-integrations-credentials-types) to check out the specific credential type schema to find out what fields are required. ",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func CredentialExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: gcloud.GetAllWithPooledClient(getAllCredentials),
		RefAttrs:         map[string]*resourceExporter.RefAttrSettings{}, // No Reference
		UnResolvableAttributes: map[string]*schema.Schema{
			"fields": ResourceCredential().Schema["fields"],
		},
	}
}

func DataSourceIntegrationCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud integration credential. Select an integration credential by name",
		ReadContext: gcloud.ReadWithPooledClient(dataSourceIntegrationCredentialRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the integration credential",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
