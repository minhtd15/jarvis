package panther

type userService struct {
	userStore UserStore
}

func (u userService) GetByUserName(userName string) (UserEntity, error) {
	return u.userStore.GetByUserName(userName)
}

type UserStoreCfg struct {
	UserStore UserStore
}

func NewUserService(cfg UserStoreCfg) UserService {
	return &userService{
		userStore: cfg.UserStore,
	}
}
