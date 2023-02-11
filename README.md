# Check TLS Versions

`tlsversion` is a command line tools allowing to display the TLS versions supported by one or more hosts.
It only supports TLS 1.0 through TLS 1.3.

![image](https://github.com/andreburgaud/tlsversion/releases/download/0.3.0/Screenshot.from.2023-02-10.22-31-03.png)

## Usage

```bash
$ tlsversion google.com
Host        TLS1.0  TLS1.1  TLS1.2  TLS1.3  Error  
google.com  Y       Y       Y       Y       -     
```

`tlsversion` can take one or more hosts at the command line. The port defaults to 443,
but you can override the default port by adding `:some_port` to the host. For example:

```bash
$ tlsversion google.com:443
Host        TLS1.0  TLS1.1  TLS1.2  TLS1.3  Error  
google.com  Y       Y       Y       Y       -     
```

`tlsversion` can also take a file as argument with the option `--file`. The file should include
one host per line. At the beginning of a line, a `#` prefix denotes a comment.

```bash
$ cat data/servers.txt
burgaud.com
google.com
microsoft.com
cloudflare.com
# Server supporting only TLS1.0
tls-v1-0.badssl.com:1010
# Server supporting only TLS1.1
tls-v1-1.badssl.com:1011
# Server supporting only TLS1.2
tls-v1-2.badssl.com:1012
www.mozilla.org
# The following host does not exist
www.serverdoesnot.exist
```

```
$ tlsversion --file data/servers.txt
Host                      TLS1.0  TLS1.1  TLS1.2  TLS1.3  Error                                                                              
burgaud.com               Y       Y       Y       Y       -                                                                                  
cloudflare.com            Y       Y       Y       Y       -                                                                                  
google.com                Y       Y       Y       Y       -                                                                                  
microsoft.com             N       N       Y       Y       -                                                                                  
tls-v1-0.badssl.com:1010  Y       N       N       N       -                                                                                  
tls-v1-1.badssl.com:1011  -       Y       N       N       read tcp 192.168.86.31:49978->104.154.89.105:1011: read: ...
tls-v1-2.badssl.com:1012  -       -       Y       N       read tcp 192.168.86.31:43370->104.154.89.105:1012: read: ...  
www.mozilla.org           Y       Y       Y       Y       -                                                                                  
www.serverdoesnot.exist   -       -       -       -       dial tcp: lookup www.serverdoesnot.exist on 127.0.0.53:53... 
```

## Build

The build is managed with [justfile](https://github.com/casey/just) and [goreleaser](https://goreleaser.com/).
Provided that you have [Go](https://go.dev/) installed on your machine, you can build `tlsversion` 
with the following commands:

```
$ go build tlsversion/cmd/tlsversion                   # Debug build
...
$ go build -ldflags="-s -w" tlsversion/cmd/tlsversion  # Relase build
...
```

## Licenses

`tlsversion` is released under the MIT license.

It uses the following libraries also available under MIT licenses:
* https://github.com/fatih/color
* https://github.com/rodaine/table

