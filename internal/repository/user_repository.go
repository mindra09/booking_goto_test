package repository

import (
	"booking_togo/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	GetAll(ctx context.Context) ([]*model.UserDetailResponse, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userID int) error
	DeleteFamily(ctx context.Context, userID int, familyID int) error
	GetUserDetail(ctx context.Context, userID int) (*model.UserDetailResponse, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*model.UserDetailResponse, error) {
	families := model.FamiliesJSON{}
	usersResponse := []*model.UserDetailResponse{}
	var (
		nationalityName, nationalityCode string
	)

	queryStatment := `select 
	 		cust.customer_id as user_id,
      cust.cst_name as name,
      cust.cst_dob as dob,
      cust.nationality_id,
			COALESCE(nat.nationality_name,''),
			COALESCE( nat.nationality_code,''),
			COALESCE(
                (
                    SELECT JSON_AGG(
                        JSON_BUILD_OBJECT(
														'family_id', fl.fl_id::int,
														'user_id', fl.cst_id::int,
                            'name', fl.fl_name,
                            'dob', fl.fl_dob
                        ) ORDER BY fl.fl_id ASC
                    ) 
                    FROM family_list fl 
                    WHERE fl.cst_id = cust.customer_id
										
                ), 
                '[]'::JSON
            ) as families
			from customer cust
			left join nationality nat on cust.nationality_id = nat.nationality_id
			`

	rows, err := r.db.Query(ctx, queryStatment)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.UserDetailResponse
		err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.Dob,
			&user.NationalityID,
			&nationalityName,
			&nationalityCode,
			&families.Families,
		)

		if err != nil {
			return nil, err
		}

		user.Nationality.NationalityName = nationalityName
		user.Nationality.NationalityCode = nationalityCode
		user.Nationality.NationalityID = user.NationalityID

		if len(families.Families) > 0 {
			familyUnmarshallErr := json.Unmarshal(families.Families, &user.Families)
			if familyUnmarshallErr != nil {
				log.Error("GetAll - JSON Unmarshal error: ", familyUnmarshallErr)
				return nil, familyUnmarshallErr
			}
		}

		usersResponse = append(usersResponse, &user)
	}

	return usersResponse, nil

}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	var customerID int
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	userQuery := `INSERT INTO customer (nationality_id, cst_name, cst_dob) 
		VALUES ($1, $2, $3) 
		RETURNING customer_id
	`

	queryRowErr := tx.QueryRow(ctx, userQuery,
		user.NationalityID, user.Name, user.Dob,
	).Scan(&customerID)

	user.UserID = customerID
	if queryRowErr != nil {
		return queryRowErr
	}

	copyCount, copyCountErr := tx.CopyFrom(ctx,
		pgx.Identifier{"family_list"},
		[]string{"cst_id", "fl_name", "fl_dob"},
		pgx.CopyFromSlice(len(user.Families), func(i int) ([]any, error) {
			family := user.Families[i]
			return []any{user.UserID, family.Name, family.Dob}, nil
		}),
	)

	if copyCountErr != nil {
		return fmt.Errorf("failed to copy from: %w", copyCountErr)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Infof("✅ COPY inserted %d family members for user %d", copyCount, user.UserID)

	return nil
}

