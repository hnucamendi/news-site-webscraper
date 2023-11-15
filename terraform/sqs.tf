resource "aws_sqs_queue" "ws_colly_sqs_queue" {
  name                      = "ws-colly-sqs-queue"
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.ws_colly_sqs_dead_queue.arn
    maxReceiveCount     = 4
  })

  tags = {
    Name = "ws-colly-sqs-queue"
    Environment = "production"
  }
}

resource "aws_sqs_queue" "ws_colly_sqs_dead_queue" {
  name = "ws-colly-sqs-dead-queue"

  tags = {
    Name = "ws-colly-sqs-dead-queue"
    Environment = "production"
  }
}