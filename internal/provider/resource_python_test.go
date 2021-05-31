package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePythonWrapperResource(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
  requirements_path = "../../examples/serverless-wrapper/src-function/"
  artifact_name = "foo-layer"
  artifact_type = "layer"
}
`
