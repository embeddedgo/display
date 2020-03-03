package eve

// DCI defines Data and Control Interface to the EVE Display Controller.
type DCI interface {
	// SetClk sets the communication interface clock speed. If clkHz = 0 the
	// communication interface can be disabled to save power.
	SetClk(clkHz int)

	// Read reads len(p) bytes into p.
	Read(p []byte)

	// Write writes len(p) bytes from p.
	Write(p []byte)

	// End finishes the current read/write transaction.
	End()

	// Err returns an error encountered while executing the last Read, Write,
	// End command. DCI should stop executing commands until the error will be
	// cleared.
	Err(clear bool) error

	// SetPDN sets the PD_N pin to the least significant bit of pdn.
	SetPDN(pdn int)
}
