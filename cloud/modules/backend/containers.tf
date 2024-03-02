resource "aws_ecr_repository" "langlearn_ecr" {
    name = "langlearn_ecr"
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_ecs_cluster" "langlearn_cluster" {
    name = "langlearn_cluster"
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_ecs_task_definition" "langlearn_td" {
    family = "langlearn_td_family"
    requires_compatibilities = ["FARGATE"]
    network_mode = "awsvpc"
    cpu = 256
    memory = 512
    task_role_arn = ""
    execution_role_arn = ""
    runtime_platform {
        cpu_architecture = "X86_64"
        operating_system_family = "LINUX"
    }
    container_definitions = jsonencode([
        {
            essential = true,
            name = "langlearn_backend",
            image = format(
                "%s:%s", 
                aws_ecr_repository.langlearn_ecr.repository_url, 
                "backend"
            ),
            cpu = 256,
            memory = 512,
            portMappings = [
                {
                    containerPort = 80,
                    hostPort = 80,
                    protocol = "tcp"
                }
            ]
        }
    ])
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_security_group" "allow_web_traffic" {
    name = "langlearn_allow_web_traffic"
    description = "allow web requests"
    vpc_id = aws_vpc.backend_vpc.id
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_vpc_security_group_ingress_rule" "allow_http" {
    security_group_id = aws_security_group.allow_web_traffic.id
    cidr_ipv4 = "0.0.0.0/0"
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_vpc_security_group_ingress_rule" "allow_https" {
    security_group_id = aws_security_group.allow_web_traffic.id
    cidr_ipv4 = "0.0.0.0/0"
    from_port = 443
    to_port = 443
    ip_protocol = "tcp"
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_vpc_security_group_egress_rule" "allow_any_traffic" {
    security_group_id = aws_security_group.allow_web_traffic.id
    cidr_ipv4 = "0.0.0.0/0"
    ip_protocol = "-1"
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_ecs_service" "langlearn_backend" {
    name = "langlearn_backend"
    cluster = aws_ecs_cluster.langlearn_cluster.id
    task_definition = aws_ecs_task_definition.langlearn_td.arn
    desired_count = 1
    force_new_deployment = true
    health_check_grace_period_seconds = 30
    launch_type = "FARGATE"
    wait_for_steady_state = true
    network_configuration {
      subnets = [aws_subnet.subnet_private_one.id, aws_subnet.subnet_private_two.id]
      security_groups = [aws_security_group.allow_web_traffic.id]
      assign_public_ip = true
    }
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}