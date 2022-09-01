// Manage the object in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html
// Supported version:
//  - v7

package kb

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kibana "github.com/kaminskip88/go-kibana-rest/v8"
	"github.com/kaminskip88/go-kibana-rest/v8/kbapi"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle kibana save object
func resourceKibanaObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaObjectCreate,
		Read:   resourceKibanaObjectRead,
		Update: resourceKibanaObjectUpdate,
		Delete: resourceKibanaObjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
			},
			"space": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew: true,
				Default: "default",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// DiffSuppressFunc: suppressEquivalentNDJSON,
			},
			"attributes": {
				Type:     schema.TypeString,
				Required: true,
				// DiffSuppressFunc: suppressEquivalentNDJSON,
			},
			"reference": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceKibanaObjectCreate(d *schema.ResourceData, meta interface{}) error {
	title := d.Get("name").(string)
	objType := d.Get("type").(string)
	dataJSON := d.Get("attributes").(string)
	overwrite := d.Get("overwrite").(bool)
	references := d.Get("reference").(*schema.Set).List()

	client := meta.(*kibana.Client)

	data, err := unmarshalAttrs([]byte(dataJSON))
	if err != nil {
		return err
	}

	// add title to attributes
	data["title"] = title

	refs, err := buildReferences(references)
	if err != nil {
		return err
	}
	so := &kbapi.KibanaSavedObject{
		Type:       objType,
		Attributes: data,
		References: refs,
	}

	so, err = client.API.KibanaSavedObjectV2.Create(so, overwrite)
	if err != nil {
		return err
	}

	d.SetId(so.ID)

	log.Debugf("Saved object created. id: %s", so.ID)

	return resourceKibanaObjectRead(d, meta)
}

func resourceKibanaObjectRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	objType := d.Get("type").(string)

	client := meta.(*kibana.Client)

	so, err := client.API.KibanaSavedObjectV2.Get(id, objType)
	if err != nil {
		return err
	}

	if so == nil {
		log.Debugf("Sved object %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	jsonAttr, err := marshalAttrs(so.Attributes)
	if err != nil {
		return err
	}

	log.Debugf("Saved object %s found", id)

	d.Set("name", so.Attributes["title"])
	d.Set("type", so.Type)
	// d.Set("space", so.Space)
	d.Set("attributes", jsonAttr)

	return nil
}

func resourceKibanaObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	title := d.Get("name").(string)
	objType := d.Get("type").(string)
	dataJSON := d.Get("attributes").(string)
	references := d.Get("reference").(*schema.Set).List()

	client := meta.(*kibana.Client)

	data, err := unmarshalAttrs([]byte(dataJSON))
	if err != nil {
		return err
	}

	// add title to attributes
	data["title"] = title

	refs, err := buildReferences(references)
	if err != nil {
		return err
	}

	so := &kbapi.KibanaSavedObject{
		ID:         id,
		Type:       objType,
		Attributes: data,
		References: refs,
	}

	so, err = client.API.KibanaSavedObjectV2.Update(so)
	if err != nil {
		return err
	}

	d.SetId(so.ID)

	log.Debugf("Saved object %s updated", id)

	return resourceKibanaObjectRead(d, meta)
}

func resourceKibanaObjectDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	objType := d.Get("type").(string)

	client := meta.(*kibana.Client)

	err := client.API.KibanaSavedObjectV2.Delete(id, objType)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}

func marshalAttrs(attrs map[string]interface{}) ([]byte, error) {
	json, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func unmarshalAttrs(str []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func buildReferences(raws []interface{}) ([]kbapi.Reference, error) {
	// if len(raws) == 0 {
	// 	return make([]*kbapi.Reference, 0), nil
	// }

	var result []kbapi.Reference

	for _, i := range raws {
		raw := i.(map[string]interface{})

		var ref kbapi.Reference

		n, ok := raw["name"]
		if !ok {
			return nil, fmt.Errorf("reference \"name\" required")
		}
		ref.Name = n.(string)

		n, ok = raw["type"]
		if !ok {
			return nil, fmt.Errorf("reference \"type\" required")
		}
		ref.Type = n.(string)

		n, ok = raw["id"]
		if !ok {
			return nil, fmt.Errorf("reference \"id\" required")
		}
		ref.ID = n.(string)

		result = append(result, ref)
	}
	return result, nil
}
