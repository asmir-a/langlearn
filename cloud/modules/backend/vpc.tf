resource "aws_vpc" "backend_vpc" {
    cidr_block = "10.0.0.0/16"
    enable_dns_support = true
    enable_dns_hostnames = true
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

data "aws_availability_zones" "available" {
    state = "available"
}

resource "aws_subnet" "subnet_public_one" {
    vpc_id = aws_vpc.backend_vpc.id
    cidr_block = "10.0.0.0/24"
    availability_zone = data.aws_availability_zones.available.names[0]
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_subnet" "subnet_public_two" {
    vpc_id = aws_vpc.backend_vpc.id
    cidr_block = "10.0.1.0/24"
    availability_zone = data.aws_availability_zones.available.names[1]
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_subnet" "subnet_private_one" {
    vpc_id = aws_vpc.backend_vpc.id
    cidr_block = "10.0.2.0/24"
    availability_zone = data.aws_availability_zones.available.names[0]
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_subnet" "subnet_private_two" {
    vpc_id = aws_vpc.backend_vpc.id
    cidr_block = "10.0.3.0/24"
    availability_zone = data.aws_availability_zones.available.names[1]
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_internet_gateway" "backend_ig" {
    vpc_id = aws_vpc.backend_vpc.id
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_route_table" "public_rt" {
    vpc_id = aws_vpc.backend_vpc.id
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.backend_ig.id
    }
    tags = {
        Application = "langlearn"
        Side = "backend"
    }
}

resource "aws_route_table_association" "public_rta_one" {
    subnet_id = aws_subnet.subnet_public_one.id
    route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_rta_two" {
    subnet_id = aws_subnet.subnet_public_two.id
    route_table_id = aws_route_table.public_rt.id
}