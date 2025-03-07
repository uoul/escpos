package escpos

import (
	"fmt"
	"io"
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

// GetErrorState implements IPrinter.
func (n *Printer) GetErrorState() (ErrorState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 3})
	if err != nil {
		return ErrorState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return ErrorState{}, err
	}
	return ErrorState{
		AutoCutterError:    (b[0] & 0x08) != 0,
		UnrecoverableError: (b[0] & 0x20) != 0,
		TemperatureAndVoltageOfPrintHeadOutOfRange: (b[0] & 0x40) != 0,
	}, nil
}

// GetFeedState implements IPrinter.
func (n *Printer) GetFeedState() (FeedState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return FeedState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return FeedState{}, err
	}
	return FeedState{
		PaperEnd:     (b[0] & 0x0C) != 0,
		PaperPresent: (b[0] & 0x60) != 0,
	}, nil
}

// GetOffLineState implements IPrinter.
func (n *Printer) GetOffLineState() (OffLineState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return OffLineState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return OffLineState{}, err
	}
	return OffLineState{
		TopCoverOpen:     (b[0] & 0x04) != 0,
		FeedByFeedButton: (b[0] & 0x08) != 0,
		ShortageOfPaper:  (b[0] & 0x20) != 0,
		Error:            (b[0] & 0x40) != 0,
	}, nil
}

// GetPrinterState implements IPrinter.
func (n *Printer) GetPrinterState() (PrinterState, error) {
	err := n.WriteRaw([]byte{_DLE, _EOT, 4})
	if err != nil {
		return PrinterState{}, err
	}
	b := make([]byte, 1)
	_, err = n.rw.Read(b)
	if err != nil {
		return PrinterState{}, err
	}
	return PrinterState{
		DrawerClosed:         (b[0] & 0x04) != 0,
		Offline:              (b[0] & 0x08) != 0,
		WaitForOnlineRecover: (b[0] & 0x20) != 0,
	}, nil
}

// WriteRaw implements IPrinter.
func (n *Printer) WriteRaw(b []byte) error {
	_, err := n.rw.Write(b)
	return err
}

// Cut implements IPrinter.
func (n *Printer) Cut() error {
	if err := n.WriteRaw([]byte{_GS, 'V', 66, 30}); err != nil {
		return err
	}
	n.Print("")
	return nil
}

// Print implements IPrinter.
func (n *Printer) Print(text string, opts ...func(IPrinter) error) error {
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

// PrintCodabar implements IPrinter.
func (n *Printer) PrintCodabar(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(73, code, 36, 68, 1, 255, opts...)
}

// PrintCode128 implements IPrinter.
func (n *Printer) PrintCode128(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(73, code, 0, 127, 2, 255, opts...)
}

// PrintCode39 implements IPrinter.
func (n *Printer) PrintCode39(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(69, code, 32, 90, 1, 255, opts...)
}

// PrintCode93 implements IPrinter.
func (n *Printer) PrintCode93(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(72, code, 0, 127, 1, 255, opts...)
}

// PrintEan13 implements IPrinter.
func (n *Printer) PrintEan13(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(67, code, '0', '9', 12, 13, opts...)
}

// PrintEan8 implements IPrinter.
func (n *Printer) PrintEan8(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(68, code, '0', '9', 7, 8, opts...)
}

// PrintItf implements IPrinter.
func (n *Printer) PrintItf(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(70, code, '0', '9', 1, 255, opts...)
}

// PrintUpcA implements IPrinter.
func (n *Printer) PrintUpcA(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(65, code, '0', '9', 11, 12, opts...)
}

// PrintUpcE implements IPrinter.
func (n *Printer) PrintUpcE(code string, opts ...func(IPrinter) error) error {
	return n.printBarcode(66, code, '0', '9', 11, 12, opts...)
}

// --------------------------------------------------------------------------------
// Options
// --------------------------------------------------------------------------------

func WithNegativ() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'B', 1})
	}
}

func WithFontA() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'M', '0'})
	}
}

func WithFontB() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'M', '1'})
	}
}

func WithUnderline(thickness int) func(IPrinter) error {
	return func(i IPrinter) error {
		if thickness < 0 || thickness > 2 {
			return fmt.Errorf("underline thickness has to be between 0 and 2")
		}
		return i.WriteRaw([]byte{_ESC, '-', byte(thickness)})
	}
}

func WithEmphasize() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'E', 1})
	}
}

func WithRotation() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'V', '1'})
	}
}

func WithJustifyLeft() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '0'})
	}
}

func WithJustifyCenter() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '1'})
	}
}

func WithJustifyRight() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 'a', '2'})
	}
}

func WithSize(height, width uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, '!', (((width - 1) << 4) | (height - 1))})
	}
}

func WithLineSpacing(space uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, '3', space})
	}
}

func WithBarcodeHight(hight uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'h', hight})
	}
}

func WithBarcodeWidth(width uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'w', width})
	}
}

func WithBarcodeStartingPos(pos uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'x', pos})
	}
}

func WithBarcodeHriFontA() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'f', '0'})
	}
}

func WithBarcodeHriFontB() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'f', '1'})
	}
}

func WithBarcodeNoHri() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '0'})
	}
}

func WithBarcodeHriTop() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '1'})
	}
}

func WithBarcodeHriBottom() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, 'H', '2'})
	}
}

func WithBarcodeHriTopAndBottom() func(IPrinter) error {
	return func(i IPrinter) error {
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

func (n *Printer) printBarcode(codeType uint8, code string, minChar, maxChar rune, minLen, maxLen int, opts ...func(IPrinter) error) error {
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

func NewPrinter(rw io.ReadWriter) IPrinter {
	return &Printer{
		rw: rw,
	}
}
