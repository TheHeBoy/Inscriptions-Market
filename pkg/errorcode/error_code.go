package errorcode

import "fmt"

type ErrorCode struct {
	Code int    `json:"code"`
	Msg  string `json:"error"`
}

func (errorCode *ErrorCode) Error() string {
	return fmt.Sprint("code:", errorCode.Code, " msg:", errorCode.Msg)
}
