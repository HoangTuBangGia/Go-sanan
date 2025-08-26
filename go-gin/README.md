## Tìm hiểu Gin (Go)

### 1. Gin là gì?

Gin là một web framework viết bằng Go, tối ưu cho hiệu năng (dựa trên `httprouter`) và trải nghiệm phát triển nhanh gọn. Nó cung cấp:

- HTTP router nhanh (tree-based routing)
- Middleware (logger, recovery, custom...)
- Binding & validation (JSON, form, query, header, uri)
- Context thống nhất (`*gin.Context`)
- Route grouping & versioning
- Nhiều renderer built-in (JSON, XML, HTML, file...)

### 2. Cấu trúc dự án mẫu (thư mục `route-group`)

```
route-group/
  main.go
  internal/api/
	v1/handler/{user.go, product.go}
	v2/handler/{user.go}
```

Ý nghĩa:

- `main.go`: Khởi tạo server, định nghĩa group `/api/v1`, `/api/v2` và resource.
- `internal/api/v1/handler`: Handler v1 (users, products).
- `internal/api/v2/handler`: Handler v2 (users).

### 3. Khởi tạo cơ bản

```go
r := gin.Default() // Logger + Recovery
r.Run(":8080")
```

### 4. Route, URL, Path và các khái niệm

- URL: `http://localhost:8080/api/v1/users/123?active=true`
- Path: `/api/v1/users/123`
- Route pattern: `GET /api/v1/users/:id`
- Method: GET/POST/PUT/PATCH/DELETE...
- Resource: thực thể nghiệp vụ (user, product...)

### 5. Path Param & Query Param

- Path param: khai báo bằng `:name` trong pattern, đọc bằng `ctx.Param("id")`.
- Query param: sau `?`, đọc bằng `ctx.Query("page")` hoặc `ctx.DefaultQuery("page","1")`.
  Ví dụ: `GET /api/v1/users/15?page=2&active=true` → id=15, page=2, active=true.

### 6. Route Grouping

Khái niệm: gom nhiều route chung prefix + (tuỳ chọn) middleware.
Ví dụ trong `main.go` (rút gọn):

```go
v1 := r.Group("/api/v1")
users := v1.Group("/users")
users.GET("/", h.GetUsersV1)
users.GET(":id", h.GetUsersByIdV1)
```

Lợi ích:

- DRY prefix `/api/v1`
- Áp middleware theo version / module
- Điều khiển versioning rõ ràng
- Giảm xung đột route & tăng tính tổ chức
  Có thể lồng: `api := r.Group("/api"); v1 := api.Group("/v1"); ...`

### 7. HTTP Methods & REST cơ bản

| Method | Mục đích          | Ví dụ                      | Trạng thái phổ biến |
| ------ | ----------------- | -------------------------- | ------------------- |
| GET    | Lấy danh sách     | `GET /api/v1/users`        | 200                 |
| GET    | Lấy 1             | `GET /api/v1/users/:id`    | 200 / 404           |
| POST   | Tạo mới           | `POST /api/v1/users`       | 201                 |
| PUT    | Thay thế toàn bộ  | `PUT /api/v1/users/:id`    | 200 / 404           |
| PATCH  | Cập nhật một phần | (chưa dùng)                | 200 / 404           |
| DELETE | Xóa               | `DELETE /api/v1/users/:id` | 204 / 404           |

### 8. Handler là gì & ví dụ

Handler là hàm được gắn với một route (method + path pattern). Nó nhận `*gin.Context`, đọc đầu vào (path param, query param, body), thực thi logic và ghi response.

```go
func (u *UserHandler) GetUsersByIdV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get user by ID (v1)", "id": id})
}
```

### 9. Response & Render

- JSON nhanh: `c.JSON(201, gin.H{"message":"Create"})`
- Các renderer khác:
  - `c.IndentedJSON(...)`, `c.PureJSON(...)`, `c.XML(...)`, `c.YAML(...)`, `c.ProtoBuf(...)`
  - HTML: `r.LoadHTMLGlob("templates/**/*"); c.HTML(200, "index.tmpl", gin.H{"title":"Home"})`
  - File: `c.File(path)`, `c.FileAttachment(path, name)`
  - Redirect: `c.Redirect(302, url)`
  - Raw bytes: `c.Data(200, "text/plain", []byte("hi"))`
    `gin.H` = map[string]any giản tiện.

### 10. Middleware

Khái niệm: hàm trung gian (pipeline) có thể chạy trước & sau handler.
Chữ ký: `type HandlerFunc func(*gin.Context)`
Điều khiển luồng:

- `c.Next()` tiếp tục
- `c.Abort()` / `c.AbortWithStatusJSON(...)` dừng các handler sau
  Ví dụ log thời gian:

```go
func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s %d %v", c.Request.Method, c.FullPath(), c.Writer.Status(), time.Since(start))
	}
}
```

Ví dụ auth đơn giản:

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing token"})
			return
		}
		c.Set("userID", "123")
		c.Next()
	}
}
```

Áp dụng: `r.Use(gin.Logger(), gin.Recovery(), TimingMiddleware())` hoặc `group := v1.Group("/users", AuthMiddleware())`
Use cases: Auth, logging, tracing, rate limiting, CORS, compression, metrics, security headers.

### 11. Binding & Validation (khái niệm)

- Binding: ánh xạ tự động dữ liệu request (JSON, query, form, uri, header) vào struct Go qua tag (ví dụ `json:"name"`) thay vì parse thủ công.
- Validation: kiểm tra ràng buộc dữ liệu (required, email, min, max...) khai báo trong tag `binding:"..."`; nếu vi phạm trả về 400.
  Nguồn thường dùng: `ShouldBindJSON`, `ShouldBindQuery`, `ShouldBind`, `ShouldBindUri`, `ShouldBindHeader`.
  Ví dụ tối giản:

```go
type CreateUserDTO struct {
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=1,lte=120"`
}
func (h *UserHandler) PostUsersV1(c *gin.Context) {
	var body CreateUserDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Create user (v1)", "data": body})
}
```

### 12. API Versioning (khái niệm & lý do)

Khái niệm: tách API thành các phiên bản (v1, v2, ...) để có thể thay đổi phá vỡ mà không làm hỏng client đang dùng phiên bản cũ.
Tại sao cần: giữ tương thích ngược, cho phép tiến hóa dần, giảm rủi ro triển khai.
Ví dụ: dùng prefix `/api/v1` và `/api/v2`.

### 13. Trạng thái HTTP phổ biến

- 200 OK: Thành công chung
- 201 Created: Tạo mới
- 204 No Content: Xóa thành công
- 400 Bad Request: Dữ liệu sai
- 401 / 403: Sai/thiếu auth hoặc không đủ quyền
- 404 Not Found: Không tồn tại
- 500 Internal Server Error: Lỗi máy chủ

### 14. Chạy thử

```bash
cd route-group
go mod tidy
go run .
```

Test nhanh:

```bash
curl -s http://localhost:8080/api/v1/users/
curl -s http://localhost:8080/api/v1/users/123
curl -s -X POST http://localhost:8080/api/v1/users/ -H 'Content-Type: application/json' -d '{"name":"Alice","email":"a@b.com","age":30}'
```
