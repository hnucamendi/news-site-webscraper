# resource "aws_api_gateway_rest_api" "ws_colly_dynamo_api" {
#   name          = "ws-colly-dynamo-api"
# }
#
# resource "aws_api_gateway_resource" "ws_colly_dynamo_api_top_headlines_resource" {
#   rest_api_id = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   parent_id   = aws_api_gateway_rest_api.ws_colly_dynamo_api.root_resource_id
#   path_part   = "top-headlines"
# }
#
# resource "aws_api_gateway_method" "ws_colly_dynamo_api_top_headlines_method" {
#   rest_api_id   = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   resource_id   = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
#   http_method   = "GET"
#   authorization = "NONE"
# }
#
# resource "aws_api_gateway_integration" "ws_colly_dynamo_api_top_headlines_integration" {
#   rest_api_id             = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   resource_id             = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
#   http_method             = aws_api_gateway_method.ws_colly_dynamo_api_top_headlines_method.http_method
#   type                    = "AWS"
#   integration_http_method = "POST"
# #   uri                     = aws_dynamodb_table.ws_colly_dynamo_table.arn
#   uri                     = "arn:aws:apigateway:us-east-1:dynamodb:action/GetItem"
#   credentials             = aws_iam_role.api_gateway_assume_role.arn
#   passthrough_behavior    = "WHEN_NO_TEMPLATES"
#
#   request_templates = {
#       "application/json" = <<EOF
#   {
#       "TableName": "${aws_dynamodb_table.ws_colly_dynamo_table.name}",
#       "Key": {
#           "ID": {
#               "S": "$input.params('id')"
#           }
#       }
#   }
#   EOF
#     }
# }
#
# resource "aws_api_gateway_integration_response" "ws_colly_dynamo_api_top_headlines_integration_response" {
#   rest_api_id = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   resource_id = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
#   http_method = aws_api_gateway_method.ws_colly_dynamo_api_top_headlines_method.http_method
#   status_code = "200"
# }
#
# resource "aws_api_gateway_method_response" "ws_colly_dynamo_api_top_headlines_method_response" {
#   rest_api_id = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   resource_id = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
#   http_method = aws_api_gateway_method.ws_colly_dynamo_api_top_headlines_method.http_method
#   status_code = "200"
#
#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Origin" = true
#   }
# }
#
# resource "aws_api_gateway_deployment" "ws_colly_dynamo_api_top_headlines_deployment" {
#   rest_api_id = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#
#   lifecycle {
#     create_before_destroy = true
#   }
# }
#
# resource "aws_api_gateway_stage" "ws_colly_dynamo_api_top_headlines_stage" {
#   deployment_id = aws_api_gateway_deployment.ws_colly_dynamo_api_top_headlines_deployment.id
#   rest_api_id   = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   stage_name    = "ws-colly-dynamo-api-top-headlines-stage"
# }
#
# resource "aws_iam_role" "api_gateway_assume_role" {
#   name               = "api-gateway-assume-role"
#   assume_role_policy = data.aws_iam_policy_document.api_gateway_assume_role_policy_document.json
# }
#
# data "aws_iam_policy_document" "api_gateway_assume_role_policy_document" {
#   statement {
#     effect = "Allow"
#     principals {
#       type        = "Service"
#       identifiers = ["apigateway.amazonaws.com"]
#     }
#     actions = ["sts:AssumeRole"]
#   }
# }
#
# # resource "aws_iam_policy" "api_gateway_role_policy" {
# #   name        = "api_gateway_role_policy"
# #   description = "Allows Lambda function to receive messages to dynamodb"
# #
# #   policy = jsonencode({
# #     Version = "2012-10-17",
# #     Statement = [
# #       {
# #         Action = [
# #           "sqs:ReceiveMessage",
# #           "sqs:DeleteMessage",
# #           "sqs:GetQueueAttributes",
# #           "ssm:GetParameter",
# #           "dynamodb:*",
# #           "execute-api:Invoke",
# #           "execute-api:ManageConnections"
# #         ],
# #         Effect   = "Allow",
# #         Resource = [
# #           aws_sqs_queue.ws_colly_sqs_queue.arn,
# #           aws_sqs_queue.ws_colly_sqs_dead_queue.arn,
# #           aws_ssm_parameter.sqs-url.arn,
# #           aws_dynamodb_table.ws_colly_dynamo_table.arn,
# #         ]
# #       },
# #     ]
# #   })
# # }