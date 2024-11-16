# Brief
This is a simple plugin of type HTTP Server Plugin (see here:https://www.krakend.io/docs/extending/).  
The purpose of this plugin is to demonstrate  how to capture a specific request path
and do some action on it.  For the purpose of this demo, we look for a particular URL path
(configured in krakend.json), and if any request comes for that path, then we 
write a response to the client ourselves (see lines #58 in main.go).  For all other requests
we just forward them to be handled by KrakenD.

<br/>

# How to configure the plugin 

1. Prepare your krakend.json file.
It is essential to create your krakend.json file with additional
configuration apart from your regular routing info.  The krakenD
service loads this json during startup, and hence comes to know
of the plugin details and loads it too.

The crucial part of the krakend.json is the `plugin` configuration
which tells about the location of the plugin.  After that, the 
`extra_config` portion is meant for the plugin itself.

```json
  "plugin": {
    "pattern": ".so",
    "folder": "/opt/krakend/plugins/"
  }, 
  "extra_config": {
    "plugin/http-server": {
      "name": ["krakend-plugin"],
      "krakend-plugin": {
        "path": "/_hijack"
      }
    }
  }

```
In the example above, we are telling the plugin to intercept the 
http request having a path `/_hijack` and do something with this
request, before it is handled by the KrakenD.  In our example
implementation, the plugin sends back a reponse `Lowjack "/_hijack"` to the API client.

2. Assuming that you have created the plugin code as per the 
official docuemntation, you can use the Dockerfile present in this
repo to carry out the building of the plugin and including it as 
part of the KrakenD image.  Using this Dockerfile ensures that 
the environment (go version and architecture) needed for compiling
the plugin wwill be the same as the one used for building the 
KrakenD itself.

of go lang that is compatible with krakenD (at the time of writing
this document, it is 1.22.7).

In order to build, use the below command. Replace the name 'my_krakend' with any other name of your liking. 

```console
docker build -t my_krakend .
```

# How to use
Run the docker image to bring up krakend instance.
```console 
docker run --name myKraken -d -p 8080:8080 mykrak 

```

## Test the plugin
```console 
$ curl -i localhost:8080/_hijack
HTTP/1.1 200 OK
Date: Sat, 16 Nov 2024 11:17:06 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8

Lowjack "/_hijack"%

```


## Happy coding.
