package handlers

import (
	"github.com/iamkirillnb/finStar/internal"
	"github.com/iamkirillnb/finStar/internal/entities"
	"github.com/iamkirillnb/finStar/internal/repos"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

const (
	defaultServerReadTimeout  = 15 * time.Second
	defaultServerWriteTimeout = 30 * time.Second
)

type Handler struct {
	*http.Server

	router *echo.Echo
	config *internal.ServerConfig
	repo   *repos.DbRepo
}

func NewHandler(cfg *internal.ServerConfig, repo *repos.DbRepo) Handler {
	addr := cfg.Address()
	serv := &http.Server{
		Addr:         addr,
		ReadTimeout:  defaultServerReadTimeout,
		WriteTimeout: defaultServerWriteTimeout,
	}
	rout := echo.New()

	return Handler{
		router: rout,
		Server: serv,
		config: cfg,
		repo:   repo,
	}
}

func (h *Handler) GetAllUsers(ctx echo.Context) error {
	users := h.repo.SelectAll()
	return ctx.JSON(200, users)
}

func (h *Handler) CreateUser(ctx echo.Context) error {
	user := &entities.Wallet{}
	err := ctx.Bind(user)
	if err != nil {
		return err
	}
	err = h.repo.CreateUser(user)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}
	return ctx.JSON(http.StatusCreated, "Created")
}

func (h *Handler) GetUser(ctx echo.Context) error {
	user := &entities.Wallet{}
	err := ctx.Bind(user)

	data, err := h.repo.GetUser(user.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "user is not exist")
	}
	return ctx.JSON(http.StatusOK, data)

}

func (h *Handler) AddMoney(ctx echo.Context) error {
	user := &entities.Wallet{}
	err := ctx.Bind(user)
	if err != nil {
		return err
	}
	err = h.repo.AddMoney(user)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}
	return ctx.JSON(http.StatusOK, "Added")

}

func (h *Handler) SendMoneyToUser(ctx echo.Context) error {
	users := []*entities.Wallet{}
	err := ctx.Bind(&users)
	if err != nil {
		return err
	}
	user1, user2 := users[0], users[1]
	err = h.repo.SendMoneyToUser(user1, user2)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}
	return ctx.JSON(http.StatusOK, "Money sent successfully")
}

func (h *Handler) Start() {
	h.router.GET("/all", h.GetAllUsers)
	h.router.POST("/create_user", h.CreateUser)
	h.router.POST("/get_user", h.GetUser)
	h.router.POST("/add_money", h.AddMoney)
	h.router.POST("/send_money", h.SendMoneyToUser)

	log.Fatal(h.router.StartServer(h.Server))
}
