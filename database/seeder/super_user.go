package seeder

import (
	"context"
	"fmt"
	"os"

	"hris/internal/user/entity"
	"hris/internal/user/repository"
	"hris/shared/constants"
	"hris/shared/helper"
)

// Email: superuser@example.com
//    •  Password: SuperUser123!
//    •  Role: SUPER_USER
//    •  Status: Active

type SuperUserSeeder struct {
	repo repository.UserRepository
}

func NewSuperUserSeeder(repo repository.UserRepository) *SuperUserSeeder {
	return &SuperUserSeeder{repo: repo}
}

func (s *SuperUserSeeder) Seed(ctx context.Context) error {
	// Check if super user already exists
	existingUser, err := s.repo.FindByEmail(ctx, getSuperUserEmail())
	if err != nil && err != repository.ErrUserNotFound {
		return fmt.Errorf("failed to check existing super user: %w", err)
	}   

	if existingUser != nil {
		fmt.Println("✓ Super user already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(getSuperUserPassword())
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create super user with default company
	superUser := entity.NewUser(
		getSuperUserName(),
		getSuperUserEmail(),
		hashedPassword,
		constants.RoleSuperUser,
		"00000000-0000-0000-0000-000000000001", // Default company
	)
	superUser.IsActive = true

	// Save to database
	if err := s.repo.Create(ctx, superUser); err != nil {
		return fmt.Errorf("failed to create super user: %w", err)
	}

	fmt.Println("✓ Super user created successfully")
	fmt.Printf("  Email: %s\n", superUser.Email)
	fmt.Printf("  Password: %s\n", getSuperUserPassword())
	fmt.Printf("  Role: %s\n", superUser.Role)

	return nil
}

func getSuperUserName() string {
	name := os.Getenv("SUPER_USER_NAME")
	if name == "" {
		return "Super User"
	}
	return name
}

func getSuperUserEmail() string {
	email := os.Getenv("SUPER_USER_EMAIL")
	if email == "" {
		return "superuser@example.com"
	}
	return email
}

func getSuperUserPassword() string {
	password := os.Getenv("SUPER_USER_PASSWORD")
	if password == "" {
		return "SuperUser123!"
	}
	return password
}
