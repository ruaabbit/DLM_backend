package models

// User 定义用户模型
type User struct {
	ID       int    `json:"id" gorm:"primaryKey"` // 主键ID
	Username string `json:"username"`
	Password string `json:"password"` // 正式环境下请使用加密存储
	Role     string `json:"role"`     // 例如 "admin" 或 "keeper"
	Name     string `json:"name"`     // 个人姓名
	Phone    string `json:"phone"`    // 手机号
}
