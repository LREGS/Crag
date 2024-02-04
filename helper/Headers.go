package helpers

// type HttpHeaders interface {
// 	ReturnHeaders()
// }

// type MetOfficeHeaders struct {
// 	CLIENT_ID     string `json:"X-IBM-Client-Id"`
// 	CLIENT_SECRET string `json:"X-IBM-Client-Secret"`
// 	Accept        string `json:"accept"`
// }

// simplest way now but maybe we want to make other requests or expand headers and would
// want structs/interfaces I dont know
func ReturnHeaders() map[string]string {
	envVariables, err := GetEnv([]string{"CLIENT_ID", "CLIENT_SECRET"})
	CheckError(err)

	// m.CLIENT_ID = envVariables[0]
	// m.CLIENT_SECRET = envVariables[1]
	// m.Accept = "application/json"

	return map[string]string{
		"X-IBM-Client-Id":     envVariables[0],
		"X-IBM-Client-Secret": envVariables[1],
		"accept":              "application/json",
	}

}