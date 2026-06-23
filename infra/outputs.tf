output "api_url" {
  description = "API Gateway endpoint URL"
  value       = aws_apigatewayv2_api.trama.api_endpoint
}

output "function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.trama.function_name
}
