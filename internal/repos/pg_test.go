package repos

import (
	"github.com/iamkirillnb/finStar/internal"
	"github.com/iamkirillnb/finStar/internal/entities"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"testing"
)

type walletSuite struct {
	suite.Suite

	db   *sqlx.DB
	repo *DbRepo
}

func TestWallteSuit(t *testing.T) {
	suite.Run(t, new(walletSuite))
}

func (w *walletSuite) SetupSuite() {
	cfg := internal.GetConfig("../../dev.yaml")

	db := NewPostgres(&cfg.Db)

	repo := NewDbRepo(db)

	w.db = db.DB
	w.repo = repo
}

func FirstUser() *entities.Wallet {
	// TODO: создавать в базе
	return &entities.Wallet{
		Id: "a9b41fd8-6cd2-451a-9b32-1a0b92970f0e",
		Amount: 1000,
	}
}

func SecondUser() *entities.Wallet {
	// TODO: создавать в базе
	return &entities.Wallet{
		Id: "87a03be9-1fe9-472b-bcd0-3a172ecd236d",
		Amount: 1000,
	}
}

func (w *walletSuite) TestCreateUser() {
	one := FirstUser()
	err := w.repo.CreateUser(one)
	w.Require().NoError(err)
}

func (w *walletSuite) TestGetUser() {
	user := FirstUser()
	result, err := w.repo.GetUser(user.Id)
	w.Require().NoError(err)
	w.Assert().Equal(user.Id, result.Id)
	w.Assert().Equal(user.Amount, result.Amount)
}

func (w *walletSuite) TestAddMoney() {
	beforeAdd := FirstUser()
	afterAdd := FirstUser()
	err := w.repo.AddMoney(afterAdd)
	w.Require().NoError(err)
	w.Assert().NotEqual(beforeAdd.Amount, afterAdd.Amount)
	w.Assert().Equal(beforeAdd.Id, afterAdd.Id)
}

func (w *walletSuite) TestSendMoneyToUser() {
	user1 := FirstUser()
	user1.Amount = 5

	user2 := SecondUser()

	err := w.repo.SendMoneyToUser(user1, user2)
	w.Require().NoError(err)
}