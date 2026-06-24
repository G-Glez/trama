resource "aws_dynamodb_table" "users" {
  name         = "${var.tags["Project"]}-users-${var.tags["Environment"]}"
  billing_mode = var.dynamodb_billing_mode
  hash_key     = "username"

  attribute {
    name = "username"
    type = "S"
  }

  tags = var.tags
}
