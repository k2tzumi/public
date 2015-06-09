# Groupcache DB Experiment - Revisited
This is the revisited version of Capotej's [Groupcache DB Experiment](https://github.com/capotej/groupcache-db-experiment).

I decided to replay his experiment using newer techniques, such the use of [gb](http://getgb.io), `init()` niladic functions and fail soon error handling.

This project simulates a scenario wherein a few frontends running [groupcache](http://github.com/golang/groupcache) are fronting a slow database. See his [blog post](http://capotej.com/blog/2013/07/28/playing-with-groupcache/) about it for more details.

# Getting it running
The following commands will set up this topology:
![groupcache topology](https://raw.github.com/capotej/groupcache-db-experiment/master/topology.png)

### Build everything

1. ```git clone https://github.com/dericofilho/groupcache-experiment.git```
2. ```go get github.com/constabulary/gb/...```
3. ```gb build all```

### Start DB server

1. ```./bin/backend```

This starts a deliberately slow k/v datastore on :8080

### Start Multiple Frontends

1. ```./bin/frontend -listen "http://localhost:8001" -frontend "localhost:9001" ```
2. ```./bin/frontend -listen "http://localhost:8002" -frontend "localhost:9002" ```
3. ```./bin/frontend -listen "http://localhost:8003" -frontend "localhost:9003" ```

### Use the CLI to set/get values

1. ```./cli -set -k foo -v bar```
2. ```./cli -get -k foo``` should see bar in 300 ms
3. ```./cli -get -k foo``` should see bar instantly
