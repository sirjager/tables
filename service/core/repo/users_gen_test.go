package core_repo

// func createRandomUser(t *testing.T) CoreUser {
// 	arg := CreateUserParams{
// 		Email:    utils.RandomEmail(),
// 		Username: utils.RandomUserName(),
// 		Password: utils.RandomPassword(),
// 	}
// 	user, err := testQueries.CreateUser(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, user)

// 	require.Equal(t, arg.Email, user.Email)
// 	require.Equal(t, arg.Username, user.Username)
// 	require.Equal(t, arg.Password, user.Password)
// 	require.NotZero(t, user.ID)
// 	require.False(t, user.Public)
// 	require.False(t, user.Blocked)
// 	require.False(t, user.Verified)
// 	require.NotZero(t, user.Created)
// 	require.NotZero(t, user.Updated)
// 	return user
// }

// func TestAddCoreUser(t *testing.T) {
// 	createRandomUser(t)
// }

// func TestGetUser(t *testing.T) {
// 	user := createRandomUser(t)

// 	userInDb, err := testQueries.GetUser(context.Background(), user.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, userInDb)
// 	require.Equal(t, user.Email, userInDb.Email)
// 	require.Equal(t, user.Username, userInDb.Username)
// 	require.Equal(t, user.Password, userInDb.Password)
// 	require.Equal(t, user.ID, userInDb.ID)
// 	require.Equal(t, user.Public, userInDb.Public)
// 	require.Equal(t, user.Blocked, userInDb.Blocked)
// 	require.Equal(t, user.Verified, userInDb.Verified)
// 	require.Equal(t, user.Created, userInDb.Created)
// 	require.Equal(t, user.Updated, userInDb.Updated)

// 	userNotInDb, err := testQueries.GetUser(context.Background(), utils.RandomInt(1, 99999999999))
// 	require.Error(t, err)
// 	require.Empty(t, userNotInDb)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// }

// func TestGetUserWhereEmail(t *testing.T) {
// 	user := createRandomUser(t)

// 	userInDb, err := testQueries.GetUserWhereEmail(context.Background(), user.Email)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, userInDb)
// 	require.Equal(t, user.Email, userInDb.Email)
// 	require.Equal(t, user.Username, userInDb.Username)
// 	require.Equal(t, user.Password, userInDb.Password)
// 	require.Equal(t, user.ID, userInDb.ID)
// 	require.Equal(t, user.Public, userInDb.Public)
// 	require.Equal(t, user.Blocked, userInDb.Blocked)
// 	require.Equal(t, user.Verified, userInDb.Verified)
// 	require.Equal(t, user.Created, userInDb.Created)
// 	require.Equal(t, user.Updated, userInDb.Updated)

// 	userNotInDb, err := testQueries.GetUserWhereEmail(context.Background(), utils.RandomEmail())
// 	require.Error(t, err)
// 	require.Empty(t, userNotInDb)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// }

// func TestGetUserWhereUsername(t *testing.T) {
// 	user := createRandomUser(t)

// 	userInDb, err := testQueries.GetUserWhereUsername(context.Background(), user.Username)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, userInDb)
// 	require.Equal(t, user.Email, userInDb.Email)
// 	require.Equal(t, user.Username, userInDb.Username)
// 	require.Equal(t, user.Password, userInDb.Password)
// 	require.Equal(t, user.Fullname, userInDb.Fullname)
// 	require.Equal(t, user.ID, userInDb.ID)
// 	require.Equal(t, user.Public, userInDb.Public)
// 	require.Equal(t, user.Blocked, userInDb.Blocked)
// 	require.Equal(t, user.Verified, userInDb.Verified)
// 	require.Equal(t, user.Created, userInDb.Created)
// 	require.Equal(t, user.Updated, userInDb.Updated)

// 	userNotInDb, err := testQueries.GetUserWhereUsername(context.Background(), utils.RandomUserName())
// 	require.Error(t, err)
// 	require.Empty(t, userNotInDb)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// }

// func TestGetAllUsers(t *testing.T) {
// 	var createUsers []CoreUser
// 	const totalUsers = 10
// 	for i := 0; i < totalUsers; i++ {
// 		user := createRandomUser(t)
// 		createUsers = append(createUsers, user)
// 	}

// 	listedUsers, err := testQueries.GetAllUsers(context.Background())
// 	require.NoError(t, err)
// 	require.GreaterOrEqual(t, len(listedUsers), totalUsers)

// 	var foundUserUids []int64
// 	for _, gotUser := range listedUsers {
// 		for _, createdUser := range createUsers {
// 			checked := false
// 			for _, foundUid := range foundUserUids {
// 				if foundUid == createdUser.ID {
// 					checked = true
// 				}
// 			}
// 			if !checked {
// 				if gotUser.ID == createdUser.ID {
// 					foundUserUids = append(foundUserUids, gotUser.ID)
// 					require.Equal(t, gotUser.Email, createdUser.Email)
// 					require.Equal(t, gotUser.Username, createdUser.Username)
// 					require.Equal(t, gotUser.Password, createdUser.Password)
// 					require.Equal(t, gotUser.ID, createdUser.ID)
// 					require.Equal(t, gotUser.Fullname, createdUser.Fullname)
// 					require.Equal(t, gotUser.Public, createdUser.Public)
// 					require.Equal(t, gotUser.Blocked, createdUser.Blocked)
// 					require.Equal(t, gotUser.Verified, createdUser.Verified)
// 					require.Equal(t, gotUser.Created, createdUser.Created)
// 					require.Equal(t, gotUser.Updated, createdUser.Updated)
// 				}
// 			}

// 		}
// 	}

// 	require.Equal(t, len(foundUserUids), totalUsers)
// }

