package integration

import (
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"
)

func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterResource("genesyscloud_integration", ResourceIntegration())
	regInstance.RegisterResource("genesyscloud_integration_credential", ResourceCredential())

	regInstance.RegisterDataSource("genesyscloud_integration", DataSourceIntegration())
	regInstance.RegisterDataSource("genesyscloud_integration_credential", DataSourceIntegrationCredential())

	regInstance.RegisterExporter("genesyscloud_integration", IntegrationExporter())
	regInstance.RegisterExporter("genesyscloud_integration_credential", CredentialExporter())
}
