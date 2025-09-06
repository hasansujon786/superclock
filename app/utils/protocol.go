package utils

import (
	"encoding/gob"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/models"
)

// sendCmd connects to daemon and sends a command
func sendCmd[T any](cmd constants.Command, payload any) tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("tcp", constants.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)

		req := models.Request{Cmd: cmd, Payload: payload}

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

func DaemonCmd(cmd constants.Command, payload any) tea.Cmd {
	return sendCmd[models.DaemonStateMsg](cmd, payload)
}

func InitDaemonState() any {
	return sendCmd[models.DaemonStateMsg](constants.CmdGet, nil)()
}
