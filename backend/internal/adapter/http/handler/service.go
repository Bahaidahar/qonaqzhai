package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/vendor"
)

// Service handles vendor service (menu item) endpoints.
type Service struct {
	svc *vendor.Service
}

// NewService constructs a Service handler.
func NewService(svc *vendor.Service) *Service { return &Service{svc: svc} }

type serviceReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Unit        string `json:"unit"`
	IsActive    *bool  `json:"isActive,omitempty"`
}

func (req serviceReq) toInput() domain.ServiceInput {
	return domain.ServiceInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Unit:        domain.ServiceUnit(req.Unit),
		IsActive:    req.IsActive,
	}
}

// List returns active services for the vendor specified in the path.
func (h *Service) List(w http.ResponseWriter, r *http.Request) {
	vendorID := r.PathValue("id")
	role, _ := middleware.RoleFrom(r.Context())
	activeOnly := role != domain.RoleAdmin
	items, err := h.svc.ListServices(r.Context(), vendorID, activeOnly)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// MyList returns ALL services (active + inactive) for the caller's vendor profile.
func (h *Service) MyList(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	// Resolve caller's vendor via service helper (proxies through vendor repo).
	v, err := h.svc.MyVendor(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	items, err := h.svc.ListServices(r.Context(), v.ID, false)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Create publishes a new service on the caller's vendor profile.
func (h *Service) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req serviceReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	s, err := h.svc.CreateService(r.Context(), uid, req.toInput())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, s)
}

// Update modifies an existing service the caller owns.
func (h *Service) Update(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	var req serviceReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	s, err := h.svc.UpdateService(r.Context(), uid, id, req.toInput())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, s)
}

// Delete removes a service the caller owns.
func (h *Service) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.svc.DeleteService(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
