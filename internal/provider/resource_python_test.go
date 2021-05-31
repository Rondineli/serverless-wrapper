package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWrapperPython(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"wrapper.foo-layer", "runtime", regexp.MustCompile("^python3.8")),
				),
			},
		},
	})
}

const testAccResourceScaffolding = `
resource "wrapper" "foo-layer" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./layer/"
  artifact_name = "foo-layer"
  artifact_type = "layer"
}
`
