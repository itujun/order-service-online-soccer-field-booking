package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	FindAllWithPagination(context.Context, *dto.OrderRequestParam) ([]models.Order, int64, error)
	FindByUserID(context.Context, string) ([]models.Order, error)
	FindByUUID(context.Context, string) (*models.Order, error)
	Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error)
	Update(context.Context, *gorm.DB, *models.Order, uuid.UUID) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) FindAllWithPagination(
	ctx context.Context,
	param *dto.OrderRequestParam,
) ([]models.Order, int64, error) {
	var (
		orders []models.Order
		sort   string
		total  int64
	)
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := o.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&orders).
		Error
	if err != nil {
		return nil, 0, err
	}

	err = o.db.
		WithContext(ctx).
		Model(&models.Order{}).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, total, nil
}

func (o *OrderRepository) FindByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	var order models.Order
	err := o.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&order).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &order, nil
}

func (o *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&orders).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, nil
}

func (o *OrderRepository) incrementCode(ctx context.Context) (*string, error) {
	var (
		order  *models.Order
		result string
		today  = time.Now().Format("20060102")
	)
	err := o.db.WithContext(ctx).Order("id desc").First(&order).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	if order.ID != 0 {
		orderCode := order.Code
		splitOrderName, _ := strconv.Atoi(orderCode[4:9])
		code := splitOrderName + 1
		result = fmt.Sprintf("ORD-%05d-%s", code, today)
	} else {
		result = fmt.Sprintf("ORD-%05d-%s", 1, today)
	}

	return &result, nil
}
