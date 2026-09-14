like [proxytunnel](https://github.com/proxytunnel/proxytunnel), but uses HTTP/3 instead

## example

```sh
ssh -o ProxyCommand='proxytunnel-http3 -p proxy.example.org:443 -d %h:%p' examplehost
```