# WSC2019-Kazan-Project
Simulation of test project for the WSC2019 Kazan Test Project. It has a gameday on day 1 & 2, but it is combined into just 1 here as they are extremely similar. Only difference is day 1 requires efs and uses MySQL + redis, day 2 requires WAF and uses Postgresql + Memcached and has 2 server binaries.

## Initial Infrastructure
There will just be a single EC2 Instance on a VPC with 1 public subnet. Participants are expected to create the rest of the infrastructure however they want

## Server Binary
It is a simple api that has a single endpoint that accepts a query parameter `id` which is an integer. It will check cache first for the corresponding token for this id (both the tokens in redis and in filesystem).

If the cache is missed, it will simulate some expensive work to slow it down before checking the mysql db for the value.

If even the db does not have it, it will create a new entry and enter into mysql and the caches.

## Common issues when doing gameday
1. Can't mount EFS to EC2: Ensure gameday VPC has dns resolution and hostname turned on
2. Some instances do not have access to internet: Ensure pubilc subnet 2 is associated to the route table for IGW
3. Application logs not appearing in cloudwatch: EC2 instances have NAT gateway to access internet and has SG that allows outbound HTTPS and has iam role with the cloudwatch server policy.

## Final Infrastructure
Roughly there will be the following:
1. A VPC with multiple subnets and atleast 2 AZs used for availability (different subnets for different applications/services)
2. Proper IGW and NGW setup for public and private subnets
3. Proper SG configuration according to least privilege
4. EC2 instances are in ASG and can scale horizontally with an ALB to distribute load
5. RDS is used to host mysql/postgresql
6. Elasticache is used to host redis/memcached
7. All EC2 instances in ASG can access the shared EFS filesystem
8. A WAF is used to stop SQL Injection attacks
9. (Optional) A cloudfront distribution infront of the ALB (reject requests going into alb directly, only allow from cloudfront)
