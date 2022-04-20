# vars
variable "region" {
  type = string
  default = "us-west-1"
}

# provider
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  profile = "jds"
  region     = var.region
}

# import
data "terraform_remote_state" "stinkyfingers" {
  backend = "s3"
  config = {
    bucket  = "remotebackend"
    key     = "stinkyfingers/terraform.tfstate"
    region  = "us-west-1"
    profile = "jds"
  }
}

# Lambda
resource "aws_lambda_permission" "badlibs_server" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.badlibs_server.arn
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    =  aws_lb_target_group.badlibs.arn
}

resource "aws_lambda_permission" "badlibs_server_live" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_alias.badlibs_server_live.arn
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn = aws_lb_target_group.badlibs.arn
}

resource "aws_lambda_alias" "badlibs_server_live" {
  name             = "live"
  description      = "set a live alias"
  function_name    = aws_lambda_function.badlibs_server.arn
  function_version = aws_lambda_function.badlibs_server.version
}

resource "aws_lambda_function" "badlibs_server" {
  filename         = "../lambda.zip"
  function_name    = "badlibs_server"
  role             = aws_iam_role.lambda_role.arn
  handler          = "lambda-lambda"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("../lambda.zip")
  environment {
    variables = {
      PROFILE = "jds"
    }
  }
}

# IAM
resource "aws_iam_role" "lambda_role" {
  name = "badlibs-lambda-role"
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

resource "aws_iam_role_policy_attachment" "cloudwatch-attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "s3-policy" {
  name = "badlibs-lambda-s3-policy"
  description = "Grants lambda access to s3"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": "arn:aws:s3:::*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "s3-policy-attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.s3-policy.arn
}

# ALB
resource "aws_lb_target_group" "badlibs" {
  name        = "badlibs"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "badlibs_server" {
  target_group_arn  = aws_lb_target_group.badlibs.arn
  target_id         = aws_lambda_alias.badlibs_server_live.arn
  depends_on        = [aws_lambda_permission.badlibs_server_live]
}

resource "aws_lb_listener_rule" "badlibs_server" {
listener_arn = data.terraform_remote_state.stinkyfingers.outputs.stinkyfingers_https_listener
priority = 2
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.badlibs.arn
  }
  condition {
    path_pattern {
      values = ["/badlibs/*"]
    }
  }
  depends_on = [aws_lb_target_group.badlibs]
}

# db
resource "aws_s3_bucket" "badlibs" {
  bucket = "badlibs"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Sid": "Lambda Read",
          "Effect": "Allow",
          "Principal": {
              "AWS": "${aws_iam_role.lambda_role.arn}"
          },
          "Action": [
            "s3:*"
          ],
          "Resource": [
            "arn:aws:s3:::badlibs",
            "arn:aws:s3:::badlibs/*"
          ]
      }
  ]
}
EOF
}

# backend
terraform {
  backend "s3" {
    bucket = "remotebackend"
    key    = "badlibs/terraform.tfstate"
    region = "us-west-1"
    profile = "jds"
  }
}

data "terraform_remote_state" "badlibs" {
  backend = "s3"
  config = {
    bucket  = "remotebackend"
    key     = "badlibs/terraform.tfstate"
    region  = "us-west-1"
    profile = "jds"
  }
}
