package infrastructure

import (
	"database/sql"
	"ddd-bottomup/domain"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type MySQLCircleRepository struct {
	db *sql.DB
}

func NewMySQLCircleRepository(db *sql.DB) domain.CircleRepository {
	return &MySQLCircleRepository{
		db: db,
	}
}

func (r *MySQLCircleRepository) FindByID(id *domain.CircleID) (*domain.Circle, error) {
	query := `
		SELECT id, name, owner_id, created_at 
		FROM circles 
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id.Value())

	var circleID, name, ownerID string
	var createdAt time.Time

	if err := row.Scan(&circleID, &name, &ownerID, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// エンティティの再構成
	reconstructedID, _ := domain.ReconstructCircleID(circleID)
	circleName, _ := domain.NewCircleName(name)
	reconstructedOwnerID, _ := domain.ReconstructUserID(ownerID)

	// メンバーIDを取得
	memberIDs, err := r.getMemberIDs(reconstructedID)
	if err != nil {
		return nil, err
	}

	return domain.ReconstructCircle(reconstructedID, circleName, reconstructedOwnerID, memberIDs, createdAt), nil
}

func (r *MySQLCircleRepository) FindByName(name *domain.CircleName) (*domain.Circle, error) {
	query := `
		SELECT id, name, owner_id, created_at 
		FROM circles 
		WHERE name = ?
	`

	row := r.db.QueryRow(query, name.Value())

	var circleID, circleName, ownerID string
	var createdAt time.Time

	if err := row.Scan(&circleID, &circleName, &ownerID, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// エンティティの再構成
	reconstructedID, _ := domain.ReconstructCircleID(circleID)
	reconstructedName, _ := domain.NewCircleName(circleName)
	reconstructedOwnerID, _ := domain.ReconstructUserID(ownerID)

	// メンバーIDを取得
	memberIDs, err := r.getMemberIDs(reconstructedID)
	if err != nil {
		return nil, err
	}

	return domain.ReconstructCircle(reconstructedID, reconstructedName, reconstructedOwnerID, memberIDs, createdAt), nil
}

func (r *MySQLCircleRepository) FindAll() ([]*domain.Circle, error) {
	query := `
		SELECT id, name, owner_id, created_at 
		FROM circles 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCircles(rows)
}

func (r *MySQLCircleRepository) Save(circle *domain.Circle) error {
	// トランザクション開始
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// サークル保存（UPSERT）
	query := `
		INSERT INTO circles (id, name, owner_id, created_at, member_count) 
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		name = VALUES(name), 
		owner_id = VALUES(owner_id),
		member_count = VALUES(member_count)
	`

	_, err = tx.Exec(query,
		circle.ID().Value(),
		circle.Name().Value(),
		circle.OwnerID().Value(),
		circle.CreatedAt(),
		circle.GetMemberCount())
	if err != nil {
		return err
	}

	// 既存のメンバー関係を削除
	_, err = tx.Exec("DELETE FROM circle_members WHERE circle_id = ?", circle.ID().Value())
	if err != nil {
		return err
	}

	// 新しいメンバー関係を挿入
	if len(circle.GetMemberIDs()) > 0 {
		memberQuery := "INSERT INTO circle_members (circle_id, user_id) VALUES "
		values := make([]string, len(circle.GetMemberIDs()))
		args := make([]interface{}, 0, len(circle.GetMemberIDs())*2)

		for i, memberID := range circle.GetMemberIDs() {
			values[i] = "(?, ?)"
			args = append(args, circle.ID().Value(), memberID.Value())
		}

		memberQuery += strings.Join(values, ", ")
		_, err = tx.Exec(memberQuery, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLCircleRepository) Delete(id *domain.CircleID) error {
	query := "DELETE FROM circles WHERE id = ?"
	_, err := r.db.Exec(query, id.Value())
	return err
}

// getMemberIDs はサークルのメンバーIDを取得します
func (r *MySQLCircleRepository) getMemberIDs(circleID *domain.CircleID) ([]*domain.UserID, error) {
	query := "SELECT user_id FROM circle_members WHERE circle_id = ?"
	rows, err := r.db.Query(query, circleID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []*domain.UserID
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		memberID, _ := domain.ReconstructUserID(userID)
		memberIDs = append(memberIDs, memberID)
	}

	return memberIDs, rows.Err()
}

// scanCircles は複数のサークルをスキャンします
func (r *MySQLCircleRepository) scanCircles(rows *sql.Rows) ([]*domain.Circle, error) {
	var circles []*domain.Circle

	for rows.Next() {
		var circleID, name, ownerID string
		var createdAt time.Time

		if err := rows.Scan(&circleID, &name, &ownerID, &createdAt); err != nil {
			return nil, err
		}

		// エンティティの再構成
		reconstructedID, _ := domain.ReconstructCircleID(circleID)
		circleName, _ := domain.NewCircleName(name)
		reconstructedOwnerID, _ := domain.ReconstructUserID(ownerID)

		// メンバーIDを取得
		memberIDs, err := r.getMemberIDs(reconstructedID)
		if err != nil {
			return nil, err
		}

		circle := domain.ReconstructCircle(reconstructedID, circleName, reconstructedOwnerID, memberIDs, createdAt)
		circles = append(circles, circle)
	}

	return circles, rows.Err()
}
