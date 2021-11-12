package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/autoscaling"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go/v6"
)

var Version = "development"

type SdkBundle struct {
	CloudApiClient    *ionoscloud.APIClient
	AutoscalingClient *autoscalingService.Client
}

// Provider returns a schema.Provider for ionoscloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(ionoscloud.IonosUsernameEnvVar, nil),
				Description:   "IonosCloud username for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(ionoscloud.IonosPasswordEnvVar, nil),
				Description:   "IonosCloud password for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(ionoscloud.IonosTokenEnvVar, nil),
				Description:   "IonosCloud bearer token for API operations.",
				ConflictsWith: []string{"username", "password"},
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ionoscloud.IonosApiUrlEnvVar, ""),
				Description: "IonosCloud REST API URL.",
			},
			"retries": {
				Type:       schema.TypeInt,
				Optional:   true,
				Default:    50,
				Deprecated: "Timeout is used instead of this functionality",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			DatacenterResource:               resourceDatacenter(),
			IpBLockResource:                  resourceIPBlock(),
			FirewallResource:                 resourceFirewall(),
			LanResource:                      resourceLan(),
			"ionoscloud_loadbalancer":        resourceLoadbalancer(),
			NicResource:                      resourceNic(),
			ServerResource:                   resourceServer(),
			VolumeResource:                   resourceVolume(),
			GroupResource:                    resourceGroup(),
			ShareResource:                    resourceShare(),
			UserResource:                     resourceUser(),
			SnapshotResource:                 resourceSnapshot(),
			ResourceIpFailover:               resourceLanIPFailover(),
			K8sClusterResource:               resourcek8sCluster(),
			K8sNodePoolResource:              resourceK8sNodePool(),
			PCCResource:                      resourcePrivateCrossConnect(),
			BackupUnitResource:               resourceBackupUnit(),
			S3KeyResource:                    resourceS3Key(),
			NatGatewayResource:               resourceNatGateway(),
			NatGatewayRuleResource:           resourceNatGatewayRule(),
			"ionoscloud_networkloadbalancer": resourceNetworkLoadBalancer(),
			"ionoscloud_networkloadbalancer_forwardingrule": resourceNetworkLoadBalancerForwardingRule(),
			"ionoscloud_autoscaling_group":                  resourceAutoscalingGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			DatacenterResource:               dataSourceDataCenter(),
			"ionoscloud_location":            dataSourceLocation(),
			"ionoscloud_image":               dataSourceImage(),
			"ionoscloud_resource":            dataSourceResource(),
			SnapshotResource:                 dataSourceSnapshot(),
			LanResource:                      dataSourceLan(),
			PCCResource:                      dataSourcePcc(),
			ServerResource:                   dataSourceServer(),
			K8sClusterResource:               dataSourceK8sCluster(),
			K8sNodePoolResource:              dataSourceK8sNodePool(),
			NatGatewayResource:               dataSourceNatGateway(),
			NatGatewayRuleResource:           dataSourceNatGatewayRule(),
			"ionoscloud_networkloadbalancer": dataSourceNetworkLoadBalancer(),
			"ionoscloud_networkloadbalancer_forwardingrule": dataSourceNetworkLoadBalancerForwardingRule(),
			"ionoscloud_template":                           dataSourceTemplate(),
			BackupUnitResource:                              dataSourceBackupUnit(),
			FirewallResource:                                dataSourceFirewall(),
			S3KeyResource:                                   dataSourceS3Key(),
			GroupResource:                                   dataSourceGroup(),
			UserResource:                                    dataSourceUser(),
			IpBLockResource:                                 dataSourceIpBlock(),
			VolumeResource:                                  dataSourceVolume(),
			NicResource:                                     dataSourceNIC(),
			ShareResource:                                   dataSourceShare(),
			ResourceIpFailover:                              dataSourceIpFailover(),
			"ionoscloud_autoscaling_group":                  dataSourceAutoscalingGroup(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		log.Printf("[DEBUG] Setting terraformVersion to %s", terraformVersion)

		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {

	username, usernameOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")
	token, tokenOk := d.GetOk("token")

	if !tokenOk {
		if !usernameOk {
			diags := diag.FromErr(fmt.Errorf("neither IonosCloud token, nor IonosCloud username has been provided"))
			return nil, diags
		}

		if !passwordOk {
			diags := diag.FromErr(fmt.Errorf("neither IonosCloud token, nor IonosCloud password has been provided"))
			return nil, diags
		}
	} else {
		if usernameOk || passwordOk {
			diags := diag.FromErr(fmt.Errorf("only provide IonosCloud token OR IonosCloud username/password"))
			return nil, diags
		}
	}

	cleanedUrl := cleanURL(d.Get("endpoint").(string))

	newConfig := ionoscloud.NewConfiguration(username.(string), password.(string), token.(string), cleanedUrl)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfig.Debug = true
	}

	newClient := ionoscloud.NewAPIClient(newConfig)

	// todo: add ionoscloud.Version when added to sdk-go/v6
	//newConfig.UserAgent = fmt.Sprintf("HashiCorp Terraform/%s Terraform Plugin SDK/%s Terraform Provider Ionoscloud/%s Ionoscloud SDK Go/%s", terraformVersion, meta.SDKVersionString(), Version, newClient.Version)
	newConfig.UserAgent = fmt.Sprintf("HashiCorp Terraform/%s Terraform Plugin SDK/%s Terraform Provider Ionoscloud/%s", terraformVersion, meta.SDKVersionString(), Version)

	// configuring autoscaling client
	autoscalingClient := autoscalingService.NewClientService(username.(string), password.(string), token.(string), cleanedUrl)

	return SdkBundle{
		CloudApiClient:    newClient,
		AutoscalingClient: autoscalingClient.Get(),
	}, nil

}

// cleanURL makes sure trailing slash does not corrupt the state
func cleanURL(url string) string {
	length := len(url)
	if length > 1 && url[length-1] == '/' {
		url = url[:length-1]
	}

	return url
}

// getStateChangeConf gets the default configuration for tracking a request progress
func getStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *resource.StateChangeConf {
	stateConf := &resource.StateChangeConf{
		Pending:        resourcePendingStates,
		Target:         resourceTargetStates,
		Refresh:        resourceStateRefreshFunc(meta, location),
		Timeout:        d.Timeout(timeoutType),
		MinTimeout:     5 * time.Second,
		Delay:          0,   // Don't delay the start
		NotFoundChecks: 600, //Setting high number, to support long timeouts
	}

	return stateConf
}

type RequestFailedError struct {
	msg string
}

func (e RequestFailedError) Error() string {
	return e.msg
}

func IsRequestFailed(err error) bool {
	_, ok := err.(RequestFailedError)
	return ok
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(SdkBundle).CloudApiClient

		fmt.Printf("[INFO] Checking PATH %s\n", path)
		if path == "" {
			return nil, "", fmt.Errorf("can not check a state when path is empty")
		}

		request, _, err := client.GetRequestStatus(context.Background(), path)

		if err != nil {
			return nil, "", fmt.Errorf("request failed with following error: %s", err)
		}

		if *request.Metadata.Status == "FAILED" {
			var msg string
			if request.Metadata.Message != nil {
				msg = fmt.Sprintf("Request failed with following error: %s", *request.Metadata.Message)
			} else {
				msg = "Request failed with an unknown error"
			}
			return nil, "", RequestFailedError{msg}
		}

		if *request.Metadata.Status == "DONE" {
			return request, "DONE", nil
		}

		return nil, *request.Metadata.Status, nil
	}
}

// resourcePendingStates defines states of working in progress
var resourcePendingStates = []string{
	"RUNNING",
	"QUEUED",
}

// resourceTargetStates defines states of completion
var resourceTargetStates = []string{
	"DONE",
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
