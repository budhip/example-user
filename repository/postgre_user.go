package repository

import (
	"log"

	"github.com/budhip/example-user/model"
)

func (pu *postgreRepository) StoreUser(a *model.User) error {
	sqlStatement := `
		INSERT INTO "user" ("id", "first_name", "last_name", "email", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pu.Conn.Exec(sqlStatement, a.ID, a.FirstName, a.LastName, a.Email, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (pu *postgreRepository) fetch(query string, args ...interface{}) ([]*model.User, error) {

	rows, err := pu.Conn.Query(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*model.User, 0)
	for rows.Next() {
		t := new(model.User)
		err = rows.Scan(
			&t.ID,
			&t.FirstName,
			&t.LastName,
			&t.Email,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (pu *postgreRepository) GetUserByID(id int64) (*model.User, error) {
	query := `SELECT "id", first_name, last_name, email, created_at, updated_at FROM "user" WHERE "id" = $1`

	list, err := pu.fetch(query, id)
	if err != nil {
		return nil, err
	}

	var user *model.User

	if len(list) > 0 {
		user = list[0]
	} else {
		return nil, err
	}

	return user, nil
}