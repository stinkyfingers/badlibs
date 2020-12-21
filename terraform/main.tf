# vars
variable "region" {
  type = "string"
  default = "us-west-1"
}

# provider
provider "aws" {
  profile = "jds"
  region     = "${var.region}"
}

# import
module "stinkyfingers" {
  source = "../../../../../../infrastructure/stinkyfingers"
}

# Lambda
resource "aws_lambda_permission" "badlibs_server" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.badlibs_server.arn}"
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn = "${aws_lb_target_group.badlibs.arn}"
}

resource "aws_lambda_permission" "badlibs_server_live" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_alias.badlibs_server_live.arn}"
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn = "${aws_lb_target_group.badlibs.arn}"
}

resource "aws_lambda_alias" "badlibs_server_live" {
  name             = "live"
  description      = "set a live alias"
  function_name    = "${aws_lambda_function.badlibs_server.arn}"
  function_version = "${aws_lambda_function.badlibs_server.version}"
}

resource "aws_lambda_function" "badlibs_server" {
  filename         = "../lambda.zip"
  function_name    = "badlibs_server"
  role             = "${aws_iam_role.lambda_role.arn}"
  handler          = "lambda-lambda"
  runtime          = "go1.x"
  source_code_hash = "${filebase64sha256("../lambda.zip")}"
  environment {
    variables = {
      OPENSHIFT_MONGODB_DB_HOST     =  "${data.aws_ssm_parameter.db_host.value}"
      OPENSHIFT_MONGODB_DB_PASSWORD =  "${data.aws_ssm_parameter.db_password.value}"
      OPENSHIFT_MONGODB_DB_PORT     =  "${data.aws_ssm_parameter.db_port.value}"
      OPENSHIFT_MONGODB_DB_USERNAME =  "${data.aws_ssm_parameter.db_username.value}"
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

# ALB
resource "aws_lb_target_group" "badlibs" {
  name        = "badlibs"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "badlibs_server" {
  target_group_arn  = "${aws_lb_target_group.badlibs.arn}"
  target_id         = "${aws_lambda_alias.badlibs_server_live.arn}"
  depends_on        = ["aws_lambda_permission.badlibs_server_live"]
}

resource "aws_lb_listener_rule" "badlibs_server" {
listener_arn = "${module.stinkyfingers.stinkyfingers_http_listener}"
priority = 2
  action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.badlibs.arn}"
  }
  condition {
    field = "path-pattern"
    values = ["/badlibs/*"]
  }
  depends_on = ["aws_lb_target_group.badlibs"]
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

data "aws_ssm_parameter" "db_host" {
  name = "/badlibs/MONGODB_DB_HOST"
}
data "aws_ssm_parameter" "db_password" {
  name = "/badlibs/MONGODB_DB_PASSWORD"
}
data "aws_ssm_parameter" "db_port" {
  name = "/badlibs/MONGODB_DB_PORT"
}
data "aws_ssm_parameter" "db_username" {
  name = "/badlibs/MONGODB_DB_USERNAME"
}
