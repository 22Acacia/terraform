package google

import (
	"fmt"
	"google.golang.org/api/pubsub/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTopicCreate,
		Read:   resourceTopicRead,
		Delete: resourceTopicDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

		},
	}
}

func resourceTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := fmt.Sprintf("projects/%s/topics/%s", config.Project, d.Get("name").(string))
	topic := &pubsub.Topic{}

	call := config.clientPubsub.Projects.Topics.Create(name, topic)
	res, err := call.Do()
	if err != nil {
		return err
	}
	
	d.SetId(res.Name)

	return nil
}

func resourceTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	
	name := d.Id()
	call := config.clientPubsub.Projects.Topics.Get(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}


func resourceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Topics.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err 
	}
	
	return nil
}
