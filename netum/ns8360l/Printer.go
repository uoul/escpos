package ns8360l

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/uoul/escpos"
)

// --------------------------------------------------------------------------------
// Constants
// --------------------------------------------------------------------------------
const (
	_ESC = 0x1b
	_GS  = 0x1d
	_DLE = 0x10
	_EOT = 0x04
)

// --------------------------------------------------------------------------------
// Types
// --------------------------------------------------------------------------------
type Printer struct {
	rw io.ReadWriter
}

// --------------------------------------------------------------------------------
// Public
// --------------------------------------------------------------------------------

// GetErrorState implements escpos.IPrinter.
func (n *Printer) GetErrorState() (escpos.ErrorState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 3})
	if err != nil {
		return escpos.ErrorState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return escpos.ErrorState{}, err
	}
	return escpos.ErrorState{
		AutoCutterError:    (b[0] & 0x08) != 0,
		UnrecoverableError: (b[0] & 0x20) != 0,
		TemperatureAndVoltageOfPrintHeadOutOfRange: (b[0] & 0x40) != 0,
	}, nil
}

// GetFeedState implements escpos.IPrinter.
func (n *Printer) GetFeedState() (escpos.FeedState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return escpos.FeedState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return escpos.FeedState{}, err
	}
	return escpos.FeedState{
		PaperEnd:     (b[0] & 0x0C) != 0,
		PaperPresent: (b[0] & 0x60) != 0,
	}, nil
}

// GetOffLineState implements escpos.IPrinter.
func (n *Printer) GetOffLineState() (escpos.OffLineState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return escpos.OffLineState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return escpos.OffLineState{}, err
	}
	return escpos.OffLineState{
		TopCoverOpen:     (b[0] & 0x04) != 0,
		FeedByFeedButton: (b[0] & 0x08) != 0,
		ShortageOfPaper:  (b[0] & 0x20) != 0,
		Error:            (b[0] & 0x40) != 0,
	}, nil
}

// GetPrinterState implements escpos.IPrinter.
func (n *Printer) GetPrinterState() (escpos.PrinterState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return escpos.PrinterState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return escpos.PrinterState{}, err
	}
	return escpos.PrinterState{
		DrawerClosed:         (b[0] & 0x04) != 0,
		Offline:              (b[0] & 0x08) != 0,
		WaitForOnlineRecover: (b[0] & 0x20) != 0,
	}, nil
}

// WriteRaw implements escpos.IPrinter.
func (n *Printer) WriteRaw(b []byte) error {
	_, err := n.rw.Write(b)
	return err
}

// Cut implements escpos.IPrinter.
func (n *Printer) Cut() error {
	if err := n.WriteRaw([]byte{_GS, 'V', 66, 30}); err != nil {
		return err
	}
	n.Print("")
	return nil
}

// Print implements escpos.IPrinter.
func (n *Printer) Print(text string, opts ...func(escpos.IPrinter) error) error {
	err := n.WriteRaw([]byte{_ESC, '@'})
	if err != nil {
		return err
	}
	for _, o := range opts {
		if err := o(n); err != nil {
			return err
		}
	}
	return n.WriteRaw([]byte(text))
}

// PrintCodabar implements escpos.IPrinter.
func (n *Printer) PrintCodabar(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(73, code, 36, 68, 1, 255, opts...)
}

// PrintCode128 implements escpos.IPrinter.
func (n *Printer) PrintCode128(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(73, code, 0, 127, 2, 255, opts...)
}

// PrintCode39 implements escpos.IPrinter.
func (n *Printer) PrintCode39(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(69, code, 32, 90, 1, 255, opts...)
}

// PrintCode93 implements escpos.IPrinter.
func (n *Printer) PrintCode93(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(72, code, 0, 127, 1, 255, opts...)
}

// PrintEan13 implements escpos.IPrinter.
func (n *Printer) PrintEan13(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(67, code, '0', '9', 12, 13, opts...)
}

// PrintEan8 implements escpos.IPrinter.
func (n *Printer) PrintEan8(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(68, code, '0', '9', 7, 8, opts...)
}

