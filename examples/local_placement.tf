# style_guide_local_placement

resource "aws_instance" "web" {
  instance_type = local.instance_type
}

locals {
  instance_type = "t3.micro"
}
