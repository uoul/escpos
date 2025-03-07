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

	// Cut
	Cut() error

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
