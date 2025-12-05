package main

import (
	client "YSKH_DMS/internal/controllers/client"
	"fmt"
	"time"
)

// Import package mới theo đường dẫn tương đối (hoặc module path)
// Ví dụ: Giả sử project của bạn là 'scheduler_app' trong GOPATH/src/
// Bạn có thể dùng đường dẫn tuyệt đối nếu project đã được khởi tạo module
// Với ví dụ đơn giản này, ta dùng đường dẫn module:
// THAY THẾ 'scheduler_app' bằng tên module của bạn
func main() {
	interval := 5 * time.Second

	// ⭐️ Gọi StartScheduler và truyền 2 hàm:
	// 1. Hàm gọi API: client.CallAPI_Post
	// 2. Hàm xử lý kết quả: PostResultProcessor
	fmt.Println("Chương trình Gọi API")
	go client.StartScheduler(interval, client.CallAPI_Post, client.PostResultProcessor)

	fmt.Println("Chương trình chính đã khởi động.")
	select {}
}

// PostResultProcessor là hàm xử lý kết quả cụ thể cho kiểu *model.Post.
// Nó nhận kết quả (post) và lỗi (err) từ hàm API đã gọi.
// Chữ ký của nó: func(result T, err error)
