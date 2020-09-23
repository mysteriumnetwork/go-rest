package paging

const (
	// DefaultSize - default page size if unspecified.
	DefaultSize = 20
	// MaxSize - maximum page size (cap).
	MaxSize = 50
)

// Request represents a paging request.
type Request struct {
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"page_size"`
}

// Pageable holds pagination information.
type Pageable struct {
	Page       uint64 `json:"page"`
	PageSize   uint64 `json:"page_size"`
	TotalItems uint64 `json:"total_items"`
	TotalPages uint64 `json:"total_pages"`
}
