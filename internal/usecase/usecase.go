package usecase

//go:generate mockgen -destination=mock/user_usecase_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain UserUseCase
//go:generate mockgen -destination=mock/product_usecase_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain ProductUseCase
//go:generate mockgen -destination=mock/order_usecase_mock.go -package=mock github.com/zercle/gofiber-skeleton/internal/domain OrderUseCase