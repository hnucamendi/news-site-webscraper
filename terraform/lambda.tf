resource "aws_lambda_function" "ws-colly_lambda" {
  filename         = "main.zip"
  function_name    = "ws-colly"
  role             = "${aws_iam_role.ws-colly_lambda_role.arn}"
  handler          = "lambda_function.lambda_handler"
  source_code_hash = "${base64sha256(file("main.zip"))}"
  runtime          = "provided.al2"
  timeout          = "60"
  memory_size      = "128"
  publish          = "true"
}