data "archive_file" "trama" {
  type        = "zip"
  source_file = "${path.module}/../build/lambda/bootstrap"
  output_path = "${path.module}/../build/lambda/function.zip"
}

resource "aws_lambda_function" "trama" {
  depends_on = [aws_cloudwatch_log_group.trama]

  function_name    = "${var.tags["Project"]}-api-${var.tags["Environment"]}"
  role             = aws_iam_role.lambda.arn
  filename         = data.archive_file.trama.output_path
  source_code_hash = data.archive_file.trama.output_base64sha256
  handler          = "bootstrap"
  runtime          = "provided.al2023"

  timeout         = var.lambda_timeout
  memory_size     = var.lambda_memory_size

  environment {
    variables = {
      GIN_MODE       = "release"
      DYNAMODB_USERS_TABLE_NAME = aws_dynamodb_table.users.name
      JWT_SECRET     = var.jwt_secret
      ENVIRONMENT    = var.tags["Environment"]
    }
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "trama" {
  name              = "/aws/lambda/${var.tags["Project"]}-api-${var.tags["Environment"]}"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.trama.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.trama.execution_arn}/*/*"
}
