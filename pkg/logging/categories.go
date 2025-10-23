package logging

type Category string

const (
	Startup       Category = "Startup"
	Shutdown      Category = "Shutdown"
	Postgres      Category = "Postgres"
	Storage       Category = "Storage"
	Service       Category = "Service"
	TransportHTTP Category = "Transport/HTTP"
	RuntimePanic  Category = "RuntimePanic"
)

func (s Category) String() string {
	return string(s)
}
