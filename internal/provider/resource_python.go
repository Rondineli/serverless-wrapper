package provider

import (
	"context"
	"fmt"
	"os/exec"
	"os"
	"bytes"
	"regexp"
	"path"
	"text/template"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const shInstallDepsTemplate = `
	cd "{{.ArtifactDir}}
	find * -name *.zip -exec rm -rf {} \;
	find * -name *.pyc -exec rm -rf {} \;
	find * -name __pycache__ -exec rm -rf {} \;
	if [ "{{.ArtifactType}}" == "layer" ]; then
		mkdir -p "./python/lib/{{.PythonRuntime}}/site-packages/" || echo "ok"
		cp *.py "./python/lib/{{.PythonRuntime}}/site-packages/" || echo "ok"
		export DIR_ARTIFACT="./python/lib/{{.PythonRuntime}}/site-packages"
	else
		export DIR_ARTIFACT="./"
	fi
	{{.PythonRuntime}} -m {{.ExecutionRuntime}} install -r {{.ArtifactDir}}/requirements.txt -t "${DIR_ARTIFACT}"
`

const dockerInstallDepsTemplate = `
	touch /tmp/debug.txt
	echo $(pwd) > /tmp/debug.txt
	docker run --rm -v $(pwd)/{{.ArtifactDir}}:/{{.ArtifactDir}}/ -t "python:3.8-alpine" \
		sh -c "apk add bash --update && bash -c '{{.InstallTemplate}}'"
`

func resourceWrapper() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceCreate,
		DeleteContext: resourceDelete,

		Schema: map[string]*schema.Schema{
			"runtime": {
				// This description is used by the documentation generator and the language server.
				Description: "module runtime ie: python3.8",
				Type:         schema.TypeString,
				Required:     true,
			},
			"build_method": {
				// This description is used by the documentation generator and the language server.
				Description: "Build method to wrapper the function, ie: pip, python3.8, pipenv, make",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "pip",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					// Matching any pip version
    				r := regexp.MustCompile("pip[0-9]+\\.[0-9]")
				   	a := val.(string)
				   	options := [4]string{"docker", "pip", "pip3", "make"}
				   	var found bool

				   	for _, word := range options {
						if word == a {
							found = true
							break
						}

						if r.MatchString(word) {
							found = true
							break
						}

				   	if !found {
				    	errs = append(errs, fmt.Errorf("%q Must be one of: %v", a, options))
				   	}
				   	return
				},
			},
			"requirements_path": {
				// This description is used by the documentation generator and the language server.
				Description: "your requirements, Makefile folder",
				Type:        schema.TypeString,
				Required:    true,
			},
			"artifact_type": {
				Description: "If your artifact is a layer or a function"
				Type:        schema.TypeString
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					a := val.(string)
					options := [2]string{"layer", "function"}
					var found bool

					for _, word := range options {
						if word == a {
							found = true
							break
						}

				   	if !found {
				    	errs = append(errs, fmt.Errorf("%q Must be one of: %v", a, options))
				   	}
				   	return
				},
			}
			"artifact_name": {
				// This description is used by the documentation generator and the language server.
				Description: "your requirements, Makefile path",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				   a := val.(string)

				   fileExtension := path.Ext(a)
				   if fileExtension != "" {
				   		errs = append(errs, fmt.Errorf("%q Must not contain extension, only the file name", a))
				   }
				   return
				},
			},
			"sha1": {
				Description: "File description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"zip_path_location": {
				Description: "Generated zip file",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	var diags diag.Diagnostics

	// initializing data metadata
	data := make(map[string]string)

	requirements_path       := d.Get("requirements_path").(string)
	runtime_execution_build := d.Get("build_method").(string)
	artifact_name           := fmt.Sprintf("%s.zip", d.Get("artifact_name").(string))
	runtime_execution       := d.Get("runtime").(string)

	s := fmt.Sprintf("which %s", runtime_execution)

	if _, err := os.Stat(requirements_path); os.IsNotExist(err) {
	  return diag.Errorf("Folder containing requirements does not exists")
	}

	_, err := exec.Command("bash", "-c", s).Output()
	if err != nil {
        return diag.Errorf("Problem to execute runtime defined: %s", err)
    }

    // Matching any pip version
    re := regexp.MustCompile("pip[0-9]+\\.[0-9]")

    switch runtimes := runtime_execution_build; runtimes {
    	case "make":
			m := fmt.Sprintf("make build")
			_, err := exec.Command("bash", "-c", m).Output()
			if err != nil {
		        return diag.Errorf("Problem to execute make build: %s", err)
		    }

		case "docker":
			buf := bytes.Buffer{}
			data = map[string]string{
				"ArtifactDir": requirements_path,
				"PythonRuntime": runtime_execution,
				"ExecutionRuntime": "pip",
			}
			t := template.Must(template.New("").Parse(shInstallDepsTemplate))
			t.Execute(&buf, data)
			var sh string = buf.String()
			data["InstallTemplate"] = sh
			buf = bytes.Buffer{}
			t = template.Must(template.New("").Parse(dockerInstallDepsTemplate))
			t.Execute(&buf, data)
			var dockerSh string = buf.String()
			out, err := exec.Command("bash", "-c", dockerSh).Output()
			if err != nil {
		        return diag.Errorf("Problem to execute docker build: %s, path: %s, out: %s", err, requirements_path, out)
		    }

		case re.FindString(runtimes):
			buf := bytes.Buffer{}

			data = map[string]string{
				"ArtifactDir": requirements_path,
				"PythonRuntime": runtime_execution,
				"ExecutionRuntime": runtime_execution_build,
			}
			t := template.Must(template.New("").Parse(shInstallDepsTemplate))
			t.Execute(&buf, data)
			var sh string = buf.String()
			out, err := exec.Command("bash", "-c", sh).Output()
			if err != nil {
		        return diag.Errorf("Problem to execute pip build: %s, path: %s, out: %s", err, requirements_path, out)
		    }
		default:
			return diag.Errorf("Not immplemented")
	}
	output_dir := fmt.Sprintf("/tmp/%s", artifact_name)

	_, err = zip_wrapper(requirements_path, output_dir)

	if err != nil {
		return diag.Errorf("Problem to generate zip: Err %s", err)
	}

	sha1, err := hash_file_sha1(output_dir)
	if err != nil {
		return diag.Errorf("Problem to generate zip shasum metadata: Err %s", err)
	}

	Id := fmt.Sprintf("%s-%s-%s", runtime_execution_build, runtime_execution, artifact_name)
	d.SetId(Id)
	d.Set("sha1", sha1)
	d.Set("zip_path_location", output_dir)
	return diags
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	var diags diag.Diagnostics

	sha1 := d.Get("sha1").(string)

	f := fmt.Sprintf("/tmp/%s.zip", d.Get("artifact_name").(string))

	if _, err := os.Stat(f); os.IsNotExist(err) {
		d.SetId("")
		d.Set("sha1", "")
		return nil
	}

	sha1_check, err := hash_file_sha1(f)
	if err != nil {
		d.SetId("")
		d.Set("sha1", "")
		return diags
	}

	if sha1_check != sha1 {
		d.SetId("")
		d.Set("sha1", "")
		return diags
	}

	return diags
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	var diags diag.Diagnostics
	// initializing data metadata
	data                    := make(map[string]string)
	artifact_name           := fmt.Sprintf("%s.zip", d.Get("artifact_name").(string))
	requirements_path       := d.Get("requirements_path").(string)
	buf                     := bytes.Buffer{}
	const shTemplate = `
		rm -rf "{{.ArtifactDir}}/python" || echo "ok"
		rm -rf "/tmp/{{.ArtifactName}}" || echo "ok"
	`

	data = map[string]string{
		"ArtifactDir": requirements_path,
		"ArtifactName": artifact_name,
	}

	t := template.Must(template.New("").Parse(shTemplate))
	t.Execute(&buf, data)
	var sh string = buf.String()
	out, err := exec.Command("bash", "-c", sh).Output()
	if err != nil {
		return diag.Errorf("Problem to remove artifacts: err %s, out: %s", err, out)
	}
	return diags
}
