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

resource "aws_ecr_repository" "ecr_repo" {
    name = "langlearn"
    image_tag_mutability = "MUTABLE"
    image_scanning_configuration {
        scan_on_push = false
    }
    tags = {
        app = "langlearn"
    }
}

resource "aws_iam_user" "langlearn_user_admin" {
    name = "admin"
}

data "aws_iam_policy_document" "langlearn_policy_pc" {
    statement {
        effect = "Allow"
        actions = ["iam:ChangePassword"]
        resources = ["arn:aws:iam::*:user/&{aws:username}"]
    }
    statement {
        effect = "Allow"
        actions = ["iam:GetAccountPasswordPolicy"]
        resources = ["*"]
    }
}

data "aws_iam_policy_document" "langlearn_policy_admin" {
    statement {
        effect = "Allow"
        actions = ["*"]
        resources = ["*"]
    }
}

resource "aws_iam_policy" "langlearn_policy_pc" {
    name = "password-change-policy"
    policy = data.aws_iam_policy_document.langlearn_policy_pc.json
}

resource "aws_iam_policy" "langlearn_policy_admin" {
    name = "admin-policy"
    policy = data.aws_iam_policy_document.langlearn_policy_admin.json
}

resource "aws_iam_policy_attachment" "langlearn_policy_admin_attachment" {
    name = "admin-policy-attachment"
    users = [aws_iam_user.langlearn_user_admin.name]
    policy_arn = aws_iam_policy.langlearn_policy_admin.arn
}

resource "aws_iam_policy_attachment" "langlearn_policy_pc_attachment" {
    name = "passowrd-change-policy-attachment"
    users = [aws_iam_user.langlearn_user_admin.name]
    policy_arn = aws_iam_policy.langlearn_policy_pc.arn
}

# resource "aws_ecs_task_definition" "langlearn_td" {
#     family = "langlearn_td_family"
#     requires_compatibilities = ["FARGATE"]
#     network_mode = "awsvpc"
#     cpu = 1024
#     memory = 3072
#     runtime_platform {
#         cpu_architecture = "X86_64"
#         operating_system_family = "LINUX"
#     }
#     container_definitions = {
#     }
# }



























