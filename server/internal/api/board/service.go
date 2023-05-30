package board

import (
	"context"
	"fmt"
	"time"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
)

var (
	defaultBoardDescription = "This is a default description for the board. Feel free to customize it and add relevant information about the purpose, goals, or any specific details related to the board."
)

type Service interface {
	CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error)
	GetBoard(ctx context.Context, boardId string) (models.Board, error)
	ListOwnedBoardsWithMembers(ctx context.Context, userId string) ([]OwnedBoardWithMembersDTO, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func (s *service) CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error) {
	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to parse user ID input into UUID: %w", err)
	}
	// create board name if none provided
	if input.Name == nil {
		boards, err := s.repo.ListOwnedBoardAndUsers(ctx, userId)
		if err != nil {
			return models.Board{}, fmt.Errorf("service: failed to get existing boards when creating board: %w", err)
		}
		numBoards := len(boards)
		boardName := fmt.Sprintf("Board #%d", numBoards+1)
		input.Name = &boardName
	}

	// use default board description if none provided
	if input.Description == nil {
		input.Description = &defaultBoardDescription
	}

	// create new board
	id := uuid.New()
	now := time.Now()
	board := models.Board{
		Id:          id,
		Name:        input.Name,
		Description: input.Description,
		UserId:      userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = s.repo.CreateBoard(ctx, board)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to create board: %w", err)
	}
	return board, nil
}

func (s *service) GetBoard(ctx context.Context, boardId string) (models.Board, error) {
	boardIdUUID, err := uuid.Parse(boardId)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: issue parsing boardId into UUID: %w", err)
	}
	board, err := s.repo.GetBoard(ctx, boardIdUUID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to get board: %w", err)
	}
	return board, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListOwnedBoards(ctx context.Context, userId string) ([]models.Board, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	boards, err := s.repo.ListOwnedBoards(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	return boards, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListOwnedBoardsWithMembers(ctx context.Context, userId string) ([]OwnedBoardWithMembersDTO, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	rows, err := s.repo.ListOwnedBoardAndUsers(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// reformat rows into nested structure
	list := []OwnedBoardWithMembersDTO{}
	boardMap := make(map[uuid.UUID]int)
	for _, row := range rows {
		_, ok := boardMap[row.Board.Id]
		if !ok {
			// If board does not exist in slice, then append board into slice. Before append,
			// must convert sqlc storage type into domain type
			boardMap[row.Board.Id] = len(list)
			newItem := OwnedBoardWithMembersDTO{
				Id:          row.Board.Id,
				Name:        row.Board.Name,
				Description: row.Board.Description,
				UserId:      row.Board.UserId,
				Members:     []BoardMemberDTO{},
				CreatedAt:   row.Board.CreatedAt,
				UpdatedAt:   row.Board.UpdatedAt,
			}
			list = append(list, newItem)
		}
		// If user and board membership record exists, append to board members field
		if row.BoardMembership != nil && row.User != nil {
			newBoardMember := BoardMemberDTO{
				Id:    row.User.Id,
				Name:  row.User.Name,
				Email: row.User.Email,
				Membership: MembershipDTO{
					Role:      string(row.BoardMembership.Role),
					CreatedAt: row.BoardMembership.CreatedAt,
					UpdatedAt: row.BoardMembership.UpdatedAt,
				},
			}
			sliceIndex := boardMap[row.Board.Id]
			itemWithNewMember := list[sliceIndex]
			itemWithNewMember.Members = append(itemWithNewMember.Members, newBoardMember)
			list[sliceIndex] = itemWithNewMember
		}
	}
	return list, nil

}

func NewService(repo Repository, validator validator.Validate) *service {
	return &service{repo: repo, validator: validator}
}
