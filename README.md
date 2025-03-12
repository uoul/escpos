# ESC-POS Library for golang
This library contains an implementation of parts of esc-pos protocol. 

## Install
```
go get github.com/uoul/escpos
```

## Usage
```go
package main

import (
	"os"
  "github.com/uoul/escpos"
)

func main() {
	f, err := os.OpenFile("/dev/usb/lp3", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := ns8360l.NewPrinter(f)

	p.Print(
		"Hello World",
		p.WithFontB(),
		p.WithSize(8, 8),
		p.WithUnderline(2),
		p.WithJustifyCenter(),
		// options only effects given text
	)

	p.Print(
		"Hello Mars",
		p.WithFontA(),
		p.WithSize(5, 5),
		p.WithUnderline(2),
		p.WithJustifyLeft(),
		p.WithEmphasize(),
		// options only effects given text
	)

	p.Cut()

}
```