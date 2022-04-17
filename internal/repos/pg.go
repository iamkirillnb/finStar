package repos

import (
	"fmt"
	"github.com/iamkirillnb/finStar/internal"
	"github.com/iamkirillnb/finStar/internal/entities"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const defaultDbName string = "wallet"

type Postgres struct {
	*sqlx.DB

	Config *internal.DbConfig
}

func NewPostgres(c *internal.DbConfig) *Postgres {
	db, err := sqlx.Connect("postgres", c.DbUrlConnection())
	if err != nil {
		log.Println("connetction to posgtres failed")
		log.Fatal(err)
	}

	return &Postgres{
		DB:     db,
		Config: c,
	}
}

type DbRepo struct {
	*Postgres
}

func NewDbRepo(p *Postgres) *DbRepo {
	return &DbRepo{p}
}

func (d *DbRepo) SelectAll() []*entities.Wallet{
	var qry = `SELECT id, amount FROM wallet`


	var data []*entities.Wallet
	err := d.Select(&data, qry)
	if err != nil {
		log.Println("select all data from DB failed")
		log.Fatal(err)
	}
	return data
}

func (d *DbRepo) CreateUser(u *entities.Wallet) error {
	const qry = `INSERT INTO wallet (amount) values (:amount) RETURNING id`

	_, err := d.NamedExec(qry, &u)
	if err != nil {
		log.Println("insert to DB failed")
		return  err
	}

	return nil
}


func (d *DbRepo) GetUser(userId string) (*entities.Wallet, error) {
	const qry = `SELECT id, amount FROM wallet WHERE id=$1`
	user := entities.Wallet{}

	err := d.Get(&user, qry, userId)
	if err != nil {
		log.Println("get user failed")
		return nil, err
	}
	return &user, nil
}

func (d *DbRepo) AddMoney(u *entities.Wallet) error {
	user, err := d.GetUser(u.Id)
	if err != nil {
		log.Printf("user %s is not exists", u.Id)
		return err
	}
	u.Amount += user.Amount

	const qry = `UPDATE wallet SET amount=:amount  WHERE id=:id`

	_, err = d.NamedExec(qry, u)
	if err != nil {
		log.Println("add money to user failed")
		return err
	}
	return nil
}


func (d *DbRepo) SendMoneyToUser(from ,to  *entities.Wallet) error {
	from_user, err := d.GetUser(from.Id)
	if err != nil {
		log.Printf("user %s is not exists", from.Id)
		return err
	}
	if from_user.Amount - from.Amount < 0 {
		log.Println("not enough funds in the account")
		return fmt.Errorf("not enough funds in the account")
	} else {
		from_user.Amount -= from.Amount
		_, err := d.NamedExec(`UPDATE wallet SET amount=:amount WHERE id=:id`, from_user)
		if err != nil {
			log.Println("update user sender amount failed")
			return err
		}

	}

	to_user, err := d.GetUser(to.Id)
	if err != nil {
		log.Printf("user %s is not exists", to.Id)
		return err
	}
	to_user.Amount += from.Amount

	_, err = d.NamedExec(`UPDATE wallet SET amount=:amount WHERE id=:id`, to_user)
	if err != nil {
		log.Println("update user sender amount failed")
		return err
	}

	return nil
}