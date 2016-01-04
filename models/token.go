package models

type Token struct {
	User  string `json:"user" xorm:"pk"`
	Token string `json:"token"`
}

func (t Token) Store() {
	// TODO Убрать костыль
	// Удаляем существующий токен перед перезаписью
	t.Delete()

	_, err := Orm.Insert(&t)
	if err != nil {
		panic(err)
	}
}

func (t Token) Delete() (int64, error) {
	return Orm.Delete(&t)
}

func GetToken(user string) (string, error) {
	token := Token{User: user}
	_, err := Orm.Get(&token)

	return token.Token, err
}

func init() {
	err := Orm.CreateTables(&Token{})
	if err != nil {
		panic(err)
	}
}