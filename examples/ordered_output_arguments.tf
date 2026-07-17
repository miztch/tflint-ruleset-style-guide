# style_guide_ordered_output_arguments

output "instance_ip_addr" {
  value       = aws_instance.web.private_ip
  description = "The private IP address of the instance"
}
