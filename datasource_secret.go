package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thycotic/dsv-sdk-go/vault"
)

func dataSourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	path := d.Get("path").(string)
	dsv, err := vault.New(meta.(vault.Configuration))

	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}
	log.Printf("[DEBUG] getting secret %s", path)

	secret, err := dsv.Secret(path)

	if err != nil {
		log.Printf("[DEBUG] unable to get secret: %s", err)
		return err
	}

	d.SetId(secret.ID)

	// if element is defined, extract it from the secrets data map and return it
	if element := d.Get("element").(string); element != "" {
		if theElement, ok := secret.Data[element]; ok {
			log.Printf("[DEBUG] returning %s from .data as the secret", element)
			d.Set("contents", theElement)
			return nil
		}
		return fmt.Errorf("element %s not in .data", element)
	}

	data, _ := json.Marshal(secret.Data)

	// just marshal the whole thing back into JSON and return that
	d.Set("contents", data)
	log.Printf("[DEBUG] returning .data as the secret")
	return nil
}

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"contents": {
				Computed:    true,
				Description: "the contents of the secret",
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"element": {
				Description: "the element to extract from the secret",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"id": {
				Computed:    true,
				Description: "the (UUID) identifier of the secret",
				Type:        schema.TypeString,
			},
			"path": {
				Description: "the path of the secret",
				Required:    true,
				Type:        schema.TypeString,
			},
			"version": {
				Computed:    true,
				Description: "the version of the secret",
				Type:        schema.TypeInt,
			},
		},
	}
}
