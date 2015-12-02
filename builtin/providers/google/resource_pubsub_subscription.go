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

			"ackDeadlineSeconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"pushConfig": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     schema.TypeString,
						},

						"pushEndpoint": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

		},
	}
}

func cleanAdditionalArgs(args map[string]interface{}) map[string]string {
	cleaned_args := make(map[string]string)
	for k,v := range  args {
		cleaned_args[k] = v.(string)
	}
	return cleaned_args
}

func resourcePubsubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := fmt.Sprintf("projects/%s/subscriptions/%s", config.Project, d.Get("name").(string))
	computed_topic_name := fmt.Sprintf("projects/%s/topics/%s", config.Project, d.Get("topic").(string))

	//  process optional parameters
	var ackDeadlineSeconds int64
	ackDeadlineSeconds = 10
	if v, ok := d.GetOk("ackDeadlineSeconds"); ok {
		ackDeadlineSeconds = v.(int64)
	}

	var subscription *pubsub.Subscription
	if v, ok := d.GetOk("pushConfig"); ok {
		pushConfigs := v.([]interface{})

		if len(pushConfigs) > 1 {
			return fmt.Errorf("At most one PushConfig is allowed per subscription!")
		}

		pushConfig := pushConfigs[0].(map[string]interface{})
		pushAttributes := pushConfig["attributes"].(map[string]interface{})
		pushAttributesClean := cleanAdditionalArgs(pushAttributes)
		pubsubPushConfig := &pubsub.PushConfig{Attributes: pushAttributesClean, PushEndpoint: pushConfig["pushEndpoint"].(string)}
		subscription = &pubsub.Subscription{AckDeadlineSeconds: ackDeadlineSeconds, Topic: computed_topic_name, PushConfig: pubsubPushConfig}
	}  else {
		subscription = &pubsub.Subscription{AckDeadlineSeconds: ackDeadlineSeconds, Topic: computed_topic_name}
	}

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
