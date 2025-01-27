package kb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	kibana "github.com/kaminskip88/go-kibana-rest/v8"
	"github.com/pkg/errors"
)

func TestAccKibanaObject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaObject,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaObjectExists("kibana_object.test"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckKibanaObjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		// space := "default"
		objType := "search"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaSavedObjectV2.Get(rs.Primary.ID, objType)
		if err != nil {
			return err
		}
		if data == nil {
			return errors.Errorf("Object %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_object" {
			continue
		}

	}

	return nil
}

var testKibanaObject = `
resource "kibana_object" "test" {
	name 	    = "test-search"
	type        = "search"
	attributes	= jsonencode({
		sort = [["@timestamp","desc"]]
		columns = ["message", "_id"]
	})
}
`
