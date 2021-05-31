# Terraform Provider for wrapper serverless

I am just a python guy, trying to solve problems using golang.

## Objective

Install python layers and requirements for aws serverless functions
It will use pip, make or docker to install your code dependencies.

ie:

Copy the `terraform-provider-wrapper` to the following dir:
*Not if the dir does not exist, try create it first, the following example we will create it*
*it it is mac, try change `linux_amd64` to `darwin_amd64`*

```bash
mkdir  ~/.terraform.d/plugins/example.com/svl/wrapper/0.3/linux_amd64/
cp -rf terraform-provider-wrapper ~/.terraform.d/plugins/example.com/svl/wrapper/0.3/linux_amd64/
```
or
```bash
make
# If mac
make OS_ARCH=darwin_amd64
```

then:

```bash
terraform {
  	required_providers {
   	wrapper = {
    	versions = ["0.1"]
    	source = "example.com/svl/wrapper"
    }
  }
}

provider "aws" {}

provider "wrapper" {}

// Instaling a layer ie: `./test_tf/layer/`
resource "wrapper" "foo-layer" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./layer/"
  artifact_name = "foo-layer"
  artifact_type = "layer"
}

// Instaling a layer ie: `./test_tf/src-function/`
resource "wrapper" "foo-function" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./src-function/"
  artifact_name = "foo-function-layer"
  artifact_type = "function"
}

// Instaling a layer ie: `./test_tf/src-layer/`
resource "wrapper" "foo-function-no-layer" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./src-layer/"
  artifact_name = "foo-function"
  artifact_type = "function"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_layer_version" "lambda_layer" {
  filename   = wrapper.foo-layer.zip_path_location
  layer_name = "foo-layer"

  compatible_runtimes = ["python3.8"]
}

resource "aws_lambda_function" "function-with-layer" {
  filename      = wrapper.foo-function-no-layer.zip_path_location
  function_name = "foo-function-no-layer"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "main.lambda_handler"

  layers = [aws_lambda_layer_version.lambda_layer.arn]

  source_code_hash = filebase64sha256("${wrapper.foo-function-no-layer.zip_path_location}")

  runtime = "python3.8"

}
```
