package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thycotic/dsv-sdk-go/vault"
	"log"
)

func dataSourceClientRead(d *schema.ResourceData, meta interface{}) error {
	clientID := d.Get("client_id").(string)
	dsv, err := vault.New(meta.(vault.Configuration))

	if err != nil {
		log.Printf("[DEBUG] configuration error", err)
		return err
	}

	log.Printf("[DEBUG] getting client %s", clientID)

	client, err := dsv.Client(clientID)

	if err != nil {
		log.Print("[DEBUG] unable to get client", err)
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
