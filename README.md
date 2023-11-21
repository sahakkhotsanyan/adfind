# adfind
## Admin Panel Finder<br />
## Depends golang
## installing adfind
#### sudo git clone https://github.com/sahakkhotsanyan/adfind.git<br />
#### cd adfind*<br />
```bash
go build cmd/adfind/adfind.go -o adfind
```
## Usage
```text
./adfind
Usage of adfind
  -b string
        base path of config files (default is /usr/share/adfind/) (default "/usr/share/adfind/")
  -h    show this help
  -s    stop when admin panel was found
  -t string
        type of admin panel (default is all) {types: php , asp, aspx, js, cfm, cgi, brf. example:adfind -u http://example.com -t php} (default "all")
  -to int
        timeout for request in milliseconds (default 1000)
  -u string
        URL of site {example: adfind -u http://example.com}
  -v    verbose mode
```

## Example
```bash
./adfind -u http://example.com -b ./ -t php
```
