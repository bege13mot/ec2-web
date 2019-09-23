#!/bin/sh
set -e -x

cd /usr/local/
wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
tar -xzf go1.13.linux-amd64.tar.gz
rm go1.13.linux-amd64.tar.gz
chown -R ec2-user:ec2-user /usr/local/

yum update -y
yum install git -y

cd  /home/ec2-user/
git clone https://github.com/bege13mot/ec2-web.git
cd ec2-web

sudo /usr/local/go/bin/go build .
chown -R ec2-user:ec2-user /home/ec2-user/ec2-web/

cp ec2-web.sh /etc/init.d/ec2-web
mkdir /var/log/ec2-web
chown ec2-user:ec2-user /var/log/ec2-web
chmod +x ec2-web
chmod +x /etc/init.d/ec2-web

service ec2-web start
chkconfig ec2-web on