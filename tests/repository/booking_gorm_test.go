package repository_test

import (
	"testing"
	"time"

	"gofiber-skeleton/internal/domain"
	"gofiber-skeleton/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	err = db.AutoMigrate(&domain.Booking{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func TestCreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookingRepository(db)

	start := time.Now().Truncate(time.Second)
	end := start.Add(2 * time.Hour)

	booking := &domain.Booking{
		UserID:    1,
		ItemID:    2,
		StartDate: start,
		EndDate:   end,
	}

	err := repo.Create(booking)
	assert.NoError(t, err, "Create should not return error")
	assert.NotZero(t, booking.ID, "Created booking should have non-zero ID")

	fetched, err := repo.GetByID(booking.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, booking.ID, fetched.ID, "Fetched booking ID should match")
	assert.Equal(t, booking.UserID, fetched.UserID, "Fetched booking UserID should match")
	assert.Equal(t, booking.ItemID, fetched.ItemID, "Fetched booking ItemID should match")
	assert.True(t, fetched.StartDate.Equal(booking.StartDate), "Fetched StartDate should match")
	assert.True(t, fetched.EndDate.Equal(booking.EndDate), "Fetched EndDate should match")
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookingRepository(db)

	booking := &domain.Booking{
		UserID:    3,
		ItemID:    4,
		StartDate: time.Now().Truncate(time.Second),
		EndDate:   time.Now().Add(time.Hour).Truncate(time.Second),
	}
	err := repo.Create(booking)
	assert.NoError(t, err)
	assert.NotZero(t, booking.ID)

	booking.ItemID = 99
	err = repo.Update(booking)
	assert.NoError(t, err, "Update should not return error")

	updated, err := repo.GetByID(booking.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint(99), updated.ItemID, "ItemID should be updated")
}

func TestListByUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookingRepository(db)

	userID := uint(42)
	for i := 0; i < 3; i++ {
		b := &domain.Booking{
			UserID:    userID,
			ItemID:    uint(i + 1),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour),
		}
		err := repo.Create(b)
		assert.NoError(t, err)
	}
	other := &domain.Booking{
		UserID:    100,
		ItemID:    7,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	}
	err := repo.Create(other)
	assert.NoError(t, err)

	list, err := repo.ListByUser(userID)
	assert.NoError(t, err, "ListByUser should not return error")
	assert.Len(t, list, 3, "should return three bookings")
	for _, b := range list {
		assert.Equal(t, userID, b.UserID, "each booking should belong to the requested user")
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookingRepository(db)

	booking := &domain.Booking{
		UserID:    5,
		ItemID:    6,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	}
	err := repo.Create(booking)
	assert.NoError(t, err)

	err = repo.Delete(booking.ID)
	assert.NoError(t, err, "Delete should not return error")

	_, err = repo.GetByID(booking.ID)
	assert.Error(t, err, "GetByID should return error for deleted booking")
}