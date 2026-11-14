package services

import "evermos-hasby/ecommerce/internal/interfaces"

type HealthCheck struct {
	HealthCheckRepository interfaces.IHealthCheckRepository
}

func (s *HealthCheck) HealthCheckServices() (string, error) {
	return "service is healthy", nil
}
