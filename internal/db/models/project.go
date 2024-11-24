package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	uuid "github.com/google/uuid"
)

type NewProject struct {
	Name     string
	Template string
}

type Project struct {
	Name      string
	Template  string
	ID        uuid.UUID
	Pools     []*Pool
	CreatedAt time.Time
}

type Foo struct {
	Name string
}

func GetProjects(db *sql.DB) ([]*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	projects := make([]*Project, 0)
	rows, err := db.QueryContext(ctx, "SELECT id, name, created_at, template  FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.Template)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, err
}

func CreateProject(db *sql.DB, req *NewProject) (*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := uuid.New()
	insertSqlStatement := "INSERT INTO projects (id, name, created_at, template) VALUES ($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, insertSqlStatement, id, req.Name, time.Now(), req.Template)
	if err != nil {
		return nil, err
	}

	project := &Project{}
	selectSqlStatement := "SELECT id, name, created_at, template FROM projects WHERE id = $1"
	err = db.QueryRowContext(ctx, selectSqlStatement, id).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.Template)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func UpdateProject(db *sql.DB, projectID uuid.UUID, updateReq *NewProject) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updateSqlStatement := `
		UPDATE projects 
		SET name = $1, template = $2 
		WHERE id = $3
	`
	result, err := db.ExecContext(ctx, updateSqlStatement, &updateReq.Name, &updateReq.Template, projectID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("no such row")
	}
	return err
}

func GetProject(db *sql.DB, projectId string) (*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	project := &Project{}
	selectSqlStatement := "SELECT id, name, created_at, template FROM projects WHERE id = $1"
	err := db.QueryRowContext(ctx, selectSqlStatement, projectId).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.Template)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func GetProjectWithPools(db *sql.DB, projectId string) (*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	project := &Project{}
	pools := make([]*Pool, 0)
	selectSqlStatement := `
		SELECT 
		pools.id, pools.name, pools.description, pools.created_at,
		projects.id, projects.name, projects.created_at, projects.template 
		FROM projects 
		LEFT JOIN pools ON pools.project_id = projects.id AND pools.project_id = $1
		ORDER BY pools.created_at
    `
	rows, err := db.QueryContext(ctx, selectSqlStatement, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pool := &Pool{}
		var poolID uuid.UUID
		var poolName sql.NullString
		var poolDescription sql.NullString
		var poolCreatedAt sql.NullTime
		err := rows.Scan(&poolID, &poolName, &poolDescription, &poolCreatedAt, &project.ID, &project.Name, &project.CreatedAt, &project.Template)
		if err != nil {
			return nil, err
		}
		if poolID.String() != "" && poolName.Valid {
			pool.ID = poolID
			pool.Name = poolName.String
			pool.Description = poolDescription.String
			pool.CreatedAt = poolCreatedAt.Time
			pools = append(pools, pool)
		}
	}
	project.Pools = pools

	return project, nil
}