// func TestGetSomeUsers(t *testing.T) {
// 	const totalUsers = 5
// 	const fetchOnly = 3
// 	const offset_ = 3

// 	for i := 0; i < totalUsers; i++ {
// 		createRandomUser(t)
// 	}

// 	listedUsersWithLimit, err := testQueries.GetSomeUsers(context.Background(),
// 		GetSomeUsersParams{Offset: offset_, Limit: fetchOnly})
// 	require.NoError(t, err)
// 	require.Equal(t, len(listedUsersWithLimit), fetchOnly)
// 	for _, user := range listedUsersWithLimit {
// 		require.NotEmpty(t, user)
// 		require.NotEmpty(t, user.Email)
// 		require.NotEmpty(t, user.Username)
// 		require.NotEmpty(t, user.Password)
// 		require.NotZero(t, user.ID)
// 		require.False(t, user.Public)
// 		require.False(t, user.Blocked)
// 		require.False(t, user.Verified)
// 		require.NotZero(t, user.Created)
// 		require.NotZero(t, user.Updated)
// 	}
// }

// func TestDeleteUser(t *testing.T) {
// 	user := createRandomUser(t)
// 	err := testQueries.DeleteUser(context.Background(), user.ID)
// 	require.NoError(t, err)
// 	userShouldBeDeleted, err := testQueries.GetUser(context.Background(), user.ID)
// 	require.Error(t, err)
// 	require.Empty(t, userShouldBeDeleted)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// }

// func TestUpdateUserBlocked(t *testing.T) {
// 	user := createRandomUser(t)
// 	updatedUser, err := testQueries.UpdateUserBlocked(context.Background(), UpdateUserBlockedParams{ID: user.ID, Blocked: true})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Blocked, updatedUser.Blocked)
// 	require.False(t, user.Blocked)
// 	require.True(t, updatedUser.Blocked)

// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Username, updatedUser.Username)
// 	require.Equal(t, user.Password, updatedUser.Password)
// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Public, updatedUser.Public)
// 	require.Equal(t, user.Verified, updatedUser.Verified)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }

// func TestUpdateUserPublic(t *testing.T) {
// 	user := createRandomUser(t)
// 	updatedUser, err := testQueries.UpdateUserPublic(context.Background(), UpdateUserPublicParams{ID: user.ID, Public: true})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Public, updatedUser.Public)
// 	require.False(t, user.Public)
// 	require.True(t, updatedUser.Public)

// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Username, updatedUser.Username)
// 	require.Equal(t, user.Password, updatedUser.Password)
// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Blocked, updatedUser.Blocked)
// 	require.Equal(t, user.Verified, updatedUser.Verified)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }

// func TestUpdateUserVerified(t *testing.T) {
// 	user := createRandomUser(t)
// 	updatedUser, err := testQueries.UpdateUserVerified(context.Background(), UpdateUserVerifiedParams{ID: user.ID, Verified: true})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Verified, updatedUser.Verified)
// 	require.False(t, user.Verified)
// 	require.True(t, updatedUser.Verified)

// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Username, updatedUser.Username)
// 	require.Equal(t, user.Password, updatedUser.Password)
// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Blocked, updatedUser.Blocked)
// 	require.Equal(t, user.Public, updatedUser.Public)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }

// func TestUpdateUserFullName(t *testing.T) {
// 	user := createRandomUser(t)
// 	newName := utils.RandomUserName()

// 	updatedUser, err := testQueries.UpdateUserFullName(context.Background(), UpdateUserFullNameParams{ID: user.ID, Fullname: newName})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Fullname, updatedUser.Fullname)
// 	require.Equal(t, newName, updatedUser.Fullname)

// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Username, updatedUser.Username)
// 	require.Equal(t, user.Password, updatedUser.Password)
// 	require.Equal(t, user.Blocked, updatedUser.Blocked)
// 	require.Equal(t, user.Public, updatedUser.Public)
// 	require.Equal(t, user.Verified, updatedUser.Verified)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }

// func TestUpdateUserUsername(t *testing.T) {
// 	user := createRandomUser(t)

// 	newUsername := utils.RandomUserName()

// 	updatedUser, err := testQueries.UpdateUserUsername(context.Background(), UpdateUserUsernameParams{ID: user.ID, Username: newUsername})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Username, updatedUser.Username)
// 	require.Equal(t, newUsername, updatedUser.Username)

// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Fullname, updatedUser.Fullname)
// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Public, updatedUser.Public)
// 	require.Equal(t, user.Blocked, updatedUser.Blocked)
// 	require.Equal(t, user.Password, updatedUser.Password)
// 	require.Equal(t, user.Verified, updatedUser.Verified)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }

// func TestUpdateUserPassword(t *testing.T) {
// 	user := createRandomUser(t)
// 	newPassword := utils.RandomPassword()

// 	updatedUser, err := testQueries.UpdateUserPassword(context.Background(), UpdateUserPasswordParams{ID: user.ID, Password: newPassword})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedUser)
// 	require.NotEqual(t, user.Password, updatedUser.Password)
// 	require.Equal(t, newPassword, updatedUser.Password)

// 	require.Equal(t, user.ID, updatedUser.ID)
// 	require.Equal(t, user.Fullname, updatedUser.Fullname)
// 	require.Equal(t, user.Email, updatedUser.Email)
// 	require.Equal(t, user.Public, updatedUser.Public)
// 	require.Equal(t, user.Blocked, updatedUser.Blocked)
// 	require.Equal(t, user.Username, updatedUser.Username)
// 	require.Equal(t, user.Verified, updatedUser.Verified)
// 	require.Equal(t, user.Created, updatedUser.Created)
// 	require.Equal(t, user.Updated, updatedUser.Updated)
// }
