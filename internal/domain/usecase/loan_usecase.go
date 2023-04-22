package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	"github.com/LuigiAzevedo/public-library-v2/internal/errs"
	r "github.com/LuigiAzevedo/public-library-v2/internal/ports/repository"
	u "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type loanUseCase struct {
	loanRepo r.LoanRepository
	userRepo r.UserRepository
	bookRepo r.BookRepository
}

// NewLoanUseCase creates a new instance of loanUseCase
func NewLoanUseCase(loan r.LoanRepository, user r.UserRepository, book r.BookRepository) u.LoanUsecase {
	return &loanUseCase{
		loanRepo: loan,
		userRepo: user,
		bookRepo: book,
	}
}

func (s *loanUseCase) BorrowBook(ctx context.Context, userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(ctx, userID, bookID)
	if err != nil {
		return errors.Wrap(err, errs.ErrBorrowBook)
	}
	if exists {
		return errors.New("return the book first before borrowing it again")
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return errors.Wrap(err, errs.ErrGetUser)
	}

	book, err := s.bookRepo.Get(ctx, bookID)
	if err != nil {
		return errors.Wrap(err, errs.ErrGetBook)
	}

	book.Amount -= 1
	if book.Amount < 0 {
		return errors.New("book unavailable at the moment")
	}

	err = s.loanRepo.BorrowTransaction(ctx, user, book)
	if err != nil {
		return errors.Wrap(err, errs.ErrBorrowBook)
	}

	return nil
}

func (s *loanUseCase) ReturnBook(ctx context.Context, userID, bookID int) error {
	exists, err := s.loanRepo.CheckNotReturned(ctx, userID, bookID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("loan does't exists or already returned")
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return errors.Wrap(err, errs.ErrGetUser)
	}

	book, err := s.bookRepo.Get(ctx, bookID)
	if err != nil {
		return errors.Wrap(err, errs.ErrGetBook)
	}

	book.Amount += 1

	err = s.loanRepo.ReturnTransaction(ctx, user, book)
	if err != nil {
		return errors.Wrap(err, "an error occurred while returning a book")
	}

	return nil
}

func (s *loanUseCase) SearchUserLoans(ctx context.Context, userID int) ([]*entity.Loan, error) {
	loans, err := s.loanRepo.Search(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "an error occurred while searching loans")
	}

	return loans, nil
}
