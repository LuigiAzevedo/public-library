package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	"github.com/LuigiAzevedo/public-library-v2/internal/mock"
)

func TestBorrowBook(t *testing.T) {
	repoL := mock.NewMockLoanRepository()
	repoU := mock.NewMockUserRepository()
	repoB := mock.NewMockBookRepository()

	uc := NewLoanUseCase(repoL, repoU, repoB)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		err := uc.BorrowBook(ctx, 1, 1)
		assert.NoError(t, err)
	})
	t.Run("Book Unavailable", func(t *testing.T) {
		err := uc.BorrowBook(ctx, 1, 2)
		assert.Equal(t, errs.ErrBookUnavailable, err.Error())
	})
	t.Run("Wrong ID", func(t *testing.T) {
		err := uc.ReturnBook(ctx, 5, 5)
		assert.Equal(t, errs.ErrLoanAlreadyReturned, err.Error())
	})
}

func TestReturnBook(t *testing.T) {
	repoL := mock.NewMockLoanRepository()
	repoU := mock.NewMockUserRepository()
	repoB := mock.NewMockBookRepository()

	uc := NewLoanUseCase(repoL, repoU, repoB)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		err := uc.ReturnBook(ctx, 2, 2)
		assert.NoError(t, err)
	})
	t.Run("Already Returned", func(t *testing.T) {
		err := uc.ReturnBook(ctx, 1, 1)
		assert.Equal(t, errs.ErrLoanAlreadyReturned, err.Error())
	})
	t.Run("Wrong ID", func(t *testing.T) {
		err := uc.ReturnBook(ctx, 5, 5)
		assert.Equal(t, errs.ErrLoanAlreadyReturned, err.Error())
	})
}

func TestSearchUserLoans(t *testing.T) {
	repoL := mock.NewMockLoanRepository()
	repoU := mock.NewMockUserRepository()
	repoB := mock.NewMockBookRepository()

	uc := NewLoanUseCase(repoL, repoU, repoB)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		loans, err := uc.SearchUserLoans(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, loans, 2)
	})
	t.Run("OK2", func(t *testing.T) {
		loans, err := uc.SearchUserLoans(ctx, 2)
		assert.NoError(t, err)
		assert.Len(t, loans, 1)
	})
	t.Run("Not Found", func(t *testing.T) {
		loans, err := uc.SearchUserLoans(ctx, 5)
		assert.Equal(t, errs.ErrNoLoansFound, errors.Unwrap(err).Error())
		assert.Nil(t, loans)
	})
}
