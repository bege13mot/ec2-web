
# -----------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# -----------------------------------------

# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# ---------------------------------------------------------------------------------------------------------------------


variable "aws_region" {
    description = "AWS Region code"
    type = string
    default = "us-west-2"
}

variable "app_port" {
    description = "Application port"
    type = number
    default = 8080
}