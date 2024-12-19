package repository

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllPending(ctx context.Context) ([]entity.Product, error)
	GetByIdPending(ctx context.Context, id int64) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
	SearchProducts(ctx context.Context, search string) ([]entity.Product, error)
	FilterProductsByAddress(ctx context.Context, address string) ([]entity.Product, error)
	FilterProductsByCategory(ctx context.Context, category string) ([]entity.Product, error)
	FilterProductsByPrice(ctx context.Context, minPrice string, maxPrice string) ([]entity.Product, error)
	FilterProductsByStatus(ctx context.Context, status string) ([]entity.Product, error)
	FilterProductsByDate(ctx context.Context, date string) ([]entity.Product, error)
	FilterProductsByTime(ctx context.Context, time string) ([]entity.Product, error)
	SortProductByNewest(ctx context.Context) ([]entity.Product, error)
	SortProductByExpensive(ctx context.Context) ([]entity.Product, error)
	SortProductByMostBought(ctx context.Context) ([]entity.Product, error)
	SortProductByCheapest(ctx context.Context) ([]entity.Product, error)
	SortProductByAvailable(ctx context.Context) ([]entity.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAllPending(ctx context.Context) ([]entity.Product, error) {
	result := make([]entity.Product, 0)
	if err := r.db.WithContext(ctx).Where("product_status = ?", "pending").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) GetByIdPending(ctx context.Context, id int64) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ? AND product_status = ?", id, "pending").First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Updates(&product).Error
}

func (r *productRepository) Delete(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Delete(&product).Error
}

func (r *productRepository) GetByVerifySubmissionToken(ctx context.Context, token string) (*entity.Product, error) {
	result := new(entity.Product)
	if err := r.db.WithContext(ctx).Where("verify_submission_token = ?", token).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// search products
func (r *productRepository) SearchProducts(ctx context.Context, search string) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("name LIKE ?", "%"+search+"%").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// filter by addres
func (r *productRepository) FilterProductsByAddress(ctx context.Context, address string) ([]entity.Product, error) // filter byAddress
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_address = ?", address).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *productRepository) FilterProductsByCategory(ctx context.Context, category string) ([]entity.Product, error) // filter byCategory
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_category = ?", category).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
func (r *productRepository) FilterProductsByTime(ctx context.Context, time string) ([]entity.Product, error) // filter byTime
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_time = ?", time).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *productRepository) FilterProductsByDate(ctx context.Context, date string) ([]entity.Product, error) // filter byDate
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_date = ?", date).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}


// filter  by price min max
func (r *productRepository) FilterProductsByPrice(ctx context.Context, minPrice string, maxPrice string) ([]entity.Product, error) // filter byPrice
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *productRepository) FilterProductsByStatus(ctx context.Context, status string) ([]entity.Product, error) // filter byStatus
{
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Where("product_status = ?", status).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

//shorting produk baru
fun (r *productRepository) SortProductByNewest(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Order("created_at DESC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}


//sort by price expensive(termahal)
fun (r *productRepository) SortProductByExpensive(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Order("product_price DESC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

//sort dari yang paling murah
func (r *productRepository) SortProductByCheapest(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Order("product_price ASC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
	}


func (r *productRepository) SortProductByMostBought(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Order("product_views DESC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *productRepository) SortProductByAvailable(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	result := r.db.WithContext(ctx).Order("product_status DESC").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
	}
