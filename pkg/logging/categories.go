package logging

type Category string

const (
	Startup       Category = "Startup"
	Shutdown      Category = "Shutdown"
	Postgres      Category = "Postgres"
	Storage       Category = "Storage"
	Service       Category = "Service"
	TransportHTTP Category = "Transport/HTTP"
	MQ            Category = "MQ" // Message Queue
	MQWorker      Category = "MQWorker"
	RuntimePanic  Category = "RuntimePanic"
)

func (s Category) String() string {
	return string(s)
}
