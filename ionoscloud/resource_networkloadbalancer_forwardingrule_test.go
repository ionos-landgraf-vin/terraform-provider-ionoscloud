//go:build all || nlb
// +build all nlb

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const networkLoadBalancerForwardingRuleResource = NetworkLoadBalancerForwardingRuleResource + "." + NetworkLoadBalancerForwardingRuleTestResource
const dataSourceNetworkLoadBalancerForwardingRuleId = DataSource + "." + NetworkLoadBalancerForwardingRuleResource + "." + NetworkLoadBalancerForwardingRuleDataSourceById
const dataSourceNetworkLoadBalancerForwardingRuleName = DataSource + "." + NetworkLoadBalancerForwardingRuleResource + "." + NetworkLoadBalancerForwardingRuleDataSourceByName

func TestAccNetworkLoadBalancerForwardingRuleBasic(t *testing.T) {
	var networkLoadBalancerForwardingRule ionoscloud.NetworkLoadBalancerForwardingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkLoadBalancerForwardingRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkLoadBalancerForwardingRuleExists(networkLoadBalancerForwardingRuleResource, &networkLoadBalancerForwardingRule),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "name", NetworkLoadBalancerForwardingRuleTestResource),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "algorithm", "SOURCE_IP"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "protocol", "TCP"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "listener_ip", "10.12.118.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "listener_port", "8081"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout", "1000"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout", "1200"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout", "1400"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.retries", "3"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.ip", "22.231.2.2"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.port", "8081"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.weight", "123"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check", "true"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval", "1000"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance", "false"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.1.ip", "22.231.2.3"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.ip", "22.231.2.4"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.port", "8081"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.weight", "123"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.health_check.0.check", "true"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.health_check.0.check_interval", "1000"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.health_check.0.maintenance", "false"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "name", networkLoadBalancerForwardingRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "algorithm", networkLoadBalancerForwardingRuleResource, "algorithm"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "protocol", networkLoadBalancerForwardingRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "listener_ip", networkLoadBalancerForwardingRuleResource, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "listener_port", networkLoadBalancerForwardingRuleResource, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.client_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.connect_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.target_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "health_check.0.retries", networkLoadBalancerForwardingRuleResource, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.ip", networkLoadBalancerForwardingRuleResource, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.port", networkLoadBalancerForwardingRuleResource, "targets.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.weight", networkLoadBalancerForwardingRuleResource, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.check", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.check_interval", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleId, "targets.0.health_check.0.maintenance", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "name", networkLoadBalancerForwardingRuleResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "algorithm", networkLoadBalancerForwardingRuleResource, "algorithm"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "protocol", networkLoadBalancerForwardingRuleResource, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "listener_ip", networkLoadBalancerForwardingRuleResource, "listener_ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "listener_port", networkLoadBalancerForwardingRuleResource, "listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.client_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.connect_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.target_timeout", networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "health_check.0.retries", networkLoadBalancerForwardingRuleResource, "health_check.0.retries"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.ip", networkLoadBalancerForwardingRuleResource, "targets.0.ip"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.port", networkLoadBalancerForwardingRuleResource, "targets.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.weight", networkLoadBalancerForwardingRuleResource, "targets.0.weight"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.check", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.check_interval", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval"),
					resource.TestCheckResourceAttrPair(dataSourceNetworkLoadBalancerForwardingRuleName, "targets.0.health_check.0.maintenance", networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance"),
				),
			},
			{
				Config:      testAccDataSourceNetworkLoadBalancerForwardingRuleWrongNameError,
				ExpectError: regexp.MustCompile(`no network load balancer forwarding rule found with the specified criteria: name`),
			},
			{
				Config: testAccCheckNetworkLoadBalancerForwardingRuleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "algorithm", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "listener_ip", "10.12.119.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "listener_port", "8080"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.client_timeout", "1010"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.connect_timeout", "1210"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.target_timeout", "1410"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "health_check.0.retries", "4"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.ip", "22.231.2.2"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.port", "8081"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.weight", "123"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check", "true"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.check_interval", "1000"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.0.health_check.0.maintenance", "false"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.1.ip", "22.231.2.3"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.2.ip", "22.231.2.4"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.ip", "22.231.2.5"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.port", "8081"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.weight", "123"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.health_check.0.check", "true"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.health_check.0.check_interval", "1010"),
					resource.TestCheckResourceAttr(networkLoadBalancerForwardingRuleResource, "targets.3.health_check.0.maintenance", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkLoadBalancerForwardingRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NetworkLoadBalancerForwardingRuleResource {
			continue
		}

		apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured at checking deletion of forwarding rule %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("network loadbalancer forwarding rule still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNetworkLoadBalancerForwardingRuleExists(n string, networkLoadBalancerForwardingRule *ionoscloud.NetworkLoadBalancerForwardingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNetworkLoadBalancerForwardingRuleExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching NetworkLoadBalancerForwardingRule: %s", rs.Primary.ID)
		}
		if *foundNetworkLoadBalancerForwardingRule.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		networkLoadBalancerForwardingRule = &foundNetworkLoadBalancerForwardingRule

		return nil
	}
}

const testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic = testAccCheckNetworkLoadBalancerConfigBasic + `
resource ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleTestResource + ` {
  	datacenter_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.datacenter_id
 	networkloadbalancer_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.id
 	name = "` + NetworkLoadBalancerForwardingRuleTestResource + `"
 	algorithm = "SOURCE_IP"
 	protocol = "TCP"
 	listener_ip = "10.12.118.224"
 	listener_port = "8081"
 	health_check {
		client_timeout = 1000
     	connect_timeout = 1200
     	target_timeout = 1400
     	retries = 3
 	}
   dynamic "targets" {
     for_each = var.IPs
     content {
        ip = targets.value
   		port = "8081"
   		weight = "123"
   		health_check {
     		check = true
     		check_interval = 1000
     		maintenance = false
   		}
 	  }
  }
}
variable IPs{
  type    = list
  default =["22.231.2.2", "22.231.2.3", "22.231.2.4"]
}
`

const testAccCheckNetworkLoadBalancerForwardingRuleConfigUpdate = testAccCheckNetworkLoadBalancerConfigUpdate + `
resource ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleTestResource + ` {
	datacenter_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.datacenter_id
	networkloadbalancer_id = ` + NetworkLoadBalancerResource + `.` + NetworkLoadBalancerTestResource + `.id
	name = "` + UpdatedResources + `"
	algorithm = "ROUND_ROBIN"
	protocol = "HTTP"
	listener_ip = "10.12.119.224"
	listener_port = "8080"
	health_check {
		client_timeout = 1010
		connect_timeout = 1210
		target_timeout = 1410
		retries = 4
 	}
 	targets {
   		ip = "22.231.2.2"
   		port = "8081"
   		weight = "123"
   		health_check {
     		check = true
     		check_interval = 1000
     		maintenance = false
   		}
 	}
 	targets {
   		ip = "22.231.2.3"
   		port = "8081"
   		weight = "123"
   		health_check {
     		check = true
     		check_interval = 1000
     		maintenance = false
   		}
 	}
 	targets {
   		ip = "22.231.2.4"
   		port = "8081"
   		weight = "123"
   		health_check {
     		check = true
     		check_interval = 1000
     		maintenance = false
   		}
 	}
 	targets {
   		ip = "22.231.2.5"
   		port = "8081"
   		weight = "123"
   		health_check {
     		check = true
     		check_interval = 1010
    		 maintenance = true
   		}
 	}
}
`
const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchId = testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic + `
data ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleDataSourceById + ` {
  datacenter_id = ` + networkLoadBalancerForwardingRuleResource + `.datacenter_id
  networkloadbalancer_id  = ` + networkLoadBalancerForwardingRuleResource + `.networkloadbalancer_id
  id			= ` + networkLoadBalancerForwardingRuleResource + `.id
}
`

const testAccDataSourceNetworkLoadBalancerForwardingRuleMatchName = testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic + `
data ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + networkLoadBalancerForwardingRuleResource + `.datacenter_id
  networkloadbalancer_id  = ` + networkLoadBalancerForwardingRuleResource + `.networkloadbalancer_id
 name			= ` + networkLoadBalancerForwardingRuleResource + `.name
}
`

const testAccDataSourceNetworkLoadBalancerForwardingRuleWrongNameError = testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic + `
data ` + NetworkLoadBalancerForwardingRuleResource + ` ` + NetworkLoadBalancerForwardingRuleDataSourceByName + ` {
  datacenter_id = ` + networkLoadBalancerForwardingRuleResource + `.datacenter_id
  networkloadbalancer_id  = ` + networkLoadBalancerForwardingRuleResource + `.networkloadbalancer_id
  name			= "wrong_name"
}
`
