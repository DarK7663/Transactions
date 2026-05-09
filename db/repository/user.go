package repository

import (
	"Transaction/db/model"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MaxTransferAmount = 100_000
	MinTransferAmount = 1
)

var (
	ErrInsufficientFunds = errors.New("недостаточно средств")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrInvalidAmount     = errors.New("некорректная сумма перевода")
	ErrSelfTransfer      = errors.New("нельзя перевести деньги самому себе")
)

type TransferRequest struct {
	SenderID    uint
	RecipientID uint
	Amount      int64
	Reference   string
	Description string
}

type TransferHandeler struct {
	repo model.Transaction
}

type TransactionRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTaskRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) CreateUser(name, email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("failde to hash password: %w", err)
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Balance:  0,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *TransactionRepository) SearchUser(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		fmt.Println("User not found")
		return nil, err
	}
	return &user, nil
}

func (r *TransactionRepository) UpdateUser(id uint, name, email string) (*model.User, error) {
	var userUp model.User

	if _, err := r.SearchUser(id); err != nil {
		return nil, err
	}

	if err := r.db.Where("id = ?", id).Update("name", name).Error; err != nil {
		return nil, err
	}
	if err := r.db.Update("email", email).Error; err != nil {
		return nil, err
	} // костыль
	return &userUp, nil
}
func (r *TransactionRepository) DeleteUser(id uint) error {
	if _, err := r.SearchUser(id); err != nil {
		return err
	}

	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) AuthenticateUser(email, password string) (*model.User, error) {

	var ErrInvalidCredentials = errors.New("неверный email или пароль")

	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}
	var user model.User

	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		r.logger.Error("db error during auth", "err", err, "email", email)
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	fmt.Println("Access open")
	return &user, nil
}

func (r *TransactionRepository) RegisterUser(name, email, password string) (*model.User, error) {

	if name == "" || email == "" || password == "" {
		return nil, errors.New("все поля должны быть заполнены")
	}

	var user = &model.User{
		Name:  name,
		Email: email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *TransactionRepository) SendMoney(ctx context.Context, req TransferRequest) (*model.Transaction, error) {

	if err := ValidateTransferRequest(req); err != nil {
		return nil, err
	}

	var resultTx *model.Transaction

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var existing model.Transaction

		if err := tx.Where("reference = ?", req.Reference).Take(&existing).Error; err != nil { // возврат перевода
			resultTx = &existing
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check idempotency: %w", err)
		}

		var sender model.User

		// блокировка отправителя
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.SenderID).Take(&sender).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("lock sender: %w", err)
		}

		var recipient model.User

		// блокировка получателя
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.RecipientID).Take(&recipient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("lock recipient: %w", err)
		}

		if sender.Balance < req.Amount {
			return ErrInsufficientFunds
		}

		sender.Balance -= req.Amount

		if err := tx.Save(&sender).Error; err != nil {
			return fmt.Errorf("update sender balance: %w", err)
		}

		recipient.Balance += req.Amount
		if err := tx.Save(&recipient).Error; err != nil {
			return fmt.Errorf("update recipient balance: %w", err)
		}

		transaction := model.Transaction{
			SenderID:    req.SenderID,
			RecipientID: req.RecipientID,
			Amount:      req.Amount,
			Status:      "completed",
			Reference:   req.Reference,
			Description: req.Description,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("create transaction record: %w", err)
		}
		resultTx = &transaction
		return nil

	})
	if err != nil {
		r.logger.Error("transfer failed",
			"reference", req.Reference,
			"sender_id", req.SenderID,
			"err", err,
		)

		if errors.Is(err, ErrInsufficientFunds) || errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidAmount) {
			return nil, err
		}
		return nil, errors.New("Ошибка обработки переводов")
	}

	r.logger.Info("transfer completed",
		"reference", req.Reference,
		"amount", req.Amount,
		"sender", req.SenderID,
		"recipient", req.RecipientID,
	)

	if req.Amount > 50_000 {
		r.logger.Warn("large transfer detected",
			"reference", req.Reference,
			"amount", req.Amount,
			"sender", req.SenderID,
		)
	}

	return resultTx, nil
}

func ValidateTransferRequest(req TransferRequest) error {
	if req.SenderID == 0 || req.RecipientID == 0 {
		return ErrUserNotFound
	}
	if req.SenderID == req.RecipientID {
		return ErrSelfTransfer
	}
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}
	if req.Reference == "" {
		return errors.New("reference required for idempotency")
	}
	if req.Amount > MaxTransferAmount || req.Amount < MinTransferAmount {
		return fmt.Errorf("сумма должна быть от %d до %d", MinTransferAmount, MaxTransferAmount)
	}
	return nil
}
