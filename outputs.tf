output "public_ip" {
  value       = aws_instance.web.public_dns 
  description = "The public IP of the web server"
}

output "clb_dns_name" {
  value       = aws_lb.web_lb.dns_name
  description = "The domain name of the load balancer"
}