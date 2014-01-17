## About

This is a port of JSMin (by Douglas Crockford) to the Go programming language
(http://golang.org/).

## Installing

### Using *go get*

    $ go get github.com/ae0000/gomin


## Example

    import (
            "github.com/ae0000/gomin"
    )

    func main() {

        rawJs := []byte("var abc = 123")
        minifiedJs, err := gomin.Js(rawJs)

        if err != nil {
            // Try harder
        }
    }
