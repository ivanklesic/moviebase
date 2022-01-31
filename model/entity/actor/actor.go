package actor

type Actor struct {
	Id        int `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ImagePath string `json:"-"`
	MovieIds []string `json:"movieIds"`
	MovieNames []string `json:"-"`
	Roles []string `json:"-"`
}

func (actor *Actor) Validate(method string) (bool, string, error) {
	actorExists, err := Exists("actor", []string{"first_name", "last_name"}, []interface{}{actor.FirstName, actor.LastName})
	if err != nil {
		return false, "", err
	}
	if len(actor.FirstName) > 255 || len(actor.LastName) > 255 {
		return false, "First or last name is too long", nil
	}

	if method == "create" {
		if actorExists {
			return false, "Actor with the same name and last name already exists", nil
		}
	}

	return true, "", nil
}