// PrintItf implements escpos.IPrinter.
func (n *Printer) PrintItf(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(70, code, '0', '9', 1, 255, opts...)
}

// PrintUpcA implements escpos.IPrinter.
func (n *Printer) PrintUpcA(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(65, code, '0', '9', 11, 12, opts...)
}

// PrintUpcE implements escpos.IPrinter.
func (n *Printer) PrintUpcE(code string, opts ...func(escpos.IPrinter) error) error {
	return n.printBarcode(66, code, '0', '9', 11, 12, opts...)
}

// PrintQrCode implements escpos.IPrinter.
func (n *Printer) PrintQrCode(code string, ec byte, componentType byte, opts ...func(escpos.IPrinter) error) error {
	err := n.WriteRaw([]byte{_ESC, '@'})
	if err != nil {
		return err
	}
	for _, o := range opts {
		if err := o(n); err != nil {
			return err
		}
	}
	codeLen := uint16(len(code))
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, codeLen)
	return n.WriteRaw(append([]byte{_ESC, 'Z', '0', ec, componentType, l[0], l[1]}, []byte(code)...))
}

// --------------------------------------------------------------------------------
// Options
// --------------------------------------------------------------------------------

func WithNegativ() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'B', 1})
	}
}

func WithFontA() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'M', '0'})
	}
}

func WithFontB() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'M', '1'})
	}
}

func WithUnderline(thickness int) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		if thickness < 0 || thickness > 2 {
			return fmt.Errorf("underline thickness has to be between 0 and 2")
		}
		return i.WriteRaw([]byte{_ESC, '-', byte(thickness)})
	}
}

func WithEmphasize() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'E', 1})
	}
}

func WithRotation() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'V', '1'})
	}
}

func WithJustifyLeft() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '0'})
	}
}

func WithJustifyCenter() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '1'})
	}
}

func WithJustifyRight() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '2'})
	}
}

func WithSize(height, width uint8) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, '!', (((width - 1) << 4) | (height - 1))})
	}
}

func WithLineSpacing(space uint8) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_ESC, '3', space})
	}
}

func WithBarcodeHight(hight uint8) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'h', hight})
	}
}

func WithBarcodeWidth(width uint8) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'w', width})
	}
}

func WithBarcodeStartingPos(pos uint8) func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'x', pos})
	}
}

func WithBarcodeHriFontA() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'f', '0'})
	}
}

func WithBarcodeHriFontB() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'f', '1'})
	}
}

func WithBarcodeNoHri() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '0'})
	}
}

func WithBarcodeHriTop() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '1'})
	}
}

func WithBarcodeHriBottom() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '2'})
	}
}

func WithBarcodeHriTopAndBottom() func(escpos.IPrinter) error {
	return func(i escpos.IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '3'})
	}
}

// --------------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------------
func checkCharRange(min, max rune, code string) error {
	for _, c := range code {
		if c < min || c > max {
			return fmt.Errorf("char(%v) in code is out of range(%v, %v)", c, min, max)
		}
	}
	return nil
}

func (n *Printer) printBarcode(codeType uint8, code string, minChar, maxChar rune, minLen, maxLen int, opts ...func(escpos.IPrinter) error) error {
	l := len(code)
	if l < minLen || l > maxLen {
		return fmt.Errorf("code length(%d) not in range (11 <= len <= 12)", l)
	}
	if err := checkCharRange(minChar, maxChar, code); err != nil {
		return err
	}
	err := n.WriteRaw([]byte{_ESC, '@'})
	if err != nil {
		return err
	}
	for _, o := range opts {
		if err := o(n); err != nil {
			return err
		}
	}
	return n.WriteRaw(append([]byte{_GS, 'k', codeType, uint8(l)}, []byte(code)...))
}

// --------------------------------------------------------------------------------
// Constructor
// --------------------------------------------------------------------------------

func NewPrinter(rw io.ReadWriter) escpos.IPrinter {
	return &Printer{
		rw: rw,
	}
}
