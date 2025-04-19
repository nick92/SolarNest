// Package vedirect implements a VE.Direct protocol parser.
// It mimics the provided Python example by processing incoming bytes
// using a state machine that builds a key/value dictionary representing a frame.
// For more information on the protocol, see:
// https://www.victronenergy.com/live/vedirect_protocol:faq#framehandler_reference_implementation
package sensors

import (
	"time"

	"github.com/tarm/serial"
)

// Define the states used in the state machine.
const (
	HEX         = iota // 0
	WAIT_HEADER        // 1
	IN_KEY             // 2
	IN_VALUE           // 3
	IN_CHECKSUM        // 4
)

// Vedirect represents a VE.Direct protocol parser.
type Vedirect struct {
	Ser       *serial.Port // the serial port
	Header1   byte         // expected to be '\r'
	Header2   byte         // expected to be '\n'
	Hexmarker byte         // expected to be ':'
	Delimiter byte         // expected to be '\t'

	Key      string            // accumulates record key
	Value    string            // accumulates record value
	BytesSum int               // running checksum accumulator
	State    int               // current parser state
	Dict     map[string]string // holds the parsed key/value pairs for one frame
}

// NewVedirect creates a new Vedirect instance and opens the serial port.
// The serial port is opened at 19200 baud with the provided timeout.
func NewVedirect(serialPort string, timeout time.Duration) (*Vedirect, error) {
	c := &serial.Config{
		Name:        serialPort,
		Baud:        19200,
		ReadTimeout: timeout,
	}
	ser, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	v := &Vedirect{
		Ser:       ser,
		Header1:   '\r', // 13
		Header2:   '\n', // 10
		Hexmarker: ':',
		Delimiter: '\t',
		Key:       "",
		Value:     "",
		BytesSum:  0,
		State:     WAIT_HEADER, // initialize to WAIT_HEADER
		Dict:      make(map[string]string),
	}
	return v, nil
}

// Input processes a single byte from the VE.Direct stream and updates the parser state.
// When a complete frame is received (after processing the checksum), it returns the
// built dictionary along with a true flag; otherwise, it returns nil, false.
func (v *Vedirect) Input(b byte) (map[string]string, bool) {
	// If the hex marker (':') is seen and we're not in the checksum state, switch state.
	if b == v.Hexmarker && v.State != IN_CHECKSUM {
		v.State = HEX
	}

	// Add to running checksum for non-HEX states.
	if v.State != HEX {
		v.BytesSum += int(b)
	}

	// Convert any lowercase letter to uppercase.
	if b >= 'a' && b <= 'z' {
		b -= ('a' - 'A')
	}

	switch v.State {
	case WAIT_HEADER:
		// In WAIT_HEADER, we accumulate the checksum.
		if b == v.Header1 {
			// Remain in WAIT_HEADER.
			v.State = WAIT_HEADER
		} else if b == v.Header2 {
			// When a header2 ('\n') is received, switch to collecting the key.
			v.State = IN_KEY
		}
		return nil, false

	case IN_KEY:
		if b == v.Delimiter {
			// End of key; if the key is "Checksum" (case-sensitive),
			// change to checksum processing.
			if v.Key == "CHECKSUM" || v.Key == "Checksum" {
				v.State = IN_CHECKSUM
			} else {
				v.State = IN_VALUE
			}
		} else {
			v.Key += string(b)
		}
		return nil, false

	case IN_VALUE:
		if b == v.Header1 {
			// End-of-record marker: add the key/value pair to the dictionary.
			v.State = WAIT_HEADER
			v.Dict[v.Key] = v.Value
			v.Key = ""
			v.Value = ""
		} else {
			v.Value += string(b)
		}
		return nil, false

	case IN_CHECKSUM:
		// Process the checksum record.
		v.State = WAIT_HEADER
		v.Key = ""
		v.Value = ""
		// If the checksum is valid (bytes sum modulo 256 equals zero), return the frame.
		if v.BytesSum%256 == 0 {
			v.BytesSum = 0
			result := v.Dict
			// Reset the dictionary for the next frame.
			v.Dict = make(map[string]string)
			return result, true
		} else {
			v.BytesSum = 0
		}
		return nil, false

	case HEX:
		// For a HEX record, reset the checksum.
		v.BytesSum = 0
		if b == v.Header2 {
			v.State = WAIT_HEADER
		}
		return nil, false

	default:
		// In an unexpected state.
		panic("invalid state")
	}
}

// ReadDataSingle continuously reads bytes from the serial port until a complete frame is received,
// and then returns the parsed dictionary.
func (v *Vedirect) ReadDataSingle() (map[string]string, error) {
	buf := make([]byte, 1)
	for {
		n, err := v.Ser.Read(buf)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			packet, complete := v.Input(buf[0])
			if complete {
				return packet, nil
			}
		}
	}
}

// ReadDataCallback continuously reads bytes from the serial port and invokes callback with
// the parsed dictionary whenever a complete frame is received.
func (v *Vedirect) ReadDataCallback(callback func(map[string]string)) error {
	buf := make([]byte, 1)
	for {
		n, err := v.Ser.Read(buf)
		if err != nil {
			return err
		}
		if n > 0 {
			packet, complete := v.Input(buf[0])
			if complete {
				callback(packet)
			}
		}
	}
}
