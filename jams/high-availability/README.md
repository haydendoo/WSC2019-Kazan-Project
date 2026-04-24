# High Availability (Day 3)
You are given app running on a single EC2 instance with an ALB with a MySQL RDS instance and 2 public and private subnets. Your task is to make it highly available (essentially make everything multi AZ)

1. ASG for EC2 across 2 AZs
2. Mutli AZ RDS deployment
3. 1 NGW in each AZ
4. ALB spans across multiple AZs
