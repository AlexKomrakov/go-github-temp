package models

type Token struct {
	User  string `json:"user" xorm:"pk"`
	Token string `json:"token"`
}

func (t *Token) Store() (int64, error) {
	// TODO Убрать костыль
	// Удаляем существующий токен перед перезаписью
	t.Delete()

	return Orm.Insert(t)
}

func (t Token) Delete() (int64, error) {
	return Orm.Delete(&t)
}

func (t *Token) FindOne() (bool, error) {
	return Orm.Get(t)
}