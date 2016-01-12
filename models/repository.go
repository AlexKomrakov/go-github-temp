package models
import "github.com/google/go-github/github"

type Repository struct {
	Id      int64  `json:"id"`
	Login   string `json:"login"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// TODO Add sorting order
// err = getDb().C(builds_collection).Find(bson.M{"login": r.Login, "name": r.Name}).Sort("-_id").All(&builds)
func (repository Repository) Builds() (builds []Build, err error) {
	err = Orm.Find(&builds, &Build{RepositoryId: repository.Id})
	return
}

// TODO Добавить проверку на уникальность репозиториев
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

func GetGithubRepositoriesIntersection(github_repos []github.Repository) (repositories []Repository, err error) {
	for _, github_repo := range github_repos {
		repo := &Repository{Login: *github_repo.Owner.Login, Name: *github_repo.Name};
		success, err := repo.FindOne()
		if err != nil {
			return repositories, err
		}
		if success == true {
			repositories = append(repositories, *repo)
		}
	}

	return
}

