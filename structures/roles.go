package structures

type Role struct {
	ID         string `db:"id"`
	Experience int64  `db:"experience"`
}

type Roles []*Role

func (r Roles) Len() int           { return len(r) }
func (r Roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Roles) Less(i, j int) bool { return r[i].Experience < r[j].Experience }
