resource "aws_s3_bucket" "lambda_package_zip_bucket" {
  bucket = "lambda-package-zip-bucket"
  tags = {
    Name = "lambda-package-zip-bucket"
    Environment = "production"
  }
}
