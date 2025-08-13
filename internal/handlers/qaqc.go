package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func InsertQAQCHandler(req models.CreateQAQCRequest) error {
	ctx := context.Background()

	sampledOn, err := time.Parse("2006-01-02", req.SampledOn)
	if err != nil {
		return err
	}
	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(ctx, `
		INSERT INTO qaqc_entries (
			process_type, process_ref,
			containers_sampled, sampled_quantity, sampled_by, sampled_on,
			ar_number, release_date, potency, moisture_content, yield_percent,
			status, analyst_remark, analysed_by, approved_by,
			created_at, updated_at
		) VALUES (
			$1, $2,
			$3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15,
			NOW(), NOW()
		)
	`,
		req.ProcessType, req.ProcessRef,
		req.ContainersSampled, req.SampledQuantity, req.SampledBy, sampledOn,
		req.ARNumber, releaseDate, req.Potency, req.MoistureContent, req.YieldPercent,
		req.Status, req.AnalystRemark, req.AnalysedBy, req.ApprovedBy,
	)
	_, err = db.DB.Exec(ctx, `
  	UPDATE qaqc_statuses
  	SET status = 'created', updated_at = NOW()
  	WHERE process_type = $1 AND process_ref = $2
	`, req.ProcessType, req.ProcessRef)

	return err
}

func GetQAQCByRef(processType, processRef string) (*models.QAQCEntry, error) {
	ctx := context.Background()

	row := db.DB.QueryRow(ctx, `
		SELECT id, process_type, process_ref,
			containers_sampled, sampled_quantity, sampled_by, sampled_on,
			ar_number, release_date, potency, moisture_content, yield_percent,
			status, analyst_remark, analysed_by, approved_by,
			created_at, updated_at
		FROM qaqc_entries
		WHERE process_type = $1 AND process_ref = $2
	`, processType, processRef)

	var entry models.QAQCEntry
	err := row.Scan(
		&entry.ID, &entry.ProcessType, &entry.ProcessRef,
		&entry.ContainersSampled, &entry.SampledQuantity, &entry.SampledBy, &entry.SampledOn,
		&entry.ARNumber, &entry.ReleaseDate, &entry.Potency, &entry.MoistureContent, &entry.YieldPercent,
		&entry.Status, &entry.AnalystRemark, &entry.AnalysedBy, &entry.ApprovedBy,
		&entry.CreatedAt, &entry.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
