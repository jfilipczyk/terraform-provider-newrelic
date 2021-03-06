package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func dataSourceNewRelicAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicAlertPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_preference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Alert Policies")

	name := d.Get("name").(string)

	params := alerts.ListPoliciesParams{
		Name: name,
	}

	policies, err := client.Alerts.ListPolicies(&params)
	if err != nil {
		return err
	}

	var policy *alerts.Policy

	for _, c := range policies {
		if strings.EqualFold(c.Name, name) {
			policy = &c
			break
		}
	}

	if policy == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert policy", name)
	}

	d.SetId(strconv.Itoa(policy.ID))

	return flattenAlertPolicyDataSource(policy, d)
}
