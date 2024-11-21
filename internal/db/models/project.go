package models

import (
	"errors"
	// "log"
	"context"
	uuid "github.com/google/uuid"
	// "github.com/uptrace/bun"
	"time"
	// "github.com/jackc/pgx/v5/pgxpool"
	"database/sql"
)

type NewProject struct {
	Name     string
	Template string
}

// type Project struct {
// 	*NewProject
// 	ID        uuid.UUID 
// 	Pools     []*Pool   
// 	CreatedAt time.Time
// }

type Project struct {
	Name     string
	Template string
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
	return nil, errors.New("not implemented")
	// ctx := context.Background()
	// id := uuid.New()
	// createdProject := &Project{
	// 	NewProject: req,
	// 	ID:         id,
	// 	CreatedAt:  time.Now(),
	// }
	// _, err := db.NewInsert().Model(createdProject).Exec(ctx)
	// if err != nil {
		// return nil, err
	// }

	// project := &Project{}

	// err = db.NewSelect().
	// 	Model(project).
	// 	Where("project.id = ?", id).
	// 	Scan(ctx)

	// return project, err
}

func GetProject(db *sql.DB, projectId string) (*Project, error) {
	return nil, errors.New("not implemented")
	// ctx := context.Background()
	// project := new(Project)

	// err := db.NewSelect().
	//     Model(project).
	// 	Relation("Pools").
	// 	Where("project.id = ?", projectId).
	// 	Limit(1).
	// 	Scan(ctx)



		// Model(project).
		// Column("project.*").
		// Join("LEFT JOIN pools ON pools.project_id = project.id").
		// // Join("JOIN pools as p ON p.project_id = project.id").
		// // Join("LEFT JOIN pools_attached ON pool.project_id = project.id").
		// Where("project.id = ?", projectId).
		// Scan(ctx)

		// ColumnExpr("project.*").
		// ColumnExpr("pool.*").
		// Join("JOIN pools ON pool.project_id = project.id").

		// Limit(1).
	
	//     log.Printf("%#v", project)

	// return project, err
}
