resource "aws_dynamodb_table" "ws_colly_dynamo_table" {
  name = "ws-colly-dynamo-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "ID"

  attribute {
    name = "ID"
    type = "S"
  }
}

