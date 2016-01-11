package models

type Repository struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func (repository Repository) Builds() (builds []Build, err error) {
	// TODO add sorting
	// err = getDb().C(builds_collection).Find(bson.M{"login": r.Login, "name": r.Name}).Sort("-_id").All(&builds)
	err = Orm.Find(&builds, &Build{RepositoryId: repository.Id})
	return
}

func (repository *Repository) Store() (int64, error) {
	return Orm.Insert(repository)
}

func (repository *Repository) FindOne() (bool, error) {
	return Orm.Get(repository)
}

func (repository *Repository) FindOrCreate() (bool, error) {
	success, err := repository.FindOne()
	if success == false {
		number, err := Orm.Insert(repository)
		return number != 0, err
	}

	return true, err
}

func (repository Repository) Delete() (int64, error) {
	return Orm.Delete(&repository)
}

