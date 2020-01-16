package main

import (
	"github.com/thycotic/dsv-sdk-go/vault"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	dsv := vault.New(meta.(vault.Configuration))
	name := d.Get("name").(string)

	log.Printf("[DEBUG] getting role %s", name)

	role, err := dsv.Role(name)

	if err != nil {
		log.Print("[DEBUG] unable to get role", err)
		return err
	}

	d.SetId(role.ID)
	d.Set("role_provider", role.Provider)
	d.Set("external_id", role.ExternalID)
	d.Set("groups", role.Groups)
	d.Set("version", role.Version)
	return nil
}

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRoleRead,

		Schema: map[string]*schema.Schema{
			"role_provider": {
				Computed:    true,
				Description: "the provider of the role",
				Type:        schema.TypeString,
			},
			"external_id": {
				Computed:    true,
				Description: "the external-id of the role",
				Type:        schema.TypeString,
			},
			"groups": {
				Computed:    true,
				Description: "the groups associated with the role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeSet,
			},
			"name": {
				Description: "the name of the role",
				Required:    true,
				Type:        schema.TypeString,
			},
			"id": {
				Computed:    true,
				Description: "the (UUID) identifier of the role",
				Type:        schema.TypeString,
			},
			"version": {
				Computed:    true,
				Description: "the version of the role",
				Type:        schema.TypeInt,
			},
		},
	}
}
