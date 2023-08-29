package media_retention_policy

import (
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRecordingMediaRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud media retention policy. Select a policy by name",
		ReadContext: gcloud.ReadWithPooledClient(dataSourceRecordingMediaRetentionPolicyRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Media retention policy name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func ResourceMediaRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Genesys Cloud Media Retention Policies",
		CreateContext: gcloud.CreateWithPooledClient(createMediaRetentionPolicy),
		ReadContext:   gcloud.ReadWithPooledClient(readMediaRetentionPolicy),
		UpdateContext: gcloud.UpdateWithPooledClient(updateMediaRetentionPolicy),
		DeleteContext: gcloud.DeleteWithPooledClient(deleteMediaRetentionPolicy),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The policy name. Changing the policy_name attribute will cause the recording_media_retention_policy to be dropped and recreated with a new ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"order": {
				Description: "The ordinal number for the policy",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"description": {
				Description: "The description for the policy",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "The policy will be enabled if true, otherwise it will be disabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"media_policies": {
				Description: "Conditions and actions per media type",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        mediaPolicies,
			},
			"conditions": {
				Description: "Conditions",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        policyConditions,
			},
			"actions": {
				Description: "Actions",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        policyActions,
			},
			"policy_errors": {
				Description: "A list of errors in the policy configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        policyErrors,
			},
		},
	}
}

func MediaRetentionPolicyExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: gcloud.GetAllWithPooledClient(getAllMediaRetentionPolicies),
		RefAttrs: map[string]*resourceExporter.RefAttrSettings{
			"media_policies.chat_policy.conditions.for_queue_ids":                                         {RefType: "genesyscloud_routing_queue", AltValues: []string{"*"}},
			"media_policies.call_policy.conditions.for_queue_ids":                                         {RefType: "genesyscloud_routing_queue", AltValues: []string{"*"}},
			"media_policies.message_policy.conditions.for_queue_ids":                                      {RefType: "genesyscloud_routing_queue", AltValues: []string{"*"}},
			"media_policies.email_policy.conditions.for_queue_ids":                                        {RefType: "genesyscloud_routing_queue", AltValues: []string{"*"}},
			"conditions.for_queue_ids":                                                                    {RefType: "genesyscloud_routing_queue", AltValues: []string{"*"}},
			"media_policies.call_policy.conditions.for_user_ids":                                          {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.chat_policy.conditions.for_user_ids":                                          {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.email_policy.conditions.for_user_ids":                                         {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.message_policy.conditions.for_user_ids":                                       {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"conditions.for_user_ids":                                                                     {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.assign_evaluations.evaluation_form_id":                    {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.call_policy.actions.assign_calibrations.evaluation_form_id":                   {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.call_policy.actions.assign_metered_evaluations.evaluation_form_id":            {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.call_policy.actions.assign_metered_assignment_by_agent.evaluation_form_id":    {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.chat_policy.actions.assign_evaluations.evaluation_form_id":                    {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.chat_policy.actions.assign_calibrations.evaluation_form_id":                   {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.chat_policy.actions.assign_metered_evaluations.evaluation_form_id":            {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.chat_policy.actions.assign_metered_assignment_by_agent.evaluation_form_id":    {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.message_policy.actions.assign_evaluations.evaluation_form_id":                 {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.message_policy.actions.assign_calibrations.evaluation_form_id":                {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.message_policy.actions.assign_metered_evaluations.evaluation_form_id":         {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.message_policy.actions.assign_metered_assignment_by_agent.evaluation_form_id": {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.email_policy.actions.assign_evaluations.evaluation_form_id":                   {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.email_policy.actions.assign_calibrations.evaluation_form_id":                  {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.email_policy.actions.assign_metered_evaluations.evaluation_form_id":           {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.email_policy.actions.assign_metered_assignment_by_agent.evaluation_form_id":   {RefType: "genesyscloud_quality_forms_evaluation"},
			"actions.assign_evaluations.evaluation_form_id":                                               {RefType: "genesyscloud_quality_forms_evaluation"},
			"actions.assign_calibrations.evaluation_form_id":                                              {RefType: "genesyscloud_quality_forms_evaluation"},
			"actions.assign_metered_evaluations.evaluation_form_id":                                       {RefType: "genesyscloud_quality_forms_evaluation"},
			"actions.assign_metered_assignment_by_agent.evaluation_form_id":                               {RefType: "genesyscloud_quality_forms_evaluation"},
			"media_policies.call_policy.actions.assign_evaluations.evaluator_ids":                         {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.assign_calibrations.evaluator_ids":                        {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.assign_metered_evaluations.evaluator_ids":                 {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.assign_metered_assignment_by_agent.evaluator_ids":         {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.chat_policy.actions.assign_evaluations.evaluator_ids":                         {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.chat_policy.actions.assign_calibrations.evaluator_ids":                        {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.chat_policy.actions.assign_metered_evaluations.evaluator_ids":                 {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.chat_policy.actions.assign_metered_assignment_by_agent.evaluator_ids":         {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.message_policy.actions.assign_evaluations.evaluator_ids":                      {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.message_policy.actions.assign_calibrations.evaluator_ids":                     {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.message_policy.actions.assign_metered_evaluations.evaluator_ids":              {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.message_policy.actions.assign_metered_assignment_by_agent.evaluator_ids":      {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.email_policy.actions.assign_evaluations.evaluator_ids":                        {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.email_policy.actions.assign_calibrations.evaluator_ids":                       {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.email_policy.actions.assign_metered_evaluations.evaluator_ids":                {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.email_policy.actions.assign_metered_assignment_by_agent.evaluator_ids":        {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"actions.assign_evaluations.evaluator_ids":                                                    {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"actions.assign_calibrations.evaluator_ids":                                                   {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"actions.assign_metered_evaluations.evaluator_ids":                                            {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"actions.assign_metered_assignment_by_agent.evaluator_ids":                                    {RefType: "genesyscloud_user", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.assign_calibrations.calibrator_id":                        {RefType: "genesyscloud_user"},
			"media_policies.chat_policy.actions.assign_calibrations.calibrator_id":                        {RefType: "genesyscloud_user"},
			"media_policies.message_policy.actions.assign_calibrations.calibrator_id":                     {RefType: "genesyscloud_user"},
			"media_policies.email_policy.actions.assign_calibrations.calibrator_id":                       {RefType: "genesyscloud_user"},
			"media_policies.call_policy.actions.assign_calibrations.expert_evaluator_id":                  {RefType: "genesyscloud_user"},
			"media_policies.chat_policy.actions.assign_calibrations.expert_evaluator_id":                  {RefType: "genesyscloud_user"},
			"media_policies.message_policy.actions.assign_calibrations.expert_evaluator_id":               {RefType: "genesyscloud_user"},
			"media_policies.email_policy.actions.assign_calibrations.expert_evaluator_id":                 {RefType: "genesyscloud_user"},
			"actions.assign_calibrations.expert_evaluator_id":                                             {RefType: "genesyscloud_user"},
			"media_policies.call_policy.conditions.language_ids":                                          {RefType: "genesyscloud_routing_language", AltValues: []string{"*"}},
			"media_policies.chat_policy.conditions.language_ids":                                          {RefType: "genesyscloud_routing_language", AltValues: []string{"*"}},
			"media_policies.message_policy.conditions.language_ids":                                       {RefType: "genesyscloud_routing_language", AltValues: []string{"*"}},
			"media_policies.email_policy.conditions.language_ids":                                         {RefType: "genesyscloud_routing_language", AltValues: []string{"*"}},
			"media_policies.call_policy.conditions.wrapup_code_ids":                                       {RefType: "genesyscloud_routing_wrapupcode", AltValues: []string{"*"}},
			"media_policies.chat_policy.conditions.wrapup_code_ids":                                       {RefType: "genesyscloud_routing_wrapupcode", AltValues: []string{"*"}},
			"media_policies.message_policy.conditions.wrapup_code_ids":                                    {RefType: "genesyscloud_routing_wrapupcode", AltValues: []string{"*"}},
			"media_policies.email_policy.conditions.wrapup_code_ids":                                      {RefType: "genesyscloud_routing_wrapupcode", AltValues: []string{"*"}},
			"conditions.wrapup_code_ids":                                                                  {RefType: "genesyscloud_routing_wrapupcode", AltValues: []string{"*"}},
			"media_policies.call_policy.actions.integration_export.integration_id":                        {RefType: "genesyscloud_integration"},
			"media_policies.chat_policy.actions.integration_export.integration_id":                        {RefType: "genesyscloud_integration"},
			"media_policies.message_policy.actions.integration_export.integration_id":                     {RefType: "genesyscloud_integration"},
			"media_policies.email_policy.actions.integration_export.integration_id":                       {RefType: "genesyscloud_integration"},
			"actions.media_transcriptions.integration_id":                                                 {RefType: "genesyscloud_integration"},
			"media_policies.call_policy.actions.assign_surveys.flow_id":                                   {RefType: "genesyscloud_flow"},
			"media_policies.chat_policy.actions.assign_surveys.flow_id":                                   {RefType: "genesyscloud_flow"},
			"media_policies.message_policy.actions.assign_surveys.flow_id":                                {RefType: "genesyscloud_flow"},
			"media_policies.email_policy.actions.assign_surveys.flow_id":                                  {RefType: "genesyscloud_flow"},
			"actions.assign_surveys.flow_id":                                                              {RefType: "genesyscloud_flow"},
			"media_policies.call_policy.actions.assign_evaluations.user_id":                               {RefType: "genesyscloud_user"},
			"media_policies.chat_policy.actions.assign_evaluations.user_id":                               {RefType: "genesyscloud_user"},
			"media_policies.message_policy.actions.assign_evaluations.user_id":                            {RefType: "genesyscloud_user"},
			"media_policies.email_policy.actions.assign_evaluations.user_id":                              {RefType: "genesyscloud_user"},
			"actions.assign_evaluations.user_id":                                                          {RefType: "genesyscloud_user"},
		},
		AllowZeroValues: []string{"order"},
		RemoveIfMissing: map[string][]string{
			"":               {"conditions", "actions"},
			"media_policies": {"call_policy", "chat_policy", "message_policy", "email_policy"},
		},
	}
}
