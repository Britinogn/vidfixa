package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/britinogn/vidfixa/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

/*
Overview handles GET /api/admin/overview.
The route middleware ensures this is only available to an administrator.
*/
func (h *AdminHandler) Overview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.adminService.GetOverview(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(overview)
}

/*
Users handles GET /api/admin/users.
*/
func (h *AdminHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.adminService.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

/*
Subscriptions handles GET /api/admin/subscriptions.
*/
func (h *AdminHandler) Subscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.adminService.ListSubscriptions(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(subscriptions)
}

/*
Payments handles GET /api/admin/payments.
*/
func (h *AdminHandler) Payments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.adminService.ListPayments(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payments)
}

/*
Downloads handles GET /api/admin/downloads.
*/
func (h *AdminHandler) Downloads(w http.ResponseWriter, r *http.Request) {
	downloads, err := h.adminService.ListDownloads(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(downloads)
}
