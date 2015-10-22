package google

import (
	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePubsub() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubCreate,
		Read:   resourcePubsubRead,
		Delete: resourcePubsubDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subscription_count": &schema.Schema{
				Type:		schema.TypeInt,
				Computed:	true,
			},
		},
	}
}

func resourcePubsubCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	fullid, err := terraformGcloud.CreatePubsub(d)
	if err != nil {
		return err
	}

	err = resourcePubsubRead(d, meta)
	if err != nil {
		return err
	}

	d.SetId(fullid)

	return nil
}

func resourcePubsubRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	exists, subcnt, err := terraformGcloud.ReadPubsub(d)
	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
	}

	d.Set("subscription_count", subcnt)
	
	return nil
}


func resourcePubsubDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	err = resourcePubsubRead(d, meta)
	if err != nil {
		return err
	}

	return terraformGcloud.DeletePubsub(d)
}
