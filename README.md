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

	printer := escpos.NewPrinter(f)

	printer.Print(
		"Hello World",
		escpos.WithFontB(),
		escpos.WithSize(8, 8),
		escpos.WithUnderline(2),
		escpos.WithJustifyCenter(),
		// options only effects given text
	)

	printer.Print(
		"Hello Mars",
		escpos.WithFontA(),
		escpos.WithSize(5, 5),
		escpos.WithUnderline(2),
		escpos.WithJustifyLeft(),
		escpos.WithEmphasize(),
		// options only effects given text
	)

	printer.Cut()

}
```