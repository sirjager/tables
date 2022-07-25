package tokens

// func TestJwtBuilder(t *testing.T) {
// 	builder, err := NewJwtBuilder(utils.RandomString(33))
// 	require.NoError(t, err)

// 	user := utils.RandomUserName()
// 	duration := time.Minute

// 	issued_at := time.Now()
// 	expired_at := issued_at.Add(duration)

// 	token, payload, err := builder.CreateToken(user, duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)
// 	require.NotEmpty(t, payload)

// 	payload, err = builder.VerifyToken(token)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	require.NotZero(t, payload.Id)
// 	require.Equal(t, user, payload.User)
// 	require.WithinDuration(t, issued_at, payload.IssuedAt, time.Second)
// 	require.WithinDuration(t, expired_at, payload.ExpiredAt, time.Second)

// }

// func TestExpiredJwtToken(t *testing.T) {
// 	builder, err := NewJwtBuilder(utils.RandomString(32))
// 	require.NoError(t, err)

// 	token, payload, err := builder.CreateToken(utils.RandomUserName(), -time.Minute)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)
// 	require.NotEmpty(t, payload)

// 	payload, err = builder.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrExpiredToken.Error())
// 	require.Nil(t, payload)

// }

// func TestInvalidJwtTokenAlgNone(t *testing.T) {
// 	payload, err := NewPayload(utils.RandomUserName(), time.Minute)
// 	require.NoError(t, err)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
// 	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 	require.NoError(t, err)

// 	builder, err := NewJwtBuilder(utils.RandomString(32))
// 	require.NoError(t, err)

// 	payload, err = builder.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)
// }
