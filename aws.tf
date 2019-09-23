provider "aws" {
    version = "~> 2.0"
    region     = "${var.aws_region}"
}

resource "aws_key_pair" "ssh" {
  key_name   = "local"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}


resource "aws_security_group" "web" {
  name        = "webserver"
  description = "Public HTTP + SSH"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = "${var.app_port}"
    to_port     = "${var.app_port}"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

}

resource "aws_security_group" "lb" {
  name        = "LB"
  description = "Public HTTP"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = "${var.app_port}"
    to_port     = "${var.app_port}"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "web" {
 
  ami                    = "ami-08d489468314a58df" #Amazon Linux 2 comes with five years support
  instance_type          = "t2.micro"
  key_name               = "${aws_key_pair.ssh.id}"
  vpc_security_group_ids = [ "${aws_security_group.web.id}" ]
  user_data = "${file("install.sh")}"

}

data "aws_vpc" "default" {
  default = true
}

resource "aws_subnet" "zone-a" {
    vpc_id = "${data.aws_vpc.default.id}"

    cidr_block = "172.31.1.0/24"
    availability_zone = "${var.aws_region}a"

    # tags {
    #     Name = "Public Subnet A"
    # }
}

resource "aws_subnet" "zone-b" {
    vpc_id = "${data.aws_vpc.default.id}"

    cidr_block = "172.31.2.0/24"
    availability_zone = "${var.aws_region}b"

    # tags {
    #     Name = "Public Subnet B"
    # }
}


resource "aws_lb" "web_lb" {
  name               = "terraform-lb"
  security_groups    =  [ "${aws_security_group.lb.id}" ]
  load_balancer_type = "application"

  subnets            = ["${aws_subnet.zone-a.id}", "${aws_subnet.zone-b.id}"]
}

resource "aws_alb_target_group" "alb_front_web" {
	name	= "target"
	vpc_id	= "${data.aws_vpc.default.id}"
	port	= "${var.app_port}"
	protocol	= "HTTP"
}

resource "aws_alb_target_group_attachment" "alb_web-1" {
  target_group_arn = "${aws_alb_target_group.alb_front_web.arn}"
  target_id        = "${aws_instance.web.id}"
  port             = "${var.app_port}"
}



resource "aws_alb_listener" "alb_front_web" {
	load_balancer_arn	=	"${aws_lb.web_lb.arn}"
	port			=	80
	protocol		=	"HTTP"

  default_action {
		target_group_arn	=	"${aws_alb_target_group.alb_front_web.arn}"
		type			=	"forward"
	}
}
