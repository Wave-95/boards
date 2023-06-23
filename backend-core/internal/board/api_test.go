package board

import (
	"context"
	"net/http"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	// Setup API
	boardRepo := NewMockRepository()
	user := test.NewUser()
	validator := validator.New()
	boardService := NewService(boardRepo, validator)
	boardAPI := NewAPI(boardService, validator)
	r := chi.NewRouter()
	jwtService := test.NewJWTService()
	authHandler := middleware.Auth(jwtService)
	boardAPI.RegisterHandlers(r, authHandler)

	// Setup data
	boardRepo.AddUser(user)
	board := test.NewBoard(user.ID)
	err := boardRepo.CreateBoard(context.Background(), board)
	if err != nil {
		assert.FailNow(t, "Failed to generate test board needed for sending board invites")
	}
	receiver1 := test.NewUser()
	receiver2 := test.NewUser()

	// Setup table tests
	token, err := jwtService.GenerateToken(user.ID.String())
	if err != nil {
		assert.FailNow(t, "Failed to generate test token needed for sending authenticated requests")
	}
	authHeader := test.AuthHeader(token)
	tt := []test.APITestCase{
		{"create board", http.MethodPost, "/boards", `{"name":"My very first board"}`, authHeader, http.StatusCreated, "*My very first board*"},
		{"create invites", http.MethodPost, `/boards/` + board.ID.String() + `/invites`, `{"invites":[{"receiver_id":"` + receiver1.ID.String() + `"}, {"receiver_id":"` + receiver2.ID.String() + `"}]}`, authHeader, http.StatusCreated, `*"status":"PENDING"*`},
	}

	for _, tc := range tt {
		test.Endpoint(t, r, tc)
	}
}