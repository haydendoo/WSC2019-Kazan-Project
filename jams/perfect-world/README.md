# Perfect World (Day 3)
You are given a lambda function, SQS queue, EC2 instance. You must configure it such that when someone logins via ssh, a script is triggered which sends a message to SQS, then lambda processes it, then it blocks ssh for the ec2 instance (via security group modification).

Your task is to modify sqs configuration and lambda code to block it properly. EC2 does not need to be touched except for testing by sshing in.
