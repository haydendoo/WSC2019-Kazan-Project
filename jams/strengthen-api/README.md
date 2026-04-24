# Strengthen your API Defense!! (Day 4)
You currently have an api that was leaked and is insecure. 

Your task is to restrict access to it first so it is only callable from same VPC via NAT gateway. (Configure access control using api gateway)

You should also create a new api gateway (private, only accessible from private VPC). It must use VPC endpoints to get access to services (cannot just expose to pubilc). Only accepts new api keys, old ones must be rejected.
