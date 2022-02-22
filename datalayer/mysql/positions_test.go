package mysql

import (
	"github.com/klaital/stock-portfolio-api/datalayer"
	"github.com/shopspring/decimal"
	"time"
)

func (suite *DataStoreSuite) TestAddPosition() {
	suite.NoError(suite.Store.AddUser("test@example.org", "nopass"), "Failed to setup user")
	u, err := suite.Store.GetUserByEmail("test@example.org")
	suite.NoError(err, "Failed to fetch test user")
	err = suite.Store.AddPosition(u.ID, "TMUS", decimal.NewFromInt(5), decimal.NewFromInt(12450), time.Now())
	suite.NoError(err, "Failed to insert new position")
}

func (suite *DataStoreSuite) TestGetPositionsBySymbol() {
	var err error
	var positions []datalayer.Position
	var u *datalayer.User
	suite.NoError(suite.Store.AddUser("test@example.org", "nopass"), "Failed to setup user")
	u, err = suite.Store.GetUserByEmail("test@example.org")
	suite.NoError(err, "Failed to fetch test user")
	suite.Store.AddPosition(u.ID, "TMUS", decimal.NewFromInt(5), decimal.NewFromInt(12450), time.Now())
	suite.Store.AddPosition(u.ID, "GOOG", decimal.NewFromInt(15), decimal.NewFromInt(12650), time.Now())
	suite.Store.AddPosition(u.ID, "GOOG", decimal.NewFromInt(25), decimal.NewFromInt(12750), time.Now())

	positions, err = suite.Store.GetPositionsBySymbol(u.ID, "TMUS")
	suite.NoError(err, "Failed to fetch positions")
	suite.NotEmpty(positions, "No positions returned")
	suite.Equal(1, len(positions))

	positions, err = suite.Store.GetPositionsBySymbol(u.ID, "GOOG")
	suite.NoError(err, "Failed to fetch positions")
	suite.NotEmpty(positions, "No positions returned")
	suite.Equal(2, len(positions))

	positions, err = suite.Store.GetPositionsBySymbol(u.ID+1, "TMUS")
	suite.NoError(err, "Failed to fetch positions")
	suite.Empty(positions, "Should be no positions for invalid user")
}

func (suite *DataStoreSuite) TestGetPositionsByUser() {
	var err error
	var positions []datalayer.Position
	var u *datalayer.User
	suite.NoError(suite.Store.AddUser("test@example.org", "nopass"), "Failed to setup user")
	u, err = suite.Store.GetUserByEmail("test@example.org")
	suite.NoError(err, "Failed to fetch test user")
	suite.Store.AddPosition(u.ID, "TMUS", decimal.NewFromInt(5), decimal.NewFromInt(12450), time.Now())
	suite.Store.AddPosition(u.ID, "GOOG", decimal.NewFromInt(15), decimal.NewFromInt(12650), time.Now())
	suite.Store.AddPosition(u.ID, "GOOG", decimal.NewFromInt(25), decimal.NewFromInt(12750), time.Now())

	positions, err = suite.Store.GetPositionsByUser(u.ID)
	suite.NoError(err, "Failed to fetch positions")
	suite.NotEmpty(positions, "No positions returned")
	suite.Equal(3, len(positions))

	positions, err = suite.Store.GetPositionsByUser(u.ID + 1)
	suite.NoError(err, "Failed to fetch positions")
	suite.Empty(positions, "Should be no positions for invalid user")
}
