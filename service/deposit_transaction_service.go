package service

import (
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"os"
	"pojok-baca-api/model"
	"pojok-baca-api/repository"
	"time"
)

type DepositTransactionService interface {
	CreateTransaction(userID uint, amount int) (*snap.Response, error)
	HandleWebhook(orderID, transactionStatus string) error
}

type depositTransactionService struct {
	repo     repository.DepositTransactionRepository
	userRepo repository.UserRepository
}

func NewDepositService(repo repository.DepositTransactionRepository, userRepo repository.UserRepository) DepositTransactionService {
	return &depositTransactionService{repo, userRepo}
}

func (s *depositTransactionService) CreateTransaction(userID uint, amount int) (*snap.Response, error) {
	orderID := fmt.Sprintf("ORDER-%d-%d", userID, time.Now().Unix())

	// Simpan transaksi ke DB
	tx := &model.DepositTransaction{
		UserID:     userID,
		OrderID:    orderID,
		Deposit:    amount,
		Status:     "pending",
		PaymentRef: "", // belum tahu tokennya
	}

	_, err := s.repo.Create(tx)
	if err != nil {
		return nil, err
	}

	// Setup Midtrans
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
	}

	snapRes, midErr := snap.CreateTransaction(snapReq)
	if midErr != nil {
		return nil, fmt.Errorf("midtrans error: %v", midErr.Message)
	}

	return snapRes, nil
}

func (s *depositTransactionService) HandleWebhook(orderID, transactionStatus string) error {
	tx, err := s.repo.GetByOrderID(orderID)
	if err != nil {
		return err
	}

	var paidAt *time.Time

	if transactionStatus == "settlement" || transactionStatus == "capture" {

		now := time.Now()
		paidAt = &now

		user, err := s.userRepo.GetByID(tx.UserID)
		if err != nil {
			return err
		}
		newDeposit := 0
		if user.Deposit != nil {
			newDeposit = *user.Deposit
		}
		total := newDeposit + tx.Deposit
		_, err = s.userRepo.UpdateDeposit(total, tx.UserID)
		if err != nil {
			return err
		}
	}

	return s.repo.UpdateStatus(orderID, transactionStatus, paidAt)
}
