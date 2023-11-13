# resource "aws_lambda_invocation" "ws_colly_lambda_invocation" {
#   function_name = aws_lambda_function.ws_colly_lambda.function_name

#   input = jsonencode({})

#   triggers = {
#     # Every time the lambda function is updated, the invocation will be triggered
#     # This is a workaround for the fact that terraform does not support
#     # aws_lambda_permission as a data source
#   }
# }

# output "result_entry" {
#   value = jsondecode(aws_lambda_invocation.ws_colly_lambda_invocation.result)
# }

# resource "aws_scheduler_schedule" "ws_colly_scheduler" {
#   name       = "ws-colly-scheduler"
#   group_name = aws_scheduler_schedule_group.ws_colly_lambda_wait_group.name

#   flexible_time_window {
#     mode = "FLEXIBLE"
#     maximum_window_in_minutes = 30
#   }

#   schedule_expression = "rate(5 minutes)"

#   target {
#     arn      = aws_lambda_function.ws_colly_lambda.arn
#     role_arn = aws_iam_role.ws_colly_lambda_role.arn
#   }
# }

# resource "aws_scheduler_schedule_group" "ws_colly_lambda_wait_group" {
#   name = "ws-colly-lambda-wait-group"

#   tags = {
#     Name = "ws-colly-lambda-wait-group"
#     Environment = "production"
#   }
# }

resource "aws_iam_role" "ws_colly_lambda_role" {
  name               = "ws-colly-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.ws_colly_lambda_policy.json
}

data "aws_iam_policy_document" "ws_colly_lambda_policy" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com", "scheduler.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_lambda_function" "ws_colly_lambda" {
  function_name    = "ws-colly"
  s3_bucket        = aws_s3_bucket.ws_colly_lambda_zip_bucket.id
  s3_key           = "bootstrap.zip"
  role             = aws_iam_role.ws_colly_lambda_role.arn
  handler          = "main.HandleRequest"
  runtime          = "provided.al2"
  # dead_letter_config {
  #   target_arn = aws_sqs_queue.ws_colly_sqs_dead_queue.arn
  # }
}

resource "aws_cloudwatch_event_rule" "ws_colly_lambda_daily_event" {
  name                = "ws-colly-lambda-daily-event"
  description         = "Fires every day"
  schedule_expression = "rate(1 day)"
}

resource "aws_cloudwatch_event_target" "invoke_ws_colly_lambda_every_five_minutes" {
  rule      = aws_cloudwatch_event_rule.ws_colly_lambda_daily_event.name
  target_id = "invoke_lambda"
  arn       = aws_lambda_function.ws_colly_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ws_colly_lambda" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ws_colly_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.ws_colly_lambda_daily_event.arn
}