package escpos

type IPrinter interface {
	// Get printer state
	GetPrinterState() (PrinterState, error)
	GetOffLineState() (OffLineState, error)
	GetErrorState() (ErrorState, error)
	GetFeedState() (FeedState, error)
	// Print text
	Print(text string, opts ...func(IPrinter) error) error
	// Print Barcodes
	PrintUpcA(code string, opts ...func(IPrinter) error) error
	PrintUpcE(code string, opts ...func(IPrinter) error) error
	PrintEan13(code string, opts ...func(IPrinter) error) error
	PrintEan8(code string, opts ...func(IPrinter) error) error
	PrintItf(code string, opts ...func(IPrinter) error) error
	PrintCodabar(code string, opts ...func(IPrinter) error) error
	PrintCode39(code string, opts ...func(IPrinter) error) error
	PrintCode93(code string, opts ...func(IPrinter) error) error
	PrintCode128(code string, opts ...func(IPrinter) error) error
	PrintQrCode(code string, ec byte, componentType byte, opts ...func(IPrinter) error) error
	// Options
	WithNegativ() func(IPrinter) error
	WithFontA() func(IPrinter) error
	WithFontB() func(IPrinter) error
	WithUnderline(thickness int) func(IPrinter) error
	WithEmphasize() func(IPrinter) error
	WithRotation() func(IPrinter) error
	WithJustifyLeft() func(IPrinter) error
	WithJustifyCenter() func(IPrinter) error
	WithJustifyRight() func(IPrinter) error
	WithSize(height, width uint8) func(IPrinter) error
	WithLineSpacing(space uint8) func(IPrinter) error
	WithBarcodeHight(hight uint8) func(IPrinter) error
	WithBarcodeWidth(width uint8) func(IPrinter) error
	WithBarcodeStartingPos(pos uint8) func(IPrinter) error
	WithBarcodeHriFontA() func(IPrinter) error
	WithBarcodeHriFontB() func(IPrinter) error
	WithBarcodeNoHri() func(IPrinter) error
	WithBarcodeHriTop() func(IPrinter) error
	WithBarcodeHriBottom() func(IPrinter) error
	WithBarcodeHriTopAndBottom() func(IPrinter) error
	// Cut
	Cut() error
	// Write binary
	WriteRaw(b []byte) error
}

type PrinterState struct {
	DrawerClosed         bool
	Offline              bool
	WaitForOnlineRecover bool
}

type OffLineState struct {
	TopCoverOpen     bool
	FeedByFeedButton bool
	ShortageOfPaper  bool
	Error            bool
}

type ErrorState struct {
	AutoCutterError                            bool
	UnrecoverableError                         bool
	TemperatureAndVoltageOfPrintHeadOutOfRange bool
}

type FeedState struct {
	PaperEnd     bool
	PaperPresent bool
}
