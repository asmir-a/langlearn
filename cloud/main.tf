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
    id = "whatever1243"
    ami = "ami-08d70e59c07c61a3a"
    instance_type = "t2.micro"
}