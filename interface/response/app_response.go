package response

type AppResponseInterface interface {
	PrintLog(err error)
	GetStatusCode() int
}

type AppResponseCreator interface {
	CreateResponse() AppResponseInterface
}
