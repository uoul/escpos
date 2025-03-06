package escpos

import (
	"fmt"
	"io"
)

// --------------------------------------------------------------------------------
// Constants
// --------------------------------------------------------------------------------
const (
	_ESC  = 0x1b
	_GS   = 0x1d
	_LF   = 0x0a
	_DLE  = 0x10
	_EOT  = 0x04
	_AT   = 0x40
	_DASH = 0x2d
	_EXC  = 0x21

	_J = 0x4a
	_G = 0x47
	_E = 0x45
	_B = 0x42
	_M = 0x4d
	_V = 0x56
	_a = 0x61
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
	return n.WriteRaw([]byte{_GS, 'V', 66, 30})
}

// Print implements IPrinter.
func (n *Printer) Print(text string, opts ...func(IPrinter) error) error {
	err := n.WriteRaw([]byte{_ESC, _AT})
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

// --------------------------------------------------------------------------------
// Options
// --------------------------------------------------------------------------------

func WithNegativ() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, _B, 0x01})
	}
}

func WithFontA() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _M, 0})
	}
}

func WithFontB() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _M, 1})
	}
}

func WithUnderline(thickness int) func(IPrinter) error {
	return func(i IPrinter) error {
		if thickness < 0 || thickness > 2 {
			return fmt.Errorf("underline thickness has to be between 0 and 2")
		}
		return i.WriteRaw([]byte{_ESC, _DASH, byte(thickness)})
	}
}

func WithEmphasize() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _E, 0x01})
	}
}

func WithRotation() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _V, 0x49})
	}
}

func WithJustifyLeft() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _a, 0x00})
	}
}

func WithJustifyCenter() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _a, 0x01})
	}
}

func WithJustifyRight() func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, _a, 0x02})
	}
}

func WithSize(height, width uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_GS, _EXC, (((width - 1) << 4) | (height - 1))})
	}
}

func WithLineSpacing(space uint8) func(IPrinter) error {
	return func(i IPrinter) error {
		return i.WriteRaw([]byte{_ESC, 0x33, space})
	}
}

// --------------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------------

// --------------------------------------------------------------------------------
// Constructor
// --------------------------------------------------------------------------------

func NewPrinter(rw io.ReadWriter) IPrinter {
	return &Printer{
		rw: rw,
	}
}
