terraform {
  required_providers {
    wrapper = {}
  }
}

provider "wrapper" {}

resource "wrapper_layer" "foo" {
  runtime = "python3.8"
  build_method = "pip3"
  requirements_path = "./src/"
  artifact_name = "foo"
}
