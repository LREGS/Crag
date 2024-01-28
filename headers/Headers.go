package helpers

import (
	h "workspaces/github.com/lregs/Crag/helper"
)

type HttpHeaders interface {
	ReturnHeaders()
}

type MetOfficeHeaders struct {
	CLIENT_ID     string `json:"X-IBM-Client-Id"`
	CLIENT_SECRET string `json:"X-IBM-Client-Secret"`
	Accept        string `json:"accept"`
}

func (m *MetOfficeHeaders) ReturnHeaders() {
	envVariables, err := h.GetEnv([]string{"CLIENT_ID", "CLIENT_SECRET"})
	h.CheckError(err)

	m.CLIENT_ID = envVariables[0]
	m.CLIENT_SECRET = envVariables[1]
	m.Accept = "application/json"

}
