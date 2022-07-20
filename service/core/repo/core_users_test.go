package core_repo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/SirJager/tables/service/core/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) CoreUser {
	arg := AddCoreUserParams{
		Email:    utils.RandomEmail(),
		Username: utils.RandomUserName(),
		Password: utils.RandomPassword(),
	}
	user, err := testQueries.AddCoreUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.NotZero(t, user.Uid)
	require.False(t, user.IsPublic)
	require.False(t, user.IsBlocked)
	require.False(t, user.IsVerified)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	return user
}

func TestAddCoreUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetCoreUserWithEmail(t *testing.T) {
	user := createRandomUser(t)

	userInDb, err := testQueries.GetCoreUserWithEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userInDb)
	require.Equal(t, user.Email, userInDb.Email)
	require.Equal(t, user.Username, userInDb.Username)
	require.Equal(t, user.Password, userInDb.Password)
	require.Equal(t, user.Uid, userInDb.Uid)
	require.Equal(t, user.IsPublic, userInDb.IsPublic)
	require.Equal(t, user.IsBlocked, userInDb.IsBlocked)
	require.Equal(t, user.IsVerified, userInDb.IsVerified)
	require.Equal(t, user.CreatedAt, userInDb.CreatedAt)
	require.Equal(t, user.UpdatedAt, userInDb.UpdatedAt)

	userNotInDb, err := testQueries.GetCoreUserWithEmail(context.Background(), utils.RandomEmail())
	require.Error(t, err)
	require.Empty(t, userNotInDb)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestGetCoreUserWithUid(t *testing.T) {
	user := createRandomUser(t)

	userInDb, err := testQueries.GetCoreUserWithUid(context.Background(), user.Uid)
	require.NoError(t, err)
	require.NotEmpty(t, userInDb)
	require.Equal(t, user.Email, userInDb.Email)
	require.Equal(t, user.Username, userInDb.Username)
	require.Equal(t, user.Password, userInDb.Password)
	require.Equal(t, user.Uid, userInDb.Uid)
	require.Equal(t, user.IsPublic, userInDb.IsPublic)
	require.Equal(t, user.IsBlocked, userInDb.IsBlocked)
	require.Equal(t, user.IsVerified, userInDb.IsVerified)
	require.Equal(t, user.CreatedAt, userInDb.CreatedAt)
	require.Equal(t, user.UpdatedAt, userInDb.UpdatedAt)

	userNotInDb, err := testQueries.GetCoreUserWithUid(context.Background(), utils.RandomInt(1, 99999999999))
	require.Error(t, err)
	require.Empty(t, userNotInDb)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestGetCoreUserWithUsername(t *testing.T) {
	user := createRandomUser(t)

	userInDb, err := testQueries.GetCoreUserWithUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userInDb)
	require.Equal(t, user.Email, userInDb.Email)
	require.Equal(t, user.Username, userInDb.Username)
	require.Equal(t, user.Password, userInDb.Password)
	require.Equal(t, user.Fullname, userInDb.Fullname)
	require.Equal(t, user.Uid, userInDb.Uid)
	require.Equal(t, user.IsPublic, userInDb.IsPublic)
	require.Equal(t, user.IsBlocked, userInDb.IsBlocked)
	require.Equal(t, user.IsVerified, userInDb.IsVerified)
	require.Equal(t, user.CreatedAt, userInDb.CreatedAt)
	require.Equal(t, user.UpdatedAt, userInDb.UpdatedAt)

	userNotInDb, err := testQueries.GetCoreUserWithUsername(context.Background(), utils.RandomUserName())
	require.Error(t, err)
	require.Empty(t, userNotInDb)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListCoreUsers(t *testing.T) {
	var createUsers []CoreUser
	const totalUsers = 10
	for i := 0; i < totalUsers; i++ {
		user := createRandomUser(t)
		createUsers = append(createUsers, user)
	}

	listedUsers, err := testQueries.ListCoreUsers(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(listedUsers), totalUsers)

	var foundUserUids []int64
	for _, gotUser := range listedUsers {
		for _, createdUser := range createUsers {
			checked := false
			for _, foundUid := range foundUserUids {
				if foundUid == createdUser.Uid {
					checked = true
				}
			}
			if !checked {
				if gotUser.Uid == createdUser.Uid {
					foundUserUids = append(foundUserUids, gotUser.Uid)
					require.Equal(t, gotUser.Email, createdUser.Email)
					require.Equal(t, gotUser.Username, createdUser.Username)
					require.Equal(t, gotUser.Password, createdUser.Password)
					require.Equal(t, gotUser.Uid, createdUser.Uid)
					require.Equal(t, gotUser.Fullname, createdUser.Fullname)
					require.Equal(t, gotUser.IsPublic, createdUser.IsPublic)
					require.Equal(t, gotUser.IsBlocked, createdUser.IsBlocked)
					require.Equal(t, gotUser.IsVerified, createdUser.IsVerified)
					require.Equal(t, gotUser.CreatedAt, createdUser.CreatedAt)
					require.Equal(t, gotUser.UpdatedAt, createdUser.UpdatedAt)
				}
			}

		}
	}

	require.Equal(t, len(foundUserUids), totalUsers)
}

func TestListCoreUsersWithLimit(t *testing.T) {
	const totalUsers = 5
	const fetchOnly = 3

	for i := 0; i < totalUsers; i++ {
		createRandomUser(t)
	}

	listedUsersWithLimit, err := testQueries.ListCoreUsersWithLimit(context.Background(), fetchOnly)
	require.NoError(t, err)
	require.Equal(t, len(listedUsersWithLimit), fetchOnly)
	for _, user := range listedUsersWithLimit {
		require.NotEmpty(t, user)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.Password)
		require.NotZero(t, user.Uid)
		require.False(t, user.IsPublic)
		require.False(t, user.IsBlocked)
		require.False(t, user.IsVerified)
		require.NotZero(t, user.CreatedAt)
		require.NotZero(t, user.UpdatedAt)
	}
}

