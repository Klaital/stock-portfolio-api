package mysql

func (suite *DataStoreSuite) TestAddPosition() {
	suite.NoError(suite.Store.AddUser("test@example.org", "nopass"), "Failed to setup user")
	u, err := suite.Store.GetUserByEmail("test@example.org")
	suite.NoError(err, "Failed to fetch test user")
	err = suite.Store.AddPosition(u.ID, "TMUS", 5, 12450, nil)
	suite.NoError(err, "Failed to insert new position")
}
