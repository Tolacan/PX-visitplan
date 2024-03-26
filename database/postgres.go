package database

import (
	"PX-visitplan/models"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) InsertVisitPlan(ctx context.Context, visit *models.VisitPlan) error {
	query := `
	INSERT INTO visitplans(uuid,nombre,ruta,uuid_ruta,lunes,martes,miercoles,jueves,viernes,sabado,domingo,responsable,fecha_registro,fecha_modificacion,activo)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
`
	_, err := repo.db.Exec(
		query,
		visit.Uuid,
		visit.Nombre,
		visit.Ruta,
		visit.UuidRuta,
		visit.Lunes,
		visit.Martes,
		visit.Miercoles,
		visit.Jueves,
		visit.Viernes,
		visit.Sabado,
		visit.Domingo,
		visit.Responsable,
		visit.FechaRegistro,
		visit.FechaModificacion,
		visit.Activo,
	)

	queryClient := `
		INSERT INTO visitplans_clients (uuid_visitplan,uuid_cliente)
		VALUES ($1,$2)
	`
	for _, client := range visit.Clientes {
		_, err = repo.db.Exec(
			queryClient,
			visit.Uuid,
			client,
		)
		if err != nil {
			return err
		}

	}

	return err
}

func (repo *PostgresRepository) ListVisitPlans(ctx context.Context) ([]*models.VisitPlan, error) {

	query := `
	SELECT uuid,nombre,ruta,uuid_ruta,lunes,martes,miercoles,jueves,viernes,sabado,domingo,responsable,fecha_registro,fecha_modificacion,activo
	FROM visitplans
`
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var visits []*models.VisitPlan
	for rows.Next() {
		visit := &models.VisitPlan{}
		err := rows.Scan(
			&visit.Uuid,
			&visit.Nombre,
			&visit.Ruta,
			&visit.UuidRuta,
			&visit.Lunes,
			&visit.Martes,
			&visit.Miercoles,
			&visit.Jueves,
			&visit.Viernes,
			&visit.Sabado,
			&visit.Domingo,
			&visit.Responsable,
			&visit.FechaRegistro,
			&visit.FechaModificacion,
			&visit.Activo,
		)
		if err != nil {
			return nil, err
		}
		visits = append(visits, visit)
	}

	for _, visitsClient := range visits {
		queryClients := `SELECT uuid_cliente FROM visitplans_clients WHERE uuid_visitplan=$1`
		rows, err := repo.db.QueryContext(ctx, queryClients, visitsClient.Uuid)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		var clients []string
		for rows.Next() {
			var client string
			err := rows.Scan(&client)
			if err != nil {
				return nil, err
			}
			clients = append(clients, client)
		}
		visitsClient.Clientes = clients
	}

	return visits, nil
}

func (repo *PostgresRepository) UpdateVisitPlan(ctx context.Context, visit *models.VisitPlan) error {
	query := `
  	UPDATE visitplans 
	SET nombre=$1,ruta=$2,uuid_ruta=$3,lunes=$4,martes=$5,miercoles=$6,jueves=$7,viernes=$8,sabado=$9,domingo=$10,responsable=$11,fecha_registro=$12,fecha_modificacion=$13,activo=$14
	WHERE uuid=$15
`
	_, err := repo.db.Exec(
		query,
		visit.Nombre,
		visit.Ruta,
		visit.UuidRuta,
		visit.Lunes,
		visit.Martes,
		visit.Miercoles,
		visit.Jueves,
		visit.Viernes,
		visit.Sabado,
		visit.Domingo,
		visit.Responsable,
		visit.FechaRegistro,
		visit.FechaModificacion,
		visit.Activo,
		visit.Uuid,
	)

	if len(visit.Clientes) > 0 {
		queryDelete := `DELETE FROM visitplans_clients WHERE uuid_visitplan=$1`
		_, errdb := repo.db.Exec(queryDelete, visit.Uuid)
		if errdb != nil {
			return errdb
		}

		queryClient := `INSERT INTO visitplans_clients (uuid_visitplan,uuid_cliente) VALUES ($1,$2)`
		for _, client := range visit.Clientes {
			_, err = repo.db.Exec(queryClient, visit.Uuid, client)
			if err != nil {
				return err
			}
		}
	}
	return err
}
