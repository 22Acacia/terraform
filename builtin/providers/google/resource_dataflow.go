package google

import (
	"fmt"
	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDataflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataflowCreate,
		Read:   resourceDataflowRead,
		Delete: resourceDataflowDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"jarfile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"staging_bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"jobids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"job_states": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}


func resourceDataflowCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}

	jobids, err := terraformGcloud.CreateDataflow(d.Get("name").(string), d.Get("jarfile").(string), d.Get("class").(string), config.Project, d.Get("staging_bucket").(string))
	if err != nil {
		return err
	}

	d.Set("jobids", jobids)
	d.SetId(d.Get("name").(string))

	err = resourceDataflowRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}

func resourceDataflowRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}


	job_states := make([]string, 0)
	for i := 0; i < d.Get("jobid.#").(int); i++ {
		jobidkey:= fmt.Sprintf("jobid.%d", i)
		job_state, err := terraformGcloud.ReadDataflow(d.Get(jobidkey).(string))
		if err != nil {
			return err
		}
		job_states = append(job_states, job_state)
	}

	d.Set("job_states", job_states)

	return nil
}

func resourceDataflowDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := terraformGcloud.InitGcloud(config.AccountFile)
	if err != nil {
		return err
	}

	err = resourceDataflowRead(d, meta)
	if err != nil {
		return err
	}

	failedCancel := make([]string, 0)
	for i := 0; i < d.Get("jobid.#").(int); i++ {
		jobidkey:= fmt.Sprintf("jobid.%d", i)
		jobstatekey := fmt.Sprintf("jobstate.%d", i)
		failedjob, err := terraformGcloud.CancelDataflow(d.Get(jobidkey).(string), d.Get(jobstatekey).(string))
		if err != nil {
			return err
		}
		if failedjob {
			failedCancel = append(failedCancel, d.Get(jobidkey).(string))
		}
	}

	if len(failedCancel) > 0 {
		return fmt.Errorf("Failed to cancel the following jobs: %v", failedCancel)
	}

	d.SetId("")
	return nil
}
