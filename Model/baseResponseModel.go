package Model

// Pagination đại diện cho cấu trúc "pagination"
type Pagination struct {
	Last             bool `json:"last"`
	First            bool `json:"first"`
	Empty            bool `json:"empty"`
	NumberOfElements int  `json:"numberOfElements"`
	TotalElements    int  `json:"totalElements"`
	TotalPages       int  `json:"totalPages"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

// ResponseGeneric là cấu trúc model generic
// T là type parameter, sẽ được thay thế bằng kiểu dữ liệu thực tế của trường 'data'
type ResponseGeneric[T any] struct {
	StatusCode string     `json:"statusCode"`
	Message    string     `json:"message"`
	Data       T          `json:"data"` // Trường Data sử dụng kiểu generic T
	Pagination Pagination `json:"pagination"`
}
