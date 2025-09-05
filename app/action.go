package app

import (
	"encoding/gob"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app/constants"
)

type DaemonStateMsg struct {
	Timeout  time.Duration // total countdown time
	Interval time.Duration // tick interval
	Elapsed  time.Duration // time passed
	Running  bool
}

// _sendCmd connects to daemon and sends a command
func _sendCmd[T any](cmd constants.Command, payload any) tea.Cmd {
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

		var msg T
		if err := dec.Decode(&msg); err != nil {
			return err
		}

		return msg
	}
}

func daemonCmd(cmd constants.Command, payload any) tea.Cmd {
	return _sendCmd[DaemonStateMsg](cmd, payload)
}

// func initCmd() tea.Cmd {
// 	return _sendCmd[DaemonStateMsg](constants.CmdGet, nil)
// }

func InitDaemonState() any {
	return _sendCmd[DaemonStateMsg](constants.CmdGet, nil)()
}
