package todolist

type User struct {
	//чтобы метод get работал нужно прописать тег db с названием поля из базы
	Id       int    `json: "-"`                      //db:"id"
	Name     string `json:"name" binding:"required"` //тег binding:"required" валидирует наличие данных полей в теле запроса
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
