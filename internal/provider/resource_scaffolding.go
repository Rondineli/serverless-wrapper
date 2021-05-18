package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWrapper() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,

		Schema: map[string]*schema.Schema{
			"runtime": {
				// This description is used by the documentation generator and the language server.
				Description: "module runtime ie: python3.8",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"build_method": {
				// This description is used by the documentation generator and the language server.
				Description: "Build method to wrapper the function, ie: pip, python3.8, pipenv, makefile",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"requirements_path": {
				// This description is used by the documentation generator and the language server.
				Description: "your requirements, Makefile path",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"artifact_name": {
				// This description is used by the documentation generator and the language server.
				Description: "your requirements, Makefile path",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"use_docker": {
				// This description is used by the documentation generator and the language server.
				Description: "Make zip file on a docker container",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	//if _, err := os.Stat("/tmp/%s.zip", ); os.IsNotExist(err) {
	  // path/to/whatever does not exist
	//}
	fmt.Println(d)


	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	return diag.Errorf("not implemented")
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
