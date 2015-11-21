package google

import (
	"log"
	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainerReplicaController() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerReplicaControllerCreate,
		Read:   resourceContainerReplicaControllerRead,
		Delete: resourceContainerReplicaControllerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"docker_image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"container_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"external_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"optional_args": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:	  schema.TypeString,
			},

			"external_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

		},
	}
}

func cleanOptionalArgs(optional_args map[string]interface{}) map[string]string {
	cleaned_opts := make(map[string]string)
	for k,v := range  optional_args {
		cleaned_opts[k] = v.(string)
	}
	return cleaned_opts
}

func resourceContainerReplicaControllerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.Credentials)
	if err != nil {
		return err
	}

	err = terraformGcloud.InitKubectl(d.Get("container_name").(string), config.Project, d.Get("zone").(string))
	if err != nil {
		return err
	}

	optional_args := cleanOptionalArgs(d.Get("optional_args").(map[string]interface{}))
	uid, err := terraformGcloud.CreateKubeRC(d.Get("name").(string), d.Get("docker_image").(string), d.Get("external_port").(string), optional_args)
	if err != nil {
		return err
	}

	d.SetId(uid)

	return nil
}

func resourceContainerReplicaControllerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.Credentials)
	if err != nil {
		return err
	}

	err = terraformGcloud.InitKubectl(d.Get("container_name").(string), config.Project, d.Get("zone").(string))
	if err != nil {
		return err
	}

	pod_count, external_ip, err := terraformGcloud.ReadKubeRC(d.Get("name").(string), d.Get("external_port").(string))
	if err != nil {
		return err
	}

	if pod_count == 0 {
		//  something has gone awry, there should always be at least one pod
		log.Printf("There are no pods associated with this Replica Controller.  This is unexpected and probably wrong.  Please investigate")
	}

	if external_ip != "" {
		d.Set("external_ip", external_ip)
	}

	return nil
}

func resourceContainerReplicaControllerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.Credentials)
	if err != nil {
		return err
	}

	err = terraformGcloud.InitKubectl(d.Get("container_name").(string), config.Project, d.Get("zone").(string))
	if err != nil {
		return err
	}

	err = terraformGcloud.DeleteKubeRC(d.Get("name").(string),) 
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
