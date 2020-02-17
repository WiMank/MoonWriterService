package response

type AppResponse interface {
	PrintLog(err error)
	GetStatusCode() int
}
