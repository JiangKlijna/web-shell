# web-shell
[what is a Web shell?](https://simple.wikipedia.org/wiki/Web_shell)

## Powered by
Web Shell Powered by [gorilla/websocket](https://github.com/gorilla/websocket), [runletapp/go-console](https://github.com/runletapp/go-console) and [xtermjs/xterm.js](https://github.com/xtermjs/xterm.js).
And windows need [rprichard/winpty](https://github.com/rprichard/winpty).

## Installation
### from source code
```bash
git clone github.com/jiangklijna/web-shell
cd web-shell
make gen
make
```
### from release
[releases](https://github.com/JiangKlijna/web-shell/releases)

## Help
```bash
$ web-shell -h
Usage:
  web-shell [-s server mode] [-c client mode]  [-P port] [-u username] [-p password] [-cmd command]

Example:
  web-shell -s -P 2020 -u admin -p admin -cmd bash
  web-shell -c -H 192.168.1.1 -P 2020 -u admin -p admin

Options:
  -C string
        crt file
  -H string
        connect to host (default "127.0.0.1")
  -K string
        key file
  -P string
        listening port (default "2020")
  -RC string
        root crt file
  -c    client mode
  -cmd string
        command cmd or bash
  -cp string
        content path
  -h    this help
  -https
        enable https
  -p string
        password (default "webshell")
  -s    server mode
  -u string
        username (default "admin")
  -v    show version and exit
```

## License
Source code in **web-shell** is available under the [GPL-3.0 License](https://github.com/JiangKlijna/web-shell/blob/master/LICENSE).