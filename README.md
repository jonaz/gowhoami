# gowhoami

Simple webserver echoing back the request information. 

### Usage

```
Usage of gowhoami:
  -p string
    	Port to listen on (default "8080")
```


### Example using curl to make the request

```
$ curl localhost:8080
RemoteAddr: [::1]:53042
Host: localhost:8080
Protocol: HTTP/1.1
Headers:
Accept:*/*
User-Agent:curl/7.48.0
```
