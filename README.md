# ec2-web

This web app should have an endpoint, `http://hostname/builds`, that accepts a POST request with the following sample payload...
  
```
{
    ru {
        "Build base AMI": {
            "Builds": [{
                "runtime_seconds": "1931",
                "build_date": "1506741166",
                "result": "SUCCESS",
                "output": "base-ami us-west-2 ami-9f0ae4e5 d1541c88258ccb3ee565fa1d2322e04cdc5a1fda"
            }, {
                "runtime_seconds": "1825",
                "build_date": "1506740166",
                "result": "SUCCESS",
                "output": "base-ami us-west-2 ami-d3b92a92 3dd2e093fc75f0e903a4fd25240c89dd17c75d66"
            }, {
                "runtime_seconds": "126",
                "build_date": "1506240166",
                "result": "FAILURE",
                "output": "base-ami us-west-2 ami-38a2b9c1 936c7725e69855f3c259c117173782f8c1e42d9a"
            }, {
                "runtime_seconds": "1842",
                "build_date": "1506240566",
                "result": "SUCCESS",
                "output": "base-ami us-west-2 ami-91a42ed5 936c7725e69855f3c259c117173782f8c1e42d9a"
            }, {
                "runtime_seconds": "5",
                "build_date": "1506250561",
            }, {
                "runtime_seconds": "215",
                "build_date": "1506250826",
                "result": "FAILURE",
                "output": "base-ami us-west-2 ami-34a42e15 936c7725e69855f3c259c117173782f8c1e42d9a"
            }]
        }
    }
}
```

  
and should return a json formatted response like this...
 
```
{
    "latest": {
        "build_date": "xxxxxxx",
        "ami_id": "ami-xxxxxx",
        "commit_hash": "xxxxxxxxxxxx"
    }
}
```

  
 
where the build date, ami id, and commit hash is the latest from the the builds in the given payload. Assume that all build dates are epoch timestamps.

  
The infrastructure to support the web app should have:

    1) all resources in us-west-2
    2) 1 LB (any type you feel is suitable)
    3) 1 EC2 instance
    4) Any security groups needed to appropriately restrict access
    5) A key pair to access the instance