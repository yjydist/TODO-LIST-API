package model

type TODO struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text"`
	UserID      int    `json:"user_id" gorm:"not null"`
	User        User   `json:"-" gorm:"foreignKey:UserID"`
	Completed   bool   `json:"completed" gorm:"default:false"`
}

// TableName 指定表名
func (t *TODO) TableName() string {
	return "todos"
}

// CreateTODO 创建待办事项
func CreateTODO(todo *TODO) error {
	return DB.Create(todo).Error
}

// UpdateTODO 更新待办事项
func UpdateTODO(todo *TODO) error {
	return DB.Save(todo).Error
}

// DeleteTODOByID 删除待办事项
func DeleteTODOByID(id int) error {
	return DB.Delete(&TODO{}, id).Error
}

// GetTODOByID 通过ID获取待办事项
func GetTODOByID(id int) (TODO, error) {
	var todo TODO
	result := DB.First(&todo, id)
	return todo, result.Error
}

// GetTODOsByUserID 获取用户的所有待办事项（带分页和筛选）
func GetTODOsByUserID(userID int, page, limit int, filter map[string]interface{}) ([]TODO, int64, error) {
	var todos []TODO
	var count int64

	// 基础查询
	query := DB.Model(&TODO{}).Where("user_id = ?", userID)

	// 应用筛选条件
	if filter != nil {
		for key, value := range filter {
			if key == "completed" {
				query = query.Where("completed = ?", value)
			} else if key == "title" {
				query = query.Where("title LIKE ?", "%"+value.(string)+"%")
			}
		}
	}

	// 计算总记录数
	query.Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&todos)

	return todos, count, result.Error
}

// IsOwner 检查用户是否为待办事项的所有者
func IsOwner(todoID, userID int) bool {
	var todo TODO
	result := DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo)
	return result.Error == nil
}
