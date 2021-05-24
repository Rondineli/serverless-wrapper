package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
  return &schema.Provider{
    ResourcesMap: map[string]*schema.Resource{
		"wrapper_layer": resourceWrapper(),
		"wrapper_function": resourceWrapper(),
		"wrapper_docker": resourceWrapper(),
	},
  }
}
