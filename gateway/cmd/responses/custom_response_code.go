package responses

var CustomResponseCodes = map[string]int{
	"create_successful": 1000,
	"fetch_successful":  1001,
	"update_successful": 1002,
	"delete_successful": 1003,
	"custom_error":      1004,
	"validation_error":  1005,
}
