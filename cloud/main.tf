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

module "backend" {
    source = "./modules/backend"
}