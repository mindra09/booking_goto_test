package usecase

import (
	"booking_togo/internal/model"
	"booking_togo/internal/repository"
	"context"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IUserUsecase interface {
	GetAll(ctx context.Context) (users []*model.UserDetailResponse, err error)
	Create(ctx context.Context, user *model.User) error
	Detail(ctx context.Context, id int) (user *model.UserDetailResponse, err error)
	Update(ctx context.Context, user *model.User) (err error)
	Delete(ctx context.Context, userID int) (err error)
	DeleteFamily(ctx context.Context, userID int, familyID int) (err error)
}

type UserUsecase struct {
	userRepository repository.IUserRepository
}

func NewUserUsecase(userRepository repository.IUserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) GetAll(ctx context.Context) (users []*model.UserDetailResponse, err error) {
	users, err = u.userRepository.GetAll(ctx)
	if err != nil {
		log.Error("User get all failed: ", err.Error())
		return
	}
	return

}

func (u *UserUsecase) Create(ctx context.Context, user *model.User) (err error) {
	err = u.validateUserFamilies(*user)
	if err != nil {
		log.Error("User validation failed: ", err.Error())
		return err
	}

	for _, v := range user.Families {
		msgErrorName := "Family validation failed for " + v.Name + ": is required and must be between 5 to 50 characters"
		msgErrorDob := "Family validation failed for " + v.Dob + ": is required"

		err = validation.ValidateStruct(&v,
			validation.Field(&v.Name, validation.Required.Error(msgErrorName), validation.Length(5, 50)),
			validation.Field(&v.Dob, validation.Required.Error(msgErrorDob), validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)).Error("Date must be in YYYY-MM-DD format"),
				validation.By(validateDOBFormat)),
		)

		if err != nil {
			log.Error("User - Family validation failed: ", err)
			return err
		}
	}

	if err = u.userRepository.Create(ctx, user); err != nil {
		log.Error("User create failed: ", err.Error())
		return
	}

	return
}

func (u *UserUsecase) Detail(ctx context.Context, id int) (user *model.UserDetailResponse, err error) {
	userDetail, userDetailErr := u.userRepository.GetUserDetail(ctx, id)
	if userDetailErr != nil {
		err = userDetailErr
		return
	}
	return userDetail, nil
}

func (u *UserUsecase) Update(ctx context.Context, user *model.User) (err error) {
	err = u.validateUserFamilies(*user)
	if err != nil {
		log.Error("User validation failed: ", err.Error())
		return err
	}

	for _, v := range user.Families {
		msgErrorName := "Family validation failed for " + v.Name + ": is required and must be between 5 to 50 characters"
		msgErrorDob := "Family validation failed for " + v.Dob + ": is required"

		err = validation.ValidateStruct(&v,
			validation.Field(&v.FamilyID, validation.Min(0).Error("Family ID must be a positive integer or zero for new family record")),
			validation.Field(&v.UserID, validation.Required.Error("User ID is required"), validation.Min(1).Error("User ID must be a positive integer")),
			validation.Field(&v.Name, validation.Required.Error(msgErrorName), validation.Length(5, 50)),
			validation.Field(&v.Dob, validation.Required.Error(msgErrorDob), validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)).Error("Date must be in YYYY-MM-DD format"),
				validation.By(validateDOBFormat)),
		)

		if err != nil {
			log.Error("User - Family validation failed: ", err)
			return err
		}
	}

	if err = u.userRepository.Update(ctx, user); err != nil {
		log.Error("User - Family Update failed: ", err)
		return err
	}

	return
}

func (u *UserUsecase) Delete(ctx context.Context, userID int) (err error) {
	if err = u.userRepository.Delete(ctx, userID); err != nil {
		log.Error("User - Family Delete failed: ", err)
		return
	}

	return
}

func (u *UserUsecase) DeleteFamily(ctx context.Context, userID int, familyID int) (err error) {
	if err = u.userRepository.DeleteFamily(ctx, userID, familyID); err != nil {
		log.Error("Family Delete failed: ", err)
		return
	}

	return
}

func (u *UserUsecase) validateUserFamilies(user model.User) error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&user.Dob, validation.Required, validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)).Error("Date must be in YYYY-MM-DD format"),
			validation.By(validateDOBFormat)),
		validation.Field(&user.NationalityID, validation.Required),
	)

	return err
}

func validateDOBFormat(value interface{}) error {
	dob, ok := value.(string)
	if !ok {
		return validation.NewError("validation_invalid_dob", "DOB must be a string")
	}

	// Parse the date to ensure it's valid
	_, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return validation.NewError("validation_invalid_date", "Invalid date format. Use YYYY-MM-DD")
	}

	return nil
}
