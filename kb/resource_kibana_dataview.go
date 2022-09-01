// Manage the object in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html
// Supported version:
//  - v7

package kb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kibana "github.com/kaminskip88/go-kibana-rest/v8"
	"github.com/kaminskip88/go-kibana-rest/v8/kbapi"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle kibana data view
func resourceKibanaDataView() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaDataViewCreate,
		Read:   resourceKibanaDataViewRead,
		Update: resourceKibanaDataViewUpdate,
		Delete: resourceKibanaDataViewDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
			},
			"space": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"time_field": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceKibanaDataViewCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	timeField := d.Get("time_field").(string)

	client := meta.(*kibana.Client)

	dv := &kbapi.KibanaDataView{
		Title:         name,
		TimeFieldName: timeField,
	}

	dv, err := client.API.KibanaDataView.Create(dv, false, false)
	if err != nil {
		return err
	}

	d.SetId(dv.ID)

	log.Debugf("Created dataview %s successfully", name)

	return resourceKibanaDataViewRead(d, meta)
}

func resourceKibanaDataViewRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	client := meta.(*kibana.Client)

	data, err := client.API.KibanaDataView.Get(id)
	if err != nil {
		return err
	}

	if data == nil {
		log.Debugf("Dataview with id: %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	d.Set("name", data.Title)
	d.Set("time_field", data.TimeFieldName)

	log.Debugf("Dataview with id: %s found", id)

	return nil
}

func resourceKibanaDataViewUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	timeField := d.Get("time_field").(string)

	client := meta.(*kibana.Client)

	dv := &kbapi.KibanaDataView{
		ID:            id,
		Title:         name,
		TimeFieldName: timeField,
	}

	dv, err := client.API.KibanaDataView.Update(dv, false)
	if err != nil {
		return err
	}

	d.SetId(dv.ID)

	log.Debugf("Dataview updated %s successfully", name)

	return resourceKibanaDataViewRead(d, meta)
}

func resourceKibanaDataViewDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	client := meta.(*kibana.Client)

	err := client.API.KibanaDataView.Delete(id)
	if err != nil {
		return err
	}

	d.SetId("")

	log.Debugf("Dataview id: %s deleted successfully", id)
	return nil

}
