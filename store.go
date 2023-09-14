package education_website

type Store interface {
	UserStore() UserStore
}
