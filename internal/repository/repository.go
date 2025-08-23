package repository

//go:generate mockgen -destination=mock/user_repository_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain UserRepository
//go:generate mockgen -destination=mock/product_repository_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain ProductRepository
//go:generate mockgen -destination=mock/order_repository_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain OrderRepository