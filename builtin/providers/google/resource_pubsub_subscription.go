package google

import (
	"fmt"
	"google.golang.org/api/pubsub/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePubsubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubSubscriptionCreate,
		Read:   resourcePubsubSubscriptionRead,
		Delete: resourcePubsubSubscriptionDelete,

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

		},
	}
}

func resourcePubsubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := fmt.Sprintf("projects/%s/subscriptions/%s", config.Project, d.Get("name").(string))
	computed_topic_name := fmt.Sprintf("projects/%s/topics/%s", config.Project, d.Get("topic").(string))
	subscription := &pubsub.Subscription{Topic: computed_topic_name}

	call := config.clientPubsub.Projects.Subscriptions.Create(name, subscription)
	res, err := call.Do()
	if err != nil {
		return err
	}
	
	d.SetId(res.Name)

	return nil
}

func resourcePubsubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	
	name := d.Id()
	call := config.clientPubsub.Projects.Subscriptions.Get(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}


func resourcePubsubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Subscriptions.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err 
	}
	
	return nil
}
