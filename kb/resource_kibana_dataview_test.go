package kb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	kibana "github.com/kaminskip88/go-kibana-rest/v8"
	"github.com/pkg/errors"
)

func TestAccKibanaDataview(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaDataviewDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaDataview,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaDataviewExists("kibana_dataview.test"),
				),
			},
		},
	})
}

func testCheckKibanaDataviewExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No dataview ID is set")
		}

		// space := "default"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaDataView.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if data == nil {
			return errors.Errorf("Dataview %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaDataviewDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_dataview" {
			continue
		}

	}

	return nil
}

var testKibanaDataview = `
resource "kibana_dataview" "test" {
	name = "test-logs-*"
	time_field = "@timestamp"
}
`
