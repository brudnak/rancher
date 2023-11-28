# Log Configs

## Table of Contents
1. [Getting Started](#Getting-Started)

## Getting Started
Your GO suite should be set to `-run ^TestLogsTestSuite$`. You can find any additional suite name(s) by checking the test file you plan to run.

In your config file, set the following:
```yaml
rancher: 
  host: "rancher_server_address"
  adminToken: "rancher_admin_token"
  ...
awsCredentials:
  accessKey: "your_access_key"
  secretKey: "your_secret_key"
```