func TestListCoreUsersWithLimitOffset(t *testing.T) {
	const totalUsers = 5
	const fetchOnly = 3
	const offset_ = 3

	for i := 0; i < totalUsers; i++ {
		createRandomUser(t)
	}

	listedUsersWithLimit, err := testQueries.ListCoreUsersWithLimitOffset(context.Background(),
		ListCoreUsersWithLimitOffsetParams{Offset: offset_, Limit: fetchOnly})
	require.NoError(t, err)
	require.Equal(t, len(listedUsersWithLimit), fetchOnly)
	for _, user := range listedUsersWithLimit {
		require.NotEmpty(t, user)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.Password)
		require.NotZero(t, user.Uid)
		require.False(t, user.IsPublic)
		require.False(t, user.IsBlocked)
		require.False(t, user.IsVerified)
		require.NotZero(t, user.CreatedAt)
		require.NotZero(t, user.UpdatedAt)
	}
}

func TestRemoveCoreUserWithEmail(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.RemoveCoreUserWithEmail(context.Background(), user.Email)
	require.NoError(t, err)

	userShouldBeDeleted, err := testQueries.GetCoreUserWithEmail(context.Background(), user.Email)
	require.Error(t, err)
	require.Empty(t, userShouldBeDeleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestRemoveCoreUserWithUid(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.RemoveCoreUserWithUid(context.Background(), user.Uid)
	require.NoError(t, err)
	userShouldBeDeleted, err := testQueries.GetCoreUserWithUid(context.Background(), user.Uid)
	require.Error(t, err)
	require.Empty(t, userShouldBeDeleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestRemoveCoreUserWithUsername(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.RemoveCoreUserWithUsername(context.Background(), user.Username)
	require.NoError(t, err)

	userShouldBeDeleted, err := testQueries.GetCoreUserWithUsername(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, userShouldBeDeleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestUpdateCoreUserBlocked(t *testing.T) {
	user := createRandomUser(t)
	updatedUser, err := testQueries.UpdateCoreUserBlocked(context.Background(), UpdateCoreUserBlockedParams{Uid: user.Uid, IsBlocked: true})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.IsBlocked, updatedUser.IsBlocked)
	require.False(t, user.IsBlocked)
	require.True(t, updatedUser.IsBlocked)

	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.IsPublic, updatedUser.IsPublic)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateCoreUserPublic(t *testing.T) {
	user := createRandomUser(t)
	updatedUser, err := testQueries.UpdateCoreUserPublic(context.Background(), UpdateCoreUserPublicParams{Uid: user.Uid, IsPublic: true})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.IsPublic, updatedUser.IsPublic)
	require.False(t, user.IsPublic)
	require.True(t, updatedUser.IsPublic)

	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.IsBlocked, updatedUser.IsBlocked)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateCoreUserVerified(t *testing.T) {
	user := createRandomUser(t)
	updatedUser, err := testQueries.UpdateCoreUserVerified(context.Background(), UpdateCoreUserVerifiedParams{Uid: user.Uid, IsVerified: true})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.IsVerified, updatedUser.IsVerified)
	require.False(t, user.IsVerified)
	require.True(t, updatedUser.IsVerified)

	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.IsBlocked, updatedUser.IsBlocked)
	require.Equal(t, user.IsPublic, updatedUser.IsPublic)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateCoreUserName(t *testing.T) {
	user := createRandomUser(t)
	newName := utils.RandomUserName()

	updatedUser, err := testQueries.UpdateCoreUserName(context.Background(), UpdateCoreUserNameParams{Uid: user.Uid, Fullname: newName})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Fullname, updatedUser.Fullname)
	require.Equal(t, newName, updatedUser.Fullname)

	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.IsBlocked, updatedUser.IsBlocked)
	require.Equal(t, user.IsPublic, updatedUser.IsPublic)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateCoreUserUsername(t *testing.T) {
	user := createRandomUser(t)

	newUsername := utils.RandomUserName()

	updatedUser, err := testQueries.UpdateCoreUserUsername(context.Background(), UpdateCoreUserUsernameParams{Uid: user.Uid, Username: newUsername})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Username, updatedUser.Username)
	require.Equal(t, newUsername, updatedUser.Username)

	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.Fullname, updatedUser.Fullname)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.IsPublic, updatedUser.IsPublic)
	require.Equal(t, user.IsBlocked, updatedUser.IsBlocked)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}

func TestUpdateCoreUserPassword(t *testing.T) {
	user := createRandomUser(t)
	newPassword := utils.RandomPassword()

	updatedUser, err := testQueries.UpdateCoreUserPassword(context.Background(), UpdateCoreUserPasswordParams{Uid: user.Uid, Password: newPassword})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Password, updatedUser.Password)
	require.Equal(t, newPassword, updatedUser.Password)

	require.Equal(t, user.Uid, updatedUser.Uid)
	require.Equal(t, user.Fullname, updatedUser.Fullname)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.IsPublic, updatedUser.IsPublic)
	require.Equal(t, user.IsBlocked, updatedUser.IsBlocked)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
}
