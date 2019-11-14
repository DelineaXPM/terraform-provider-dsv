package main

import (
	"log"

	"github.com/amigus/dsv-sdk-go/vault"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceClientCreate(d *schema.ResourceData, meta interface{}) error {
	dsv := vault.New(meta.(vault.Configuration))
	role := d.Get("role").(string)
	client := new(vault.Client)
	client.RoleName = role
	err := dsv.New(client)

	if err != nil {
		log.Printf("[DEBUG] unable to create client for role %s: %s", role, err)
		return err
	}
	d.SetId(client.ID)
	d.Set("client_id", client.ClientID)
	d.Set("client_secret", client.ClientSecret)
	d.Set("role", client.RoleName)
	return nil
}

func resourceClientDelete(d *schema.ResourceData, meta interface{}) error {
	dsv := vault.New(meta.(vault.Configuration))
	clientID := d.Get("client_id").(string)

	log.Printf("[DEBUG] getting client %s", clientID)

	client, err := dsv.Client(clientID)

	if err != nil {
		log.Printf("[DEBUG] unable to delete client %s: %s", clientID, err)
		return err
	}

	err = client.Delete()

	if err != nil {
		log.Printf("[DEBUG] unable to delete client %s: %s", clientID, err)
		return err
	}

	return nil
}

func resourceClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceClientCreate,
		Delete: resourceClientDelete,
		Read:   dataSourceClientRead,

		Schema: map[string]*schema.Schema{
			"role": {
				Description: "the role of the client",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"client_id": {
				Computed:    true,
				Description: "the client_id of the client",
				Type:        schema.TypeString,
			},
			"client_secret": {
				Computed:    true,
				Description: "the client_secret of the client",
				Type:        schema.TypeString,
			},
		},
	}
}
