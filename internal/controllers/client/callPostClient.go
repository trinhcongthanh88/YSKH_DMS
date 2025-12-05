package client

import (
	model "YSKH_DMS/internal/models"
	"fmt"
)

// Hàm thực hiện việc gọi API (giữ nguyên)

// callAPI thực hiện một HTTP GET request đến URL đã cho.
// Lưu ý: Tên hàm bắt đầu bằng chữ cái thường vì nó chỉ cần được dùng trong cùng package (main).
// Nếu bạn muốn dùng nó từ package khác, bạn phải đặt tên là CallAPI (chữ C viết hoa).

func CallAPI_Post() (*model.Post, error) {
	// Dữ liệu giả lập cho API B
	apiURL := "https://jsonplaceholder.typicode.com/posts/1"

	// Gọi hàm generics và chỉ định kiểu ResponseB
	return GenericAPICall[model.Post](apiURL)
}

func PostResultProcessor(post *model.Post, err error) {
	// 1. Kiểm tra lỗi (Quan trọng nhất)
	if err != nil {
		// In lỗi và dừng xử lý
		fmt.Printf("[Lỗi] Xảy ra lỗi khi gọi API: %v\n", err)
		return
	}

	// 2. Xử lý kết quả thành công
	// Lưu ý: Cần kiểm tra nil nếu post là con trỏ
	if post != nil {
		fmt.Printf("[Thành công] Nhận được Post ID: %d, Title: %s\n", post.ID, post.Title)
		// 💡 Ở đây, bạn có thể thêm logic khác như:
		// - Lưu post vào database
		// - Gửi thông báo
		// - ...
	} else {
		// Trường hợp API trả về thành công nhưng không có dữ liệu (ví dụ: 204 No Content)
		fmt.Println("[Thành công] API gọi thành công nhưng không có dữ liệu Post.")
	}
}
