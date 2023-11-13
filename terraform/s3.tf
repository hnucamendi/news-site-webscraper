resource "aws_s3_bucket" "ws_colly_lambda_zip_bucket" {
  bucket = "ws-colly-lambda-zip-bucket"
  tags = {
    Name = "ws-colly-lambda-zip-bucket"
    Environment = "production"
  }
}