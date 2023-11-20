# resource "aws_iam_role" "ws_colly_lambda_role" {
#   name               = "ws-colly-lambda-role"
#   assume_role_policy = data.aws_iam_policy_document.ws_colly_lambda_policy.json
# }
#
# data "aws_iam_policy_document" "ws_colly_lambda_policy" {
#   statement {
#     effect = "Allow"
#     principals {
#       type        = "Service"
#       identifiers = ["lambda.amazonaws.com", "scheduler.amazonaws.com"]
#     }
#     actions = ["sts:AssumeRole"]
#   }
# }
#
# resource "aws_iam_policy" "ws_colly_lambda_sqs_send_message_policy" {
#   name        = "ws_colly_lambda_sqs_send_message_policy"
#   description = "Allows Lambda function to send messages to SQS queue"
#
#   policy = jsonencode({
#     Version = "2012-10-17",
#     Statement = [
#       {
#         Action = [
#           "sqs:SendMessage",
#           "ssm:GetParameter"
#         ],
#         Effect   = "Allow",
#         Resource = [
#           aws_sqs_queue.ws_colly_sqs_queue.arn,
#           aws_sqs_queue.ws_colly_sqs_dead_queue.arn,
#           aws_ssm_parameter.sqs-url.arn,
#         ]
#       },
#     ]
#   })
# }
#
# resource "aws_lambda_function" "ws_colly_lambda" {
#   function_name    = "ws-colly"
#   s3_bucket        = aws_s3_bucket.lambda_package_zip_bucket.id
#   s3_key           = "ws-colly-lambda/bootstrap.zip"
#   role             = aws_iam_role.ws_colly_lambda_role.arn
#   handler          = "main.HandleRequest"
#   runtime          = "provided.al2"
#     environment {
#     variables = {
#       SQSqueueName = aws_sqs_queue.ws_colly_sqs_queue.url
#     }
#   }
# }
#
# resource "aws_cloudwatch_event_rule" "ws_colly_lambda_daily_event" {
#   name                = "ws-colly-lambda-daily-event"
#   description         = "Fires every day"
#   schedule_expression = "rate(8 hours)"
# }
#
# resource "aws_cloudwatch_event_target" "invoke_ws_colly_lambda_every_five_minutes" {
#   rule      = aws_cloudwatch_event_rule.ws_colly_lambda_daily_event.name
#   target_id = "invoke_lambda"
#   arn       = aws_lambda_function.ws_colly_lambda.arn
# }
#
# resource "aws_lambda_permission" "allow_cloudwatch_to_call_ws_colly_lambda" {
#   statement_id  = "AllowExecutionFromCloudWatch"
#   action        = "lambda:InvokeFunction"
#   function_name = aws_lambda_function.ws_colly_lambda.function_name
#   principal     = "events.amazonaws.com"
#   source_arn    = aws_cloudwatch_event_rule.ws_colly_lambda_daily_event.arn
# }
#
# resource "aws_iam_role_policy_attachment" "ws_colly_lambda_sqs_send_message_policy_attachment" {
#   role       = aws_iam_role.ws_colly_lambda_role.name
#   policy_arn = aws_iam_policy.ws_colly_lambda_sqs_send_message_policy.arn
# }
#
# # resource "aws_iam_role_policy_attachment" "ws_colly_lambda_sqs_policy_attachment" {
# #   role       = aws_iam_role.ws_colly_lambda_role.name
# #   policy_arn =  "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
# # }
#
# resource "aws_iam_role_policy_attachment" "ws_colly_lambda_cloudwatch_policy_attachment" {
#   role       = aws_iam_role.ws_colly_lambda_role.name
#   policy_arn =  "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
# }
#
# resource "aws_iam_role_policy_attachment" "lambda_sqs_role_policy" {
#   role       = aws_iam_role.ws_colly_lambda_role.name
#   policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
# }
#
# resource "aws_lambda_function_event_invoke_config" "ws_colly_lambda_event_invoke_config" {
#   function_name = aws_lambda_function.ws_colly_lambda.function_name
#
#   destination_config {
#
#       on_success {
#       destination = aws_sqs_queue.ws_colly_sqs_queue.arn
#     }
#
#     on_failure {
#       destination = aws_sqs_queue.ws_colly_sqs_dead_queue.arn
#     }
#   }
# }
#
# resource "aws_ssm_parameter" "sqs-url" {
#   name  = "/lambda/prod/ws-colly/ws-colly-lambda-sqs-url"
#   type  = "SecureString"
#   value = aws_sqs_queue.ws_colly_sqs_queue.url
# }