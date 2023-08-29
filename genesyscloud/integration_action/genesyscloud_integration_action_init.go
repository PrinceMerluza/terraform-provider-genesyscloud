package integration_action

import (
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"
)

func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterResource("genesyscloud_integration_action", ResourceIntegrationAction())
	regInstance.RegisterDataSource("genesyscloud_integration_action", DataSourceIntegrationAction())
	regInstance.RegisterExporter("genesyscloud_integration_action", IntegrationActionExporter())
}
