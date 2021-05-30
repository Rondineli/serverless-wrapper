package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccResourceProviders map[string]terraform.ResourceProvider
var testAccSchemaProvider *schema.Provider

func init() {
	testAccSchemaProvider = Provider().(*schema.Provider)
	testAccResourceProviders = map[string]terraform.ResourceProvider{
		"wrapper": testAccSchemaProvider,
	}
}

func TestResourceProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

/*
resource "cloudhealth_perspective" "acc_test_owner_tag" {
  name               = "acc test owner tag"
  include_in_reports = false
	hard_delete        = true
  group {
    name = "OwnerAccTest"
    type = "categorize"
    rule {
      asset     = "AwsAsset"
      tag_field = ["owner"]
    }
  }
}
*/

const testAccCreateConfig = `
resource "wrapper" "foo-layer" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./layer/"
  artifact_name = "foo-layer"
  artifact_type = "layer"
}
`

func TestAccCheckCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccSchemaProvider,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wrapper.foo-layer", "runtime", "python3.8"),
				),
			},
		},
	})
}