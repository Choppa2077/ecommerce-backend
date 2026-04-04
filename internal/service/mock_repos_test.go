package service_test

import (
	"context"

	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// ─── MockUserRepository ───────────────────────────────────────────────────────

type MockUserRepository struct {
	CreateFn          func(ctx context.Context, user *domain.User) error
	GetByEmailFn      func(ctx context.Context, email string) (*domain.User, error)
	GetByIDFn         func(ctx context.Context, id int) (*domain.User, error)
	UpdateFn          func(ctx context.Context, user *domain.User) error
	UpdateLastLoginFn func(ctx context.Context, id int) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.CreateFn(ctx, user)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.GetByEmailFn(ctx, email)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.UpdateFn(ctx, user)
}
func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	return m.UpdateLastLoginFn(ctx, id)
}

// ─── MockProductRepository ────────────────────────────────────────────────────

type MockProductRepository struct {
	CreateFn                  func(ctx context.Context, product *domain.Product) error
	GetByIDFn                 func(ctx context.Context, id int) (*domain.Product, error)
	GetByIDWithCategoryFn     func(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	UpdateFn                  func(ctx context.Context, product *domain.Product) error
	DeleteFn                  func(ctx context.Context, id int) error
	ListFn                    func(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error)
	ListWithCategoriesFn      func(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error)
	SearchFn                  func(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error)
	CreateCategoryFn          func(ctx context.Context, category *domain.Category) error
	GetCategoryByIDFn         func(ctx context.Context, id int) (*domain.Category, error)
	GetCategoryByNameFn       func(ctx context.Context, name string) (*domain.Category, error)
	ListCategoriesFn          func(ctx context.Context) ([]*domain.Category, error)
	UpdateCategoryFn          func(ctx context.Context, category *domain.Category) error
	DeleteCategoryFn          func(ctx context.Context, id int) error
	GetProductStatisticsFn    func(ctx context.Context, productID int) (*domain.ProductStatistics, error)
	RefreshProductStatsFn     func(ctx context.Context) error
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return m.CreateFn(ctx, product)
}
func (m *MockProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *MockProductRepository) GetByIDWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	return m.GetByIDWithCategoryFn(ctx, id)
}
func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return m.UpdateFn(ctx, product)
}
func (m *MockProductRepository) Delete(ctx context.Context, id int) error {
	return m.DeleteFn(ctx, id)
}
func (m *MockProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	return m.ListFn(ctx, filter)
}
func (m *MockProductRepository) ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error) {
	return m.ListWithCategoriesFn(ctx, filter)
}
func (m *MockProductRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error) {
	return m.SearchFn(ctx, query, limit, offset)
}
func (m *MockProductRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	return m.CreateCategoryFn(ctx, category)
}
func (m *MockProductRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	return m.GetCategoryByIDFn(ctx, id)
}
func (m *MockProductRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	return m.GetCategoryByNameFn(ctx, name)
}
func (m *MockProductRepository) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	return m.ListCategoriesFn(ctx)
}
func (m *MockProductRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	return m.UpdateCategoryFn(ctx, category)
}
func (m *MockProductRepository) DeleteCategory(ctx context.Context, id int) error {
	return m.DeleteCategoryFn(ctx, id)
}
func (m *MockProductRepository) GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error) {
	return m.GetProductStatisticsFn(ctx, productID)
}
func (m *MockProductRepository) RefreshProductStatistics(ctx context.Context) error {
	return m.RefreshProductStatsFn(ctx)
}

// ─── MockInteractionRepository ────────────────────────────────────────────────

type MockInteractionRepository struct {
	RecordViewFn                func(ctx context.Context, userID, productID int) error
	GetUserViewsFn              func(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)
	HasViewedFn                 func(ctx context.Context, userID, productID int) (bool, error)
	RecordLikeFn                func(ctx context.Context, userID, productID int) error
	RemoveLikeFn                func(ctx context.Context, userID, productID int) error
	GetUserLikesFn              func(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)
	HasLikedFn                  func(ctx context.Context, userID, productID int) (bool, error)
	RecordPurchaseFn            func(ctx context.Context, userID, productID int, quantity int, price float64) error
	GetUserPurchasesFn          func(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)
	HasPurchasedFn              func(ctx context.Context, userID, productID int) (bool, error)
	GetUserInteractionSummaryFn func(ctx context.Context, userID int) (*domain.UserInteractionSummary, error)
	GetAllUserViewsFn           func(ctx context.Context) ([]domain.UserProductView, error)
	GetAllUserLikesFn           func(ctx context.Context) ([]domain.UserProductLike, error)
	GetAllUserPurchasesFn       func(ctx context.Context) ([]domain.UserProductPurchase, error)
}

func (m *MockInteractionRepository) RecordView(ctx context.Context, userID, productID int) error {
	return m.RecordViewFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) GetUserViews(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	return m.GetUserViewsFn(ctx, userID, limit)
}
func (m *MockInteractionRepository) HasViewed(ctx context.Context, userID, productID int) (bool, error) {
	return m.HasViewedFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) RecordLike(ctx context.Context, userID, productID int) error {
	return m.RecordLikeFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) RemoveLike(ctx context.Context, userID, productID int) error {
	return m.RemoveLikeFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) GetUserLikes(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	return m.GetUserLikesFn(ctx, userID, limit)
}
func (m *MockInteractionRepository) HasLiked(ctx context.Context, userID, productID int) (bool, error) {
	return m.HasLikedFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) RecordPurchase(ctx context.Context, userID, productID int, quantity int, price float64) error {
	return m.RecordPurchaseFn(ctx, userID, productID, quantity, price)
}
func (m *MockInteractionRepository) GetUserPurchases(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	return m.GetUserPurchasesFn(ctx, userID, limit)
}
func (m *MockInteractionRepository) HasPurchased(ctx context.Context, userID, productID int) (bool, error) {
	return m.HasPurchasedFn(ctx, userID, productID)
}
func (m *MockInteractionRepository) GetUserInteractionSummary(ctx context.Context, userID int) (*domain.UserInteractionSummary, error) {
	return m.GetUserInteractionSummaryFn(ctx, userID)
}
func (m *MockInteractionRepository) GetAllUserViews(ctx context.Context) ([]domain.UserProductView, error) {
	return m.GetAllUserViewsFn(ctx)
}
func (m *MockInteractionRepository) GetAllUserLikes(ctx context.Context) ([]domain.UserProductLike, error) {
	return m.GetAllUserLikesFn(ctx)
}
func (m *MockInteractionRepository) GetAllUserPurchases(ctx context.Context) ([]domain.UserProductPurchase, error) {
	return m.GetAllUserPurchasesFn(ctx)
}

// ─── MockProfileRepository ────────────────────────────────────────────────────

type MockProfileRepository struct {
	CreateFn      func(ctx context.Context, profile *domain.Profile) error
	GetByUserIDFn func(ctx context.Context, userID int) (*domain.Profile, error)
	UpdateFn      func(ctx context.Context, profile *domain.Profile) error
	DeleteFn      func(ctx context.Context, userID int) error
}

func (m *MockProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	return m.CreateFn(ctx, profile)
}
func (m *MockProfileRepository) GetByUserID(ctx context.Context, userID int) (*domain.Profile, error) {
	return m.GetByUserIDFn(ctx, userID)
}
func (m *MockProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	return m.UpdateFn(ctx, profile)
}
func (m *MockProfileRepository) Delete(ctx context.Context, userID int) error {
	return m.DeleteFn(ctx, userID)
}
