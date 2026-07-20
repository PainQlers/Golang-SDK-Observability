package models

// User ตัวแทนโครงสร้างข้อมูลที่เราจะใช้ส่งไปมาในระบบ
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}