# style_guide_ordered_resource_arguments

resource "aws_instance" "web" {
  ebs_block_device {
    device_name = "/dev/sdh"
  }

  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  count = 2
}
