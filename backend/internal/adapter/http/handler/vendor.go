package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/vendor"
)

// Vendor HTTP handler (vendor profile + photos + catalog).
type Vendor struct {
	svc *vendor.Service
}

// NewVendor constructs a Vendor handler.
func NewVendor(svc *vendor.Service) *Vendor { return &Vendor{svc: svc} }

type vendorReq struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	City        string `json:"city"`
	Description string `json:"description"`
	PriceFrom   int64  `json:"priceFrom"`
}

// Upsert creates / updates the caller's vendor profile.
func (h *Vendor) Upsert(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req vendorReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	v, err := h.svc.Upsert(r.Context(), uid, domain.VendorInput{
		Name: req.Name, Category: req.Category, City: req.City,
		Description: req.Description, PriceFrom: req.PriceFrom,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

// My returns the caller's vendor profile.
func (h *Vendor) My(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	v, err := h.svc.MyVendor(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

// List runs the catalog search.
func (h *Vendor) List(w http.ResponseWriter, r *http.Request) {
	q := buildVendorQuery(r)
	role, _ := middleware.RoleFrom(r.Context())

	var (
		items []*domain.Vendor
		total int
		err   error
	)
	if role == domain.RoleAdmin {
		items, total, err = h.svc.AdminSearch(r.Context(), q)
	} else {
		items, total, err = h.svc.PublicSearch(r.Context(), q)
	}
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"page":  q.Page,
		"limit": q.Limit,
	})
}

// Detail returns one vendor (with visibility check).
func (h *Vendor) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := h.svc.ByID(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	role, _ := middleware.RoleFrom(r.Context())
	uid, _ := middleware.UserIDFrom(r.Context())
	if v.Status != domain.VendorApproved && role != domain.RoleAdmin && v.UserID != uid {
		httpx.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

// UploadPhoto attaches an image to the caller's vendor profile.
func (h *Vendor) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	if err := r.ParseMultipartForm(domain.MaxPhotoSize + 1024); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "form parse: "+err.Error())
		return
	}
	file, header, err := r.FormFile("photo")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "missing 'photo' file field")
		return
	}
	defer file.Close()
	if header.Size > domain.MaxPhotoSize {
		httpx.WriteError(w, http.StatusRequestEntityTooLarge, "max 5MB")
		return
	}
	data, err := io.ReadAll(io.LimitReader(file, domain.MaxPhotoSize))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "read failed")
		return
	}
	mime := header.Header.Get("Content-Type")
	p, err := h.svc.UploadPhoto(r.Context(), uid, mime, data)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, p)
}

// ServePhoto streams a photo blob.
func (h *Vendor) ServePhoto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := h.svc.FindPhoto(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", p.MIME)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Length", strconv.FormatInt(p.Size, 10))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p.Data)
}

// DeletePhoto removes a photo (if owned).
func (h *Vendor) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.svc.DeletePhoto(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func buildVendorQuery(r *http.Request) usecase.VendorQuery {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	minP, _ := strconv.ParseInt(q.Get("price_min"), 10, 64)
	maxP, _ := strconv.ParseInt(q.Get("price_max"), 10, 64)
	minR, _ := strconv.ParseFloat(q.Get("rating_min"), 64)
	return usecase.VendorQuery{
		Q:         strings.TrimSpace(q.Get("q")),
		Category:  q.Get("category"),
		City:      q.Get("city"),
		MinPrice:  minP,
		MaxPrice:  maxP,
		MinRating: minR,
		Sort:      q.Get("sort"),
		Page:      page,
		Limit:     limit,
	}
}
