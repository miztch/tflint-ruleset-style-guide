resource "aws_instance" "web" {
  count         = 2
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_instance" "incomplete" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  count = 2
}
