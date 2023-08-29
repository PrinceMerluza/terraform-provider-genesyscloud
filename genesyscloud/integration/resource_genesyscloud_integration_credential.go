package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	gcloud "terraform-provider-genesyscloud/genesyscloud"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v105/platformclientv2"
)

func getAllCredentials(_ context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(clientConfig)

	for pageNum := 1; ; pageNum++ {
		const pageSize = 100
		credentials, _, err := integrationAPI.GetIntegrationsCredentials(pageNum, pageSize)
		if err != nil {
			return nil, diag.Errorf("Failed to get page of credentials: %v", err)
		}

		if credentials.Entities == nil || len(*credentials.Entities) == 0 {
			break
		}

		for _, cred := range *credentials.Entities {
			if cred.Name != nil { // Credential is possible to have no name
				resources[*cred.Id] = &resourceExporter.ResourceMeta{Name: *cred.Name}
			}
		}
	}

	return resources, nil
}

func createCredential(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	name := d.Get("name").(string)
	cred_type := d.Get("credential_type_name").(string)

	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	createCredential := platformclientv2.Credential{
		Name: &name,
		VarType: &platformclientv2.Credentialtype{
			Name: &cred_type,
		},
		CredentialFields: buildCredentialFields(d),
	}

	credential, _, err := integrationAPI.PostIntegrationsCredentials(createCredential)

	if err != nil {
		return diag.Errorf("Failed to create credential %s : %s", name, err)
	}

	d.SetId(*credential.Id)

	log.Printf("Created credential %s, %s", name, *credential.Id)
	return readCredential(ctx, d, meta)
}

func readCredential(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	log.Printf("Reading credential %s", d.Id())

	return gcloud.WithRetriesForRead(ctx, d, func() *resource.RetryError {
		currentCredential, resp, getErr := integrationAPI.GetIntegrationsCredential(d.Id())
		if getErr != nil {
			if gcloud.IsStatus404(resp) {
				return resource.RetryableError(fmt.Errorf("Failed to read credential %s: %s", d.Id(), getErr))
			}
			return resource.NonRetryableError(fmt.Errorf("Failed to read credential %s: %s", d.Id(), getErr))
		}

		cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceCredential())
		d.Set("name", *currentCredential.Name)
		d.Set("credential_type_name", *currentCredential.VarType.Name)

		log.Printf("Read credential %s %s", d.Id(), *currentCredential.Name)

		return cc.CheckState()
	})
}

func updateCredential(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	cred_type := d.Get("credential_type_name").(string)

	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	if d.HasChanges("name", "credential_type_name", "fields") {

		log.Printf("Updating credential %s", name)

		_, _, putErr := integrationAPI.PutIntegrationsCredential(d.Id(), platformclientv2.Credential{
			Name: &name,
			VarType: &platformclientv2.Credentialtype{
				Name: &cred_type,
			},
			CredentialFields: buildCredentialFields(d),
		})
		if putErr != nil {
			return diag.Errorf("Failed to update credential %s: %s", name, putErr)
		}
	}

	log.Printf("Updated credential %s %s", name, d.Id())
	return readCredential(ctx, d, meta)
}

func deleteCredential(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	integrationAPI := platformclientv2.NewIntegrationsApiWithConfig(sdkConfig)

	_, err := integrationAPI.DeleteIntegrationsCredential(d.Id())
	if err != nil {
		return diag.Errorf("Failed to delete the credential %s: %s", d.Id(), err)
	}

	return gcloud.WithRetries(ctx, 30*time.Second, func() *resource.RetryError {
		_, resp, err := integrationAPI.GetIntegrationsCredential(d.Id())
		if err != nil {
			if gcloud.IsStatus404(resp) {
				// Integration credential deleted
				log.Printf("Deleted Integration credential %s", d.Id())
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting credential action %s: %s", d.Id(), err))
		}
		return resource.RetryableError(fmt.Errorf("Integration credential %s still exists", d.Id()))
	})
}

func buildCredentialFields(d *schema.ResourceData) *map[string]string {
	results := make(map[string]string)
	if fields, ok := d.GetOk("fields"); ok {
		fieldMap := fields.(map[string]interface{})
		for k, v := range fieldMap {
			results[k] = v.(string)
		}
		return &results
	}
	return &results
}
