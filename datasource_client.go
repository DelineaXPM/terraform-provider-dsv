package main

import (
	"log"

	"github.com/DelineaXPM/dsv-sdk-go/v2/vault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClientRead(d *schema.ResourceData, meta interface{}) error {
	clientID := d.Get("client_id").(string)
	dsv, err := vault.New(meta.(vault.Configuration))

	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}

	log.Printf("[DEBUG] getting client %s", clientID)

	client, err := dsv.Client(clientID)

	if err != nil {
		log.Printf("[DEBUG] unable to get client: %s", err)
		return err
	}

	d.SetId(client.ClientID) // use the ClientID as the (Terraform state) ID
	d.Set("client_id", client.ClientID)
	d.Set("role", client.RoleName)
	return nil
}

func dataSourceClient() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClientRead,

		Schema: map[string]*schema.Schema{
			"role": {
				Computed:    true,
				Description: "the role of the client",
				Type:        schema.TypeString,
			},
			"client_id": {
				Description: "the client_id of the client",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}
