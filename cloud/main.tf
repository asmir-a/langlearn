terraform {
    cloud {
        organization = "asmir4development-org"
        workspaces {
            project = "langlearn"
            name = "main"
        }
    }
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 4.16"
        }
    }
    required_version = ">= 1.2.0"
}

provider "aws" {
    region = "ap-northeast-2"
}

resource "aws_instance" "app_server" {
    ami = "ami-0f3a440bbcff3d043"
    instance_type = "t2.micro"
}
