# LHA
... is a Go implementation of `ls -lha` with some additional features.

# Installation
```sh
CGO_ENABLED=0 go build .
cp lha /usr/local/bin/
chmod +x /usr/local/bin/lha
```

# Usage
```sh
lha <path>
```
If no path is given the current directory will be inspected.

