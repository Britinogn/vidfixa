package service

import (
	"context"

	"gitlab.com/britinogn/vidfixa/internal/repository"
)

/*
AdminService contains read-only operations for the administrator
dashboard. Authorization remains in the route layer, so these methods
only describe the data an already-authorized administrator can view.
*/
type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) GetOverview(ctx context.Context) (*repository.AdminOverview, error) {
	return repository.GetAdminOverview(ctx)
}

func (s *AdminService) ListUsers(ctx context.Context) ([]repository.AdminUser, error) {
	return repository.ListAdminUsers(ctx)
}

func (s *AdminService) ListSubscriptions(ctx context.Context) ([]repository.AdminSubscription, error) {
	return repository.ListAdminSubscriptions(ctx)
}

func (s *AdminService) ListPayments(ctx context.Context) ([]repository.AdminPayment, error) {
	return repository.ListAdminPayments(ctx)
}

func (s *AdminService) ListDownloads(ctx context.Context) ([]repository.AdminDownload, error) {
	return repository.ListAdminDownloads(ctx)
}