func (r *UserRepository) GetUserDetail(ctx context.Context, userID int) (*model.UserDetailResponse, error) {
	families := model.FamiliesJSON{}
	UserDetailResponse := model.UserDetailResponse{}
	var (
		nationalityName, nationalityCode string
	)

	queryStatment := `select 
	 		cust.customer_id as user_id,
      cust.cst_name as name,
      cust.cst_dob as dob,
      cust.nationality_id,
			COALESCE(nat.nationality_name,''),
			COALESCE( nat.nationality_code,''),
			COALESCE(
                (
                    SELECT JSON_AGG(
                        JSON_BUILD_OBJECT(
														'family_id', fl.fl_id::int,
														'user_id', fl.cst_id::int,
                            'name', fl.fl_name,
                            'dob', fl.fl_dob
                        ) ORDER BY fl.fl_id ASC
                    ) 
                    FROM family_list fl 
                    WHERE fl.cst_id = cust.customer_id
										
                ), 
                '[]'::JSON
            ) as families
			from customer cust
			left join nationality nat on cust.nationality_id = nat.nationality_id
			where  cust.customer_id = $1
			`

	err := r.db.QueryRow(ctx, queryStatment, userID).Scan(
		&UserDetailResponse.UserID,
		&UserDetailResponse.Name,
		&UserDetailResponse.Dob,
		&UserDetailResponse.NationalityID,
		&nationalityName,
		&nationalityCode,
		&families.Families,
	)

	if err != nil {
		log.Error("GetUserDetail - QueryRow error: ", err)
		return &UserDetailResponse, err
	}

	if len(families.Families) > 0 {
		familyUnmarshallErr := json.Unmarshal(families.Families, &UserDetailResponse.Families)
		if familyUnmarshallErr != nil {
			log.Error("GetUserDetail - JSON Unmarshal error: ", familyUnmarshallErr)
			return &UserDetailResponse, familyUnmarshallErr
		}
	}

	UserDetailResponse.Nationality.NationalityName = nationalityName
	UserDetailResponse.Nationality.NationalityCode = nationalityCode
	UserDetailResponse.Nationality.NationalityID = UserDetailResponse.NationalityID

	return &UserDetailResponse, nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	deleteUserQuery := `DELETE FROM customer WHERE customer_id = $1`
	_, deleteUserErr := tx.Exec(ctx, deleteUserQuery, userID)
	if deleteUserErr != nil {
		return deleteUserErr

	}

	deleteFamilyQuery := `DELETE FROM family_list WHERE cst_id = $1`
	_, deleteFamilyErr := tx.Exec(ctx, deleteFamilyQuery, userID)
	if deleteFamilyErr != nil {
		return deleteFamilyErr
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteFamily(ctx context.Context, userID int, familyID int) error {

	deleteFamilyQuery := `DELETE FROM family_list WHERE cst_id = $1 and fl_id = $2`
	_, deleteFamilyErr := r.db.Exec(ctx, deleteFamilyQuery, userID, familyID)
	if deleteFamilyErr != nil {
		return deleteFamilyErr
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	var (
		customerID int
		updatedAt  time.Time
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	userQuery := `UPDATE customer 
    SET nationality_id = $1, cst_name = $2, cst_dob = $3 , updated_at = CURRENT_TIMESTAMP
    WHERE customer_id = $4
    RETURNING customer_id, updated_at`

	queryRowErr := tx.QueryRow(ctx, userQuery,
		user.NationalityID, user.Name, user.Dob, user.UserID,
	).Scan(&customerID, &updatedAt)

	user.UserID = customerID
	if queryRowErr != nil {
		return queryRowErr
	}

	if err = r.upsertFamilyMembers(ctx, tx, user.Families); err != nil {
		return fmt.Errorf("failed to upsert family members: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}

func (r *UserRepository) upsertFamilyMembers(ctx context.Context, tx pgx.Tx, families []model.Family) error {
	if len(families) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, family := range families {
		query := `
		INSERT INTO family_list (fl_id,cst_id, fl_name, fl_dob) 
		VALUES ( 
			CASE WHEN $1 = 0 THEN nextval('family_list_fl_id_seq') ELSE $1 END,
		$2, $3, $4)
		ON CONFLICT (fl_id) 
		DO UPDATE SET 
				cst_id = EXCLUDED.cst_id,
				fl_name = EXCLUDED.fl_name,
				fl_dob = EXCLUDED.fl_dob`

		batch.Queue(query, family.FamilyID, family.UserID, family.Name, family.Dob)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(families); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to upsert family member %d: %w", i, err)
		}
	}

	log.Infof("✅ upsert user family successfully: %d family members processed", len(families))

	return results.Close()
}
