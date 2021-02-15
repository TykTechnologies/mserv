# Mserv Tyk plugin server

Mserv is an asset server and gRPC host for the [Tyk open source API Gateway](https://tyk.io). The Tyk API gateway can be extended using middleware plugins, these extensions can run in the same process as the gateway when written in JavaScript, Python or Lua, and can be added as "nanoservices" using gRPC-enabled languages.

When using in-process plugins, plugin "bundles" (cryptographically signed zip files including all the necessary code to be run), these files must be stored in a third-party asset server (such as S3) and Tyk looks for them on a base URL.

In gRPC mode, the gateway still requires a downloadable bundle, but no code, instead the gRPC service is run as a side-car.

Tyk Mserv provides two services to the Tyk Gateway user:

1. Acts as an Asset server for plugins: use the Mserv REST API to push plugin bundles to a secure back end (currently only S3 and Local filesystem are supported). Tyk can then be pointed at MServ to retrieve bundle files. instead of wiring directly into S3.
2. Act as a middleware gRPC server for golang-based gRPC plugins: Write middleware plugins aa simple go functions and compile them as golang modules (`.so` files), these can then be pushed live into MServ to provide a dynamic, out-of-the-box middleware server.

## Pre-requisites

1. A Golang environment (1.8 and over, 1.10 is recommended)
2. Tyk-cli

To get Tyk CLI:

```
go get -u github.com/TykTechnologies/tyk-cli
```

## Pushing live gRPC plugins:

In order to build plugins for Mserv, they need ot be built against the Mserv code base, otherwise there are type differences which completely break the plugin import.

The best way to do this is to ensure that both mserv and the plugin are built against the same version. To do that:

```
go get TykTechnologies/mserv
cd $GOPATH/TykTechnologies/mserv/build
```

Then create your plugin in this directory, for example, assume we want a "pre-auth" middleware hook, create a file called `prehook.go`:

```
package main

import (
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
)

// MyPreHook performs a header injection:
func MyPreHook(object *coprocess.Object) (*coprocess.Object, error) {
	object.Request.SetHeaders = map[string]string{
		"Myheader": "Myvalue",
	}

	return object, nil
}
```

Then once you have created the plugin code, save it and compile it as a linux symbol file:

    go build -buildmode=plugin -o plugin.so prehook.go

Now create your manifest file (`manifest.json`):

```
{
  "file_list": [
    "plugin.so"
  ],
  "custom_middleware": {
    "pre": [
      {
        "name": "MyPreHook"
      }
    ],
    "driver": "grpc"
  }
}
```

Now you can create an unsigned bundle using:

    tyk-cli bundle build --output bundle.zip -y

You should now see a `bundle.zip` file. You can puh this to MServ using curl:

    curl -F "uploadfile=@bundle.zip" "http://localhost:8989/api/mw?api_id={API-ID}"

The `api_id` parameter is important here, Mserv will namespace all plugin functions against their calling API ID.

This command will return a bundle ID, you don't need this to invoke the MW, but you will need it to download the bundle file:

    curl "http://localhost:8989/api/mw/{bundle-id}/bundle.zip" --output bundle.zip

The above API command will fetch the bundle zip file.


## Testing plugins with HTTP calls:

If you have the option enabled, then HTTP invocation is possible using `curl`:


curl -X POST \
  http://localhost:8989/execute/{hook-name} \
  -H 'Content-Type: application/json' \
  -d '{
	"hook_type": 1,
	"hook_name": "{hook-name}",
	"request": {
		"headers": {
			"foo": "bar"
		},
		"body": "foo",
		"url": "http://localhost:8989/baz",
		"method": "POST"
	},
	"spec": {
		"APIID": "12345",
		"OrgID": ""
	}
}'

The above JSON is actually a [`coprocess object`](https://github.com/TykTechnologies/tyk/blob/master/coprocess/coprocess_object.pb.go#L20) in JSON format, the built-in `execute` endpoint just feeds this into the dispatcher and invokes your function. Since we are testing middleware, you will receive back a "MiniRequestObject" which shows what the plugin is doing - e.g. adding headers or modifying the body.

You can see something similar, but with a "real" http request using the client application in the `client` directory (still WIP).



