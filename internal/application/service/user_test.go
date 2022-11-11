package service

// func TestUserGetAll(t *testing.T) {
// 	assert := assert.New(t)

// 	userRepo := mocks.NewUserRepo(t)
// 	userRepo.EXPECT().GetAll(mock.Anything).Return([]model.User{{BaseModel: model.BaseModel{ID: 1}, Username: "aaron"}}, nil)

// 	userService := NewUserService(userRepo)
// 	users, _ := userService.GetUserList(context.Background())
// 	assert.NotEmpty(users)
// 	userRepo.AssertExpectations(t)
// }
