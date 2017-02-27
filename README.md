# turtle-proxy

Allows you to delay requests for a given pattern.

Acts as a reverse proxy. It only supports one proxy target for now.

Note that you can do this on unix by playing around with `tc`. However, if you want something a bit more control this is the tool for you.

## Usage

Get it:

```
go get github.com/lpedrosa/turtle-proxy
```

To get help:

```
$GOPATH/bin/turtle-proxy --help
usage: turtle-proxy [options]

Options:
--host      <0.0.0.0>  Hostname
--port      <6000>     Proxy bind port
--api-port  <6001>     API bind port
--target               Proxy target address (host:port)

Example:
To proxy requests to a server listening on 127.0.0.1:8080
$ turtle-proxy --target 127.0.0.1:8080
```

You just need to specify a target for proxy in the following format: `host:port`.

By default it will listen for incoming connections on port 6000 and api commands on port 6001.

Imagine your api is listening on port 9000 and you want to delay requests going to it.

```
$GOPATH/bin/turtle-proxy -target localhost:9000
```

### Rules

Delay rules have the following json format:

```javascript
{
    "method": <HTTP method>,
    "target": <path pattern, supports pattern matching>,
    "delay": {
        "request": <delay in ms to apply before the request is made to the target>,
        "response": <delay in ms to apply after the request is made to the target>
    }
}
```

### Adding Rules

You can add a rule by issuing a `POST` to `/delay`. For example:

```
curl -i -XPOST \ 
    -H 'Content-Type: application/json' \
    -d '{"method": "GET", "target": "/some-path", "delay": {"response": 5000}}' \
    "http://localhost:6001/delay"
```

This will create a rule that will delay any response from `/some-path` for 5 seconds, but only if the method is `GET`.

### Clearing Rules

You can clear all the added rules by issuing a `DELETE` to `/delay`.

```
curl -i -XDELETE "http://localhost:6001/delay"
```

### Monitoring

You can check if the proxy is alive by issuing a `GET` to `/ping`.

### License

Apache License 2.0

[License](LICENSE)
