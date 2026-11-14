package interfaces

type IHealthCheckRepository interface {
	HealthCheckRepository() (string, error)
}

type IHealthCheckServices interface {
	HealthCheckServices() (string, error)
}
