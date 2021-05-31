package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform/resource_provider"
)

var testAccResourceProviders map[string]ResourceProvider
var testAccSchemaProvider *schema.Provider

func init() {
	testAccSchemaProvider = Provider()
	testAccResourceProviders = map[string]ResourceProvider{
		"wrapper": testAccSchemaProvider,
	}
}

func TestResourceProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_impl(t *testing.T) {
	var _ ResourceProvider = Provider()
}
