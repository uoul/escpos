package escpos

type IPrinter interface {
	GetPrinterState() (PrinterState, error)
	GetOffLineState() (OffLineState, error)
	GetErrorState() (ErrorState, error)
	GetFeedState() (FeedState, error)

	Print(text string, opts ...func(IPrinter) error) error
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
