package google

import (
	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"topic": &schema.Schema{
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

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	fullid, err := terraformGcloud.CreateSubscription(d)
	if err != nil {
		return err
	}

	d.SetId(fullid)

	return nil
}

func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	exists, err := terraformGcloud.ReadSubscription(d)
	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
	}

	return nil
}


func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}
	
	err = resourceSubscriptionRead(d, meta)
	if err != nil {
		return err
	}

	return terraformGcloud.DeleteSubscription(d)
}
