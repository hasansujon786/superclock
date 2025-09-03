package timer

import (
	"encoding/gob"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app/constants"
)

// sendCmd connects to daemon and sends a command
func sendCmd(cmd constants.Command, payload any) tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("tcp", constants.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)

		req := constants.Request{Cmd: cmd, Payload: payload}

		if err := enc.Encode(req); err != nil {
			return err
		}

		var timer TimerState
		if err := dec.Decode(&timer); err != nil {
			return err
		}

		return timer
	}
}
