package main

import (
	"github.com/thycotic/dsv-sdk-go/vault"
	"github.com/hashicorp/terraform/helper/schema"
)

func providerConfig(d *schema.ResourceData) (interface{}, error) {
	return vault.Configuration{
		Tenant:       d.Get("tenant").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
	}, nil
}

// Provider is a Terraform DataSource
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"dsv_client": dataSourceClient(),
			"dsv_role":   dataSourceRole(),
			"dsv_secret": dataSourceSecret(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsv_client": resourceClient(),
		},
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DevOps Secrets Vault tenant",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DevOps Secrets Vault client_id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DevOps Secrets Vault client_secret",
			},
		},
		ConfigureFunc: providerConfig,
	}
}
