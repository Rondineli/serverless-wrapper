provider "serverless-wrapper" {}

resource "wrapper_resource" "foo" {
  runtime = "foo"
  build_method = "bla"
  requirements_path = "bar"
  artifact_name = "foo"
  use_docker = "bla"
}