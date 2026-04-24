# Identify and mitigate configuration drift! (Day 4)
You are given an app running on 2 ec2 instances between ALB with a lambda function to check for answer and an aws config rule.

One of the ec2 instances is running an old version of the application and you must upgrade it using automation.

You must use ssm to build a software inventory, find the vulnerable instance, use ssm to automatically upgrade. (AWS Config will mark it as compliant once done and lambda function can be invoked to get the answer)
