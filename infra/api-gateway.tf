locals {
  cors_origins = {
    dev  = ["*"]
    prod = ["https://trama.app"]
  }
}

resource "aws_apigatewayv2_api" "trama" {
  name          = "${var.tags["Project"]}-api-${var.tags["Environment"]}"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = local.cors_origins[var.tags["Environment"]]
    allow_methods = ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
    allow_headers = ["*"]
    max_age       = 3600
  }

  tags = var.tags
}

resource "aws_apigatewayv2_stage" "trama" {
  api_id      = aws_apigatewayv2_api.trama.id
  name        = "$default"
  auto_deploy = true

  default_route_settings {
    throttling_burst_limit = var.throttling_burst_limit
    throttling_rate_limit  = var.throttling_rate_limit
  }
}

resource "aws_apigatewayv2_integration" "trama" {
  api_id              = aws_apigatewayv2_api.trama.id
  integration_type    = "AWS_PROXY"
  integration_uri     = aws_lambda_function.trama.invoke_arn
  integration_method  = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "trama" {
  api_id    = aws_apigatewayv2_api.trama.id
  route_key = "ANY /api/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.trama.id}"
}
