package logging

type Category string

const (
	Startup  Category = "Startup"
	Shutdown Category = "Shutdown"
	Postgres Category = "Postgres"
	Service  Category = "Service"
)

func (s Category) String() string {
	return string(s)
}
