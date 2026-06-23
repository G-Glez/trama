data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "github_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }
    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repo}"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "trama-${var.tags["Environment"]}-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  tags               = var.tags
}

resource "aws_iam_role" "github" {
  name               = "trama-${var.tags["Environment"]}-github"
  assume_role_policy = data.aws_iam_policy_document.github_assume.json
  tags               = var.tags
}

resource "aws_iam_policy" "logs" {
  name        = "trama-${var.tags["Environment"]}-logs"
  description = "Allow Lambda to write CloudWatch logs"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = "arn:aws:logs:${var.aws_region}:*:*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.logs.arn
}

resource "aws_iam_role_policy_attachment" "github_tf" {
  role       = aws_iam_role.github.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
