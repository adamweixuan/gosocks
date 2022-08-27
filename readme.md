# usage 

```shell
socks5 ツ ./bin/gosocks-darwin-amd64-v3 --help
Usage of ./bin/gosocks-darwin-amd64-v3:
  -heap int
    	set memory size limit.  (default 9223372036854775807)
  -iface string
    	set specified interface to use.
  -network string
    	set network tcp or udp. (default "tcp")
  -nodelay
    	enable tcpnodelay.  (default true)
  -port uint
    	set socks server listen port.  (default 10086)
  -reuseaddr
    	enable reuseaddr.  (default true)
  -timeout duration
    	set session timeout.  (default 1s)
  -verbos
    	enable verbose log.
```

# support

- [x] ipv4 
- [ ] ipv6
- [x] tcp 
- [ ] udp 
- [ ] auth 
- [ ] gssapi 