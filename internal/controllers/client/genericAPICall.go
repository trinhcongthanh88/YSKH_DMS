package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Khai báo một Generics Constraint
// T là kiểu dữ liệu bất kỳ (any) mà ta muốn giải mã JSON vào
type AnyResponse interface{}
type GenericAPIHandler[T AnyResponse] func() (T, error)
type ResultHandler[T any] func(result T, err error)

// genericAPICall là hàm chính để thực hiện yêu cầu HTTP và giải mã JSON
// Nó chấp nhận kiểu dữ liệu trả về mong muốn là T
func GenericAPICall[T AnyResponse](url string) (*T, error) {
	// 1. Thực hiện HTTP Request (giả lập)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gọi API: %w", err)
	}
	defer resp.Body.Close()

	// 2. Khởi tạo biến để lưu trữ kết quả
	// Sử dụng new(T) để tạo một con trỏ tới một giá trị kiểu T
	var result T

	// 3. Giải mã JSON vào con trỏ (&result)
	// Marshal sẽ điền dữ liệu vào cấu trúc T
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã JSON: %w", err)
	}

	// Trả về con trỏ tới kết quả đã giải mã
	return &result, nil
}

// File: scheduler/scheduler.go (Sửa lại)

// StartScheduler (Bản Generic, nhận 2 hàm xử lý)
func StartScheduler[T any](
	interval time.Duration,
	handler GenericAPIHandler[T],
	resultProcessor ResultHandler[T]) { // ⭐️ Tham số thứ ba: hàm xử lý kết quả

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		fmt.Println("=> [Scheduler] Lịch trình Generic được kích hoạt.")

		for {
			<-ticker.C

			// 1. GỌI HÀM HANDLER ĐỂ LẤY KẾT QUẢ
			result, err := handler()

			// 2. TRUYỀN KẾT QUẢ VÀO HÀM XỬ LÝ
			// Hàm này sẽ đảm nhận việc kiểm tra 'err' và in kết quả.
			resultProcessor(result, err)
		}
	}()
}

// func StartScheduler[T AnyResponse](number time.Duration, handler GenericAPIHandler[T]) {
// 	interval := number * time.Second

// 	// 1. Khởi tạo Ticker
// 	ticker := time.NewTicker(interval)
// 	// Đảm bảo dừng Ticker khi hàm main kết thúc để giải phóng tài nguyên
// 	defer ticker.Stop()

// 	fmt.Println("Chương trình bắt đầu chạy với Ticker. Nhấn Ctrl+C để dừng.")

// 	// Vòng lặp chờ tín hiệu từ Ticker
// 	for {
// 		handler()
// 		// result, err := handler()

// 		// if err != nil {
// 		// 	fmt.Printf("[Lỗi] Xảy ra lỗi: %v\n", err)
// 		// 	continue
// 		// }

// 		// // Xử lý kết quả thành công
// 		// // Vì T là generic, chúng ta không biết chính xác nó là Post hay User,
// 		// // nên ta chỉ có thể in ra kiểu và giá trị chung (ví dụ: dùng %v)
// 		// fmt.Printf("[Thành công] Nhận được dữ liệu kiểu %T: %v\n", result, result)
// 	}

// }
