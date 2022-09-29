package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DelineaXPM/dsv-sdk-go/v2/vault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	if element := d.Get("element"); element != "" {
		// if element is defined, extract it from the secrets data map and return it
		switch v := secret.Data[element.(string)].(type) {
		case string:
			log.Printf("[DEBUG] returning string %s from .data as the secret", v)
			d.Set("contents", v)
		case map[string]interface{}:
			s, _ := json.Marshal(v)
			log.Printf("[DEBUG] returning JSON %s from .data as the secret", s)
			d.Set("contents", string(s[:]))
		default:
			return fmt.Errorf("element %s is unexpected data type %T", v, v)
		}
	} else {
		// just marshal the whole thing back into JSON and return that
		s, _ := json.Marshal(secret.Data)

		log.Printf("[DEBUG] returning .data as the secret")
		d.Set("contents", string(s[:]))
	}
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
