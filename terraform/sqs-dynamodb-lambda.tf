resource "aws_iam_role" "lambda_assume_role" {
  name               = "lambda-assume-role"
  assume_role_policy = data.aws_iam_policy_document.ws_colly_lambda_policy.json
}

data "aws_iam_policy_document" "lambda_assume_role_policy" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_policy" "lambda_sqs_receive_message_policy" {
  name        = "lambda_sqs_receive_message_policy"
  description = "Allows Lambda function to receive messages to dynamodb"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "ssm:GetParameter",
          "dynamodb:*",
        ],
        Effect   = "Allow",
        Resource = [
          aws_sqs_queue.ws_colly_sqs_queue.arn,
          aws_sqs_queue.ws_colly_sqs_dead_queue.arn,
          aws_ssm_parameter.sqs-url.arn,
          aws_dynamodb_table.ws_colly_dynamo_table.arn,
        ]
      },
    ]
  })
}
resource "aws_iam_role_policy_attachment" "lambda_cloudwatch_policy_attachment" {
  role       = aws_iam_role.lambda_assume_role.name
  policy_arn =  "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_receive_message_policy_attachment" {
  role       = aws_iam_role.lambda_assume_role.name
  policy_arn = aws_iam_policy.lambda_sqs_receive_message_policy.arn
}

# resource "aws_iam_role_policy_attachment" "lambda_sqs_policy_attachment" {
#   role       = aws_iam_role.lambda_assume_role.name
#   policy_arn =  "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
# }

resource "aws_lambda_event_source_mapping" "lambda_event_source_mapping" {
  event_source_arn = aws_sqs_queue.ws_colly_sqs_queue.arn
  function_name    = aws_lambda_function.json_parser_to_dynamodb_lambda.function_name
}

resource "aws_lambda_function" "json_parser_to_dynamodb_lambda" {
  function_name    = "json-parser-to-dynamodb-lambda"
  s3_bucket        = aws_s3_bucket.lambda_package_zip_bucket.id
  s3_key           = "sqs-to-dynamodb/bootstrap.zip"
  role             = aws_iam_role.lambda_assume_role.arn
  handler          = "main.HandleRequest"
  runtime          = "provided.al2"
}

