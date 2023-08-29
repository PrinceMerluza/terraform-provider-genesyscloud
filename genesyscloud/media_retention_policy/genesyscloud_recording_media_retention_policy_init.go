package media_retention_policy

import (
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"
)

func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterResource("genesyscloud_recording_media_retention_policy", ResourceMediaRetentionPolicy())
	regInstance.RegisterDataSource("genesyscloud_recording_media_retention_policy", DataSourceRecordingMediaRetentionPolicy())
	regInstance.RegisterExporter("genesyscloud_recording_media_retention_policy", MediaRetentionPolicyExporter())
}
