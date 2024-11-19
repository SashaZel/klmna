package models

import (
	"log"
	"context"
	uuid "github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type NewProject struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type Project struct {
	*NewProject
	ID        uuid.UUID `bun:",pk" json:"id"`
	Pools     []*Pool   `bun:"rel:has-many,join:id=project_id" json:"pools"`
	CreatedAt time.Time `json:"created_at"`
}

func GetProjects(db *bun.DB) ([]*Project, error) {
	ctx := context.Background()
	projects := make([]*Project, 0)

	err := db.NewSelect().Model(&projects).Scan(ctx)

	return projects, err
}

func CreateProject(db *bun.DB, req *NewProject) (*Project, error) {
	ctx := context.Background()
	id := uuid.New()
	createdProject := &Project{
		NewProject: req,
		ID:         id,
		CreatedAt:  time.Now(),
	}
	_, err := db.NewInsert().Model(createdProject).Exec(ctx)
	if err != nil {
		return nil, err
	}

	project := &Project{}

	err = db.NewSelect().
		Model(project).
		Where("project.id = ?", id).
		Scan(ctx)

	return project, err
}

func GetProject(db *bun.DB, projectId string) (*Project, error) {
	ctx := context.Background()
	project := new(Project)

	err := db.NewSelect().
	    Model(project).
		Relation("Pools").
		Where("project.id = ?", projectId).
		Limit(1).
		Scan(ctx)



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
	
	    log.Printf("%#v", project)

	return project, err
}
