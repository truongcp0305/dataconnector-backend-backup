package utils

const (
	STATUS_OK                = 200
	STATUS_NOT_FOUND         = 404
	STATUS_PERMISSION_DENIED = 403
	STATUS_BAD_REQUEST       = 400
	STATUS_SERVER_ERROR      = 500
)

var (
	STORE_STATUS = map[int]string{
		STATUS_OK:                "OK",
		STATUS_NOT_FOUND:         "Not found",
		STATUS_PERMISSION_DENIED: "Permission denied",
		STATUS_BAD_REQUEST:       "Bad request",
		STATUS_SERVER_ERROR:      "Server error"}
)
