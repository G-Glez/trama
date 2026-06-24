variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-west-1"
}

variable "lambda_timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 30
}

variable "lambda_memory_size" {
  description = "Lambda function memory size in MB"
  type        = number
  default     = 256
}

variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode (PAY_PER_REQUEST or PROVISIONED)"
  type        = string
  default     = "PAY_PER_REQUEST"
}

variable "github_repo" {
  description = "GitHub repo pattern for OIDC auth (e.g. g-glez/*)"
  type        = string
}

variable "jwt_secret" {
  description = "Secret key for JWT signing"
  type        = string
}

variable "throttling_burst_limit" {
  description = "API Gateway burst limit (max requests in a burst)"
  type        = number
  default     = 100
}

variable "throttling_rate_limit" {
  description = "API Gateway rate limit (requests per second)"
  type        = number
  default     = 50
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
