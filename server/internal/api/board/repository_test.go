package board

import (
	"context"
	"testing"
	"time"

	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)

	userRepo := user.NewRepository(db)
	testUser := setUpTestUser(t, userRepo)
	boardRepo := NewRepository(db)

	t.Cleanup(func() {
		cleanUpTestUser(t, userRepo, testUser.Id)
		db.Close()
	})

	t.Run("Create, get, and delete a board", func(t *testing.T) {
		testBoard := test.NewBoard(testUser.Id)
		// create board
		assert.NoError(t, boardRepo.CreateBoard(context.Background(), testBoard))

		// get board
		board, err := boardRepo.GetBoard(context.Background(), testBoard.Id)
		assert.NoError(t, err)
		assert.Equal(t, testBoard.UserId, board.UserId)

		// delete board
		assert.NoError(t, boardRepo.DeleteBoard(context.Background(), testBoard.Id))
	})

	t.Run("Get board that does not exist", func(t *testing.T) {
		randUUID := uuid.New()
		board, err := boardRepo.GetBoard(context.Background(), randUUID)
		assert.Empty(t, board)
		assert.ErrorIs(t, err, ErrBoardDoesNotExist)
	})

	t.Run("Add a user to a board", func(t *testing.T) {
		testBoard := test.NewBoard(testUser.Id)
		// insert user into new board
		membership := models.BoardMembership{
			Id:        uuid.New(),
			BoardId:   testBoard.Id,
			UserId:    testUser.Id,
			Role:      models.RoleMember,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		assert.NoError(t, boardRepo.CreateBoard(context.Background(), testBoard))
		assert.NoError(t, boardRepo.CreateMembership(context.Background(), membership))

		//TODO: check for insterted user

		// delete board
		err := boardRepo.DeleteBoard(context.Background(), testBoard.Id)
		assert.NoError(t, err)
	})

	t.Run("List boards by user", func(t *testing.T) {
		t.Run("user does not have any boards", func(t *testing.T) {
			boards, err := boardRepo.ListOwnedBoardAndUsers(context.Background(), testUser.Id)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(boards))
		})

		// t.Run("user has 5 boards", func(t *testing.T) {
		// 	boardsToCreate := 5
		// 	for i := 0; i < boardsToCreate; i++ {
		// 		err := boardRepo.CreateBoard(context.Background(), NewTestBoard(testUser.Id))
		// 		assert.NoError(t, err, "expected to create 5 test boards")
		// 	}
		// 	boards, err := boardRepo.ListOwnedBoardAndUsers(context.Background(), testUser.Id)
		// 	assert.NoError(t, err)
		// 	assert.Equal(t, boardsToCreate, len(boards))
		// 	for _, board := range boards {
		// 		fmt.Println(board.Id)
		// 		assert.NoError(t, boardRepo.DeleteBoard(context.Background(), board.Id))
		// 	}
		// })
	})
}

func setUpTestUser(t *testing.T, userRepo user.Repository) models.User {
	testUser := user.NewTestUser()
	err := userRepo.CreateUser(context.Background(), testUser)
	if err != nil {
		assert.FailNow(t, "Could not set up test user for board testing")
	}
	return testUser
}

func cleanUpTestUser(t *testing.T, userRepo user.Repository, userId uuid.UUID) {
	err := userRepo.DeleteUser(userId)
	if err != nil {
		assert.FailNow(t, "Could not clean up test user", err)
	}
}
