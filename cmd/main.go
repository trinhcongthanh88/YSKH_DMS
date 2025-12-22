package main

import (
	// client "YSKH_DMS/internal/controllers/client"
	
	routes "YSKH_DMS/internal/routes"
	db "YSKH_DMS/database"
	// "fmt"
	// "time"
)

// Import package mới theo đường dẫn tương đối (hoặc module path)
// Ví dụ: Giả sử project của bạn là 'scheduler_app' trong GOPATH/src/
// Bạn có thể dùng đường dẫn tuyệt đối nếu project đã được khởi tạo module
// Với ví dụ đơn giản này, ta dùng đường dẫn module:
// THAY THẾ 'scheduler_app' bằng tên module của bạn
func main() {
	db.InitDB()          // Kết nối DB
	// interval := 5 * time.Second

	// ⭐️ Gọi StartScheduler và truyền 2 hàm:
	// 1. Hàm gọi API: client.CallAPI_Post
	// 2. Hàm xử lý kết quả: PostResultProcessor
	// fmt.Println("Chương trình Gọi API")
	// go client.StartScheduler(interval, client.CallAPI_Post, client.PostResultProcessor)

	// fmt.Println("Chương trình chính đã khởi động.")
	r := routes.SetupRouter()
	r.Run(":4592") // Chạy server tại port 4592
	select {}
}

// PostResultProcessor là hàm xử lý kết quả cụ thể cho kiểu *model.Post.
// Nó nhận kết quả (post) và lỗi (err) từ hàm API đã gọi.
// Chữ ký của nó: func(result T, err error)
