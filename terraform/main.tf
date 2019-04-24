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
  role             = "${module.stinkyfingers.lambda_role}"
  handler          = "lambda-lambda"
  runtime          = "go1.x"
  source_code_hash = "${filebase64sha256("../lambda.zip")}"
  environment {
    variables = {
      OPENSHIFT_MONGODB_DB_HOST = "ds155862.mlab.com"
      OPENSHIFT_MONGODB_DB_PASSWORD = "Ch1ck3nbutt"
      OPENSHIFT_MONGODB_DB_PORT = "55862"
      OPENSHIFT_MONGODB_DB_USERNAME = "badlibs"
    }
  }
}

# ALB
resource "aws_lb_target_group" "badlibs" {
  name        = "badlibs"
  target_type = "lambda"

  # health_check {
  #    healthy_threshold   = "3"
  #    unhealthy_threshold = "3"
  #    timeout             = "5"
  #    path                = "/"
  #    interval            = "10"
  #    matcher             = "200"
  # }
}

resource "aws_lb_listener" "badlibs_server" {
  # load_balancer_arn = "${aws_lb.stinkyfingers_load_balancer.arn}"
  load_balancer_arn = "${module.stinkyfingers.stinkyfingers_load_balancer}"
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.badlibs.arn}"
  }
}

resource "aws_lb_target_group_attachment" "badlibs_server" {
  target_group_arn  = "${aws_lb_target_group.badlibs.arn}"
  target_id         = "${aws_lambda_alias.badlibs_server_live.arn}"
  depends_on        = ["aws_lambda_permission.badlibs_server_live"]
}

resource "aws_lb_listener_rule" "badlibs_server" {
listener_arn = "${aws_lb_listener.badlibs_server.arn}"
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
