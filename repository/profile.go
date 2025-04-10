package repository

import (
	"context"
	"fmt"
	"go-laris/dtos"
	"go-laris/lib"
	"go-laris/models"

	"github.com/jackc/pgx/v5"
)

func FindProfileByUserId(id int) (*models.Profile, error) {
	db := lib.DB()
	defer db.Close(context.Background())

	row := db.QueryRow(
		context.Background(),
		`SELECT id, picture, fullname, province, city, postal_code, gender, country, mobile, address, user_id FROM "profile" WHERE "user_id" = $1 LIMIT 1`, id,
	)

	var profile models.Profile
	err := row.Scan(
		&profile.Id,
		&profile.Picture,
		&profile.FullName,
		&profile.Province,
		&profile.City,
		&profile.PostalCode,
		&profile.Gender,
		&profile.Country,
		&profile.Mobile,
		&profile.Address,
		&profile.UserId,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		fmt.Println("Error:", err)
		return nil, err
	}

	return &profile, nil
}

func UpdateUserProfile(userID int, profile dtos.Profile) (dtos.Profile, error) {
	db := lib.DB()
	defer db.Close(context.Background())

	profile.Province = emptyStringToNil(profile.Province)
	profile.City = emptyStringToNil(profile.City)
	profile.Country = emptyStringToNil(profile.Country)
	profile.Address = emptyStringToNil(profile.Address)

	profile.PostalCode = emptyIntToNil(profile.PostalCode)
	profile.Gender = emptyIntToNil(profile.Gender)
	profile.Mobile = emptyIntToNil(profile.Mobile)

	sql := `
		UPDATE profile
		SET 
			fullname = $1,
			province = $2,
			city = $3,
			postal_code = $4,
			gender = $5,
			country = $6,
			mobile = $7,
			address = $8
		WHERE user_id = $9
		RETURNING id, user_id
	`

	var updatedID int
	var updatedUserID int
	err := db.QueryRow(context.Background(), sql,
		profile.FullName,
		profile.Province,
		profile.City,
		profile.PostalCode,
		profile.Gender,
		profile.Country,
		profile.Mobile,
		profile.Address,
		userID,
	).Scan(&updatedID, &updatedUserID)

	if err != nil {
		return dtos.Profile{}, err
	}

	profile.Id = updatedID
	profile.UserId = updatedUserID

	return profile, nil
}

func emptyStringToNil(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

func emptyIntToNil(i *int) *int {
	if i != nil && *i == 0 {
		return nil
	}
	return i
}

func UpdateProfilePicture(fullName, province, city, postalCode, gender, country, mobile, address, picturePath string, userId int) (dtos.Profile, error) {
	db := lib.DB()
	defer db.Close(context.Background())

	sql := `
		UPDATE "profile"
		SET "fullname" = $1, "province" = $2, "city" = $3, "postal_code" = $4,
			"gender" = $5, "country" = $6, "mobile" = $7, "address" = $8, "picture" = $9
		WHERE "user_id" = $10
		RETURNING id, fullname, province, city, postal_code, gender, country, mobile, address, picture, user_id
	`

	var profile dtos.Profile
	err := db.QueryRow(context.Background(), sql, fullName, province, city, postalCode, gender, country, mobile, address, picturePath, userId).Scan(
		&profile.Id,
		&profile.FullName,
		&profile.Province,
		&profile.City,
		&profile.PostalCode,
		&profile.Gender,
		&profile.Country,
		&profile.Mobile,
		&profile.Address,
		&profile.Picture,
		&profile.UserId,
	)

	if err != nil {
		fmt.Println("Error updating profile:", err)
		return dtos.Profile{}, err
	}

	return profile, nil
}
