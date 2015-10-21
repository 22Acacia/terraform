package google

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"os/exec"
	"io/ioutil"

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

			"project": &schema.Schema{
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


//  return pointer to a file that contains account information
func setAccountFile(contents string) (string, error) {
	if contents != "" {
		var account accountFile
		if err := parseJSON(&account, contents); err == nil {
			//  raw account info, write out to a file
			tmpfile, err := ioutil.TempFile("","")
			if err != nil {
				return "", err
			}
			_, err = tmpfile.WriteString(contents)
			if err != nil {
				return "", err
			}
			tmpfile.Close()
			return tmpfile.Name(), nil
		} else {
			return contents, nil
		}
	}
	return "", nil
}

func cleanupTempAccountFile(rawAccountFile, account_file string) {
	if rawAccountFile != account_file {
		os.Remove(account_file)
	}
}

//  init function will make sure that gcloud cli is installed,
//  authorized and that dataflow commands are available
func init_gcloud(config *Config) error {
	//  check that gcloud is installed
	_, err := exec.LookPath("gcloud")
	if err != nil {
		log.Println("gcloud cli is not installed.  Please install and try again")
		return err
	}

	//  ensure that the found gcloud is authorized
	account_file, err := setAccountFile(config.AccountFile)
	defer cleanupTempAccountFile(config.AccountFile, account_file)
	if err != nil {
		return err
	}
	auth_cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", account_file)
	var stdout, stderr bytes.Buffer
	auth_cmd.Stdout = &stdout
	auth_cmd.Stderr = &stderr
	err = auth_cmd.Run()
	if err != nil {
		log.Println("Dataflow auth failed with error: %q", stdout.String())
		return err 
	}
	
	// verify that datacloud functions are installed
	//  this will need to be updated when they come out of alpha
	datacloud_cmd := exec.Command("gcloud", "alpha", "dataflow" , "-h")
	err = datacloud_cmd.Run()
	if err != nil {
		log.Println("gcloud dataflow commands not installed.")
		return err
	}

	return nil
}

func resourceDataflowCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := init_gcloud(config)
	if err != nil {
		return err
	}

	//  at this point we have verified that our command line jankiness is going to work
	//  get to it
	//  I'm assuming, possibly foolishly, that java is installed on this system
	create_dataflow_cmd := exec.Command("java", "-cp", d.Get("jarfile").(string), d.Get("class").(string), "--project="+d.Get("project").(string), "--stagingLocation="+d.Get("staging_bucket").(string), "--jobName="+d.Get("name").(string))
	var stdout, stderr bytes.Buffer
	create_dataflow_cmd.Stdout = &stdout
	create_dataflow_cmd.Stderr = &stderr
	err = create_dataflow_cmd.Run()
	if err != nil {
		return fmt.Errorf("Error submitting dataflow job: %q", stderr.String())
	}
	
	//  job successfully submitted, now get the job id
	jobidRe := regexp.MustCompile("Submitted job: ([0-9-_]+)\n")
	jobidmatches := jobidRe.FindAllStringSubmatch(stdout.String(), -1)
	jobids := make([]string, 0)
	for _, match := range jobidmatches {
		jobids = append(jobids, match[1])
	}
	d.Set("jobids", jobids)
	d.SetId(d.Get("name").(string))

	err = resourceDataflowRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}

type dataflowDescription struct {
	ClientRequestId		string	`json:"clientRequestId"`
	CreateTime		string	`json:"createTime"`
	CurrentState		string	`json:"currentState"`
	CurrentStateTime	string	`json:"currentStateTime"`
	Id			string	`json:"id"`
	Name			string	`json:"name"`
	ProjectId		string	`json:"projectId"`
	Type			string	`json:"type"`
}

func resourceDataflowRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := init_gcloud(config)
	if err != nil {
		return err
	}

	//  at this point we have verified that our command line jankiness is going to work
	//  get to it
	job_states := make([]string, 0)
	for i := 0; i < d.Get("jobids.#").(int); i++ {
		key := fmt.Sprintf("jobids.%d", i)
		job_check_cmd := exec.Command("gcloud", "alpha", "dataflow", "jobs", "describe", d.Get(key).(string), "--format", "json")
		var stdout, stderr bytes.Buffer
		job_check_cmd.Stdout = &stdout
		job_check_cmd.Stderr = &stderr
		err = job_check_cmd.Run()
		if err != nil {
			return err
		}

		var jobDesc dataflowDescription
		fmt.Println(stdout.String())
		err = parseJSON(&jobDesc, stdout.String())
		if err != nil {
			return err
		}
		job_states = append(job_states, jobDesc.CurrentState)
	}

	d.Set("job_states", job_states)

	return nil
}

func resourceDataflowDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
        err := init_gcloud(config)
	if err != nil {
		return err
	}

	err = resourceDataflowRead(d, meta)
	if err != nil {
		return err
	}

	//  at this point we have verified that our command line jankiness is going to work
	//  get to it
	failedCancel := make([]string, 0)
	for i := 0; i < d.Get("jobids.#").(int); i++ {
		jobstatekey := fmt.Sprintf("job_states.%d", i)
		jobstate := d.Get(jobstatekey).(string)
		if jobstate == "JOB_STATE_RUNNING" {
			jobidkey := fmt.Sprintf("jobids.%d", i)
			job_cancel_cmd := exec.Command("gcloud", "alpha", "dataflow", "jobs", "cancel", d.Get(jobidkey).(string))
			var stdout, stderr bytes.Buffer
			job_cancel_cmd.Stdout = &stdout
			job_cancel_cmd.Stderr = &stderr
			err = job_cancel_cmd.Run()
			if err != nil {
				return err
			}

			if strings.Contains(stdout.String(), "Failed") {
				failedCancel = append(failedCancel,d.Get(jobidkey).(string))
			}
		}
	}

	if len(failedCancel) > 0 {
		errmsg := fmt.Sprintf("Failed to cancel the following jobs: %v", failedCancel)
		return errors.New(errmsg)
	}

	d.SetId("")
	return nil
}
