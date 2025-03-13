package models

// User 定义用户模型
type User struct {
	ID       int    `json:"id" gorm:"primaryKey"` // 主键ID
	Username string `json:"username"`
	Password string `json:"-" gorm:"column:password"` // 使用json:"-"来在JSON序列化时忽略该字段
	Role     string `json:"role"`                     // 例如 "admin" 或 "keeper"
	Name     string `json:"name"`                     // 个人姓名
	Phone    string `json:"phone"`                    // 手机号
}
