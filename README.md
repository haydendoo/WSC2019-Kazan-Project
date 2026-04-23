# WSC2019-Kazan-Project
Simulation of test project for the WSC2019 Kazan Test Project

## Initial Infrastructure
There will just be a single EC2 Instance on a VPC with 1 public subnet. Participants are expected to create the rest of the infrastructure however they want

## Server Binary
It is a simple api that has a single endpoint that accepts a query parameter `id` which is an integer. It will check cache first for the corresponding token for this id (both the tokens in redis and in filesystem).

If the cache is missed, it will simulate some expensive work to slow it down before checking the mysql db for the value.

If even the db does not have it, it will create a new entry and enter into mysql and the caches.
