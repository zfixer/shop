package client

import (
	"shop/logger"
	"shop/models"
)

type UserService interface {
	//GetAll() []datamodels.User
	//GetByID(id int64) (datamodels.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (*models.User, bool)
	//DeleteByID(id int64) bool
	//
	//Update(id int64, user datamodels.User) (datamodels.User, error)
	//UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	//UpdateUsername(id int64, newUsername string) (datamodels.User, error)
	//
	//Create(userPassword string, user datamodels.User) (datamodels.User, error)

	CreateUserTable()
	InsertUser(user *models.User)
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{
	}
}

func (u *userService) CreateUserTable() {
	models.UserCreateTable()
}

func (u *userService) InsertUser(user *models.User) {
	models.UserInsert(user)
}

func (u *userService) FindAllUser( users []*models.User) {
	users = models.UserFindAll()
}

func (u *userService) GetByUsernameAndPassword(username, userPassword string) (*models.User, bool) {
	user := models.UserFindByName(username)
	if user != nil && len(user.Username) > 0 {
		logger.Info.Println("找到用户名=", user.Username)
		return user, true
	} else {
		logger.Info.Println("没找到用户名=", username)
		return user, false
	}
}
