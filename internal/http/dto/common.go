package dto

type DataEnvelope[T any] struct {
	Data T `json:"data"`
}

type Meta struct {
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
	Total   int64 `json:"total"`
	HasNext bool  `json:"has_next"`
}

type ListEnvelope[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}
