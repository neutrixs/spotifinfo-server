package gettoken

type responseDataType struct {
	Token		string 	`json:"token"`
	ValidUntil 	int		`json:"validuntil"`
}

type successResponseType struct {
	Success 	bool 				`json:"success"`
	Data 		responseDataType	`json:"data"`
}

type failedResponseType struct {
	Success 	bool		`json:"success"`
	ErrorCodes 	[]string	`json:"error-codes"`
}