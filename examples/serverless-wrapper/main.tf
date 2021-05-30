terraform {
  	required_providers {
   	wrapper = {
    	versions = ["0.3"]
    	source = "hashicorp.com/edu/wrapper"
    }
  }
}

provider "aws" {}

provider "wrapper" {}

resource "wrapper" "foo-layer" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./layer/"
  artifact_name = "foo-layer"
  artifact_type = "layer"
}

resource "wrapper" "foo-function" {
  runtime = "python3.8"
  build_method = "pip"
  requirements_path = "./src-function/"
  artifact_name = "foo-function-layer"
  artifact_type = "function"
}

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