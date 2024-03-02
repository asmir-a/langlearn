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