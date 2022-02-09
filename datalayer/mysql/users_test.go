package mysql

func (suite *DataStoreSuite) TestAddUser() {
	err := suite.Store.AddUser("newuser@example.org", "testpassword")
	suite.NoError(err, "Error inserting new user")
	u, err := suite.Store.GetUserByEmail("newuser@example.org")
	suite.NoError(err, "Error fetching user back")
	suite.NotNil(u, "No user returned")
	suite.NotEqual("testpassword", u.PasswordDigest, "Password was not encrypted")
	suite.NotEmpty(u.PasswordDigest, "No password digest in response")
}
