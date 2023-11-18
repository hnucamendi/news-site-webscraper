# resource "aws_api_gateway_rest_api" "ws_colly_dynamo_api" {
#   name          = "ws-colly-dynamo-api"
# }

# resource "aws_api_gateway_resource" "ws_colly_dynamo_api_top_headlines_resource" {
#   rest_api_id = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   parent_id   = aws_api_gateway_rest_api.ws_colly_dynamo_api.root_resource_id
#   path_part   = "top-headlines"
# }

# resource "aws_api_gateway_method" "ws_colly_dynamo_api_top_headlines_method" {
#   rest_api_id   = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
#   resource_id   = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
#   http_method   = "GET"
#   authorization = "NONE"
# }

# # resource "aws_api_gateway_integration" "ws_colly_dynamo_api_top_headlines_integration" {
# #   rest_api_id          = aws_api_gateway_rest_api.ws_colly_dynamo_api.id
# #   resource_id          = aws_api_gateway_resource.ws_colly_dynamo_api_top_headlines_resource.id
# #   http_method          = aws_api_gateway_method.ws_colly_dynamo_api_top_headlines_method.http_method
# #   type                 = "MOCK"
# #   cache_key_parameters = ["method.request.path.param"]
# #   cache_namespace      = "foobar"
# #   timeout_milliseconds = 29000

# #   request_parameters = {
# #     "integration.request.header.X-Authorization" = "'static'"
# #   }

# #   # Transforms the incoming XML request to JSON
# #   request_templates = {
# #     "application/xml" = <<EOF
# # {
# #    "body" : $input.json('$')
# # }
# # EOF
# #   }
# # }