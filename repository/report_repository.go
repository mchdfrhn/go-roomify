package repository

import (
	"database/sql"
	"go-roomify/utils/query"
)

type ReportRepository interface {
	GetReportByYear(startYear string, endYear string) (*sql.Rows, error)
}

type reportRepository struct {
	db *sql.DB
}

func (self *reportRepository) GetReportByYear(startYear string, endYear string) (*sql.Rows, error) {

	query := query.QSelect{DB: self.db}

	rows, err := query.Table(
		"tx_reservation AS tr",
	).Column(
		"up.full_name AS user_full_name",
		"d.name AS user_division",
		"up.address AS user_address",
		"up.phone_number AS user_phone_number",
		"role.position AS user_position",
		"room.name AS room_name",
		"room_type.name AS room_type",
		"room.capacity",
		`(SELECT STRING_AGG(f.name, ', ')
			from mst_facility as f
			where room.id = f.room_id
		) as room_facilitys`,
		`(SELECT STRING_AGG(f.name, ', ')
			from tx_reservation_detail AS trd
			left join mst_facility as f
				on trd.facility_id = f.id
			where trd.reservation_id = tr.id
		) as additional_facilitys`,
		//"f.name AS facility_name",
		"tr.reservation_date",
		"tr.start_time",
		"tr.end_time",
		"tr.request_message",
		"trs.name AS status",
		"tr.response_message",
	).
		// LeftJoin(
		// 	"tx_reservation_detail AS trd",
		// 	"trd.reservation_id = tr.id",
		// ).
		LeftJoin(
			"mst_room AS room",
			"tr.room_id = room.id",
		).
		LeftJoin(
			"tx_reservation_status AS trs",
			"tr.reservation_status_id = trs.id",
		).LeftJoin(
		"mst_user_profile AS up",
		"tr.user_profile_id = up.id",
	).LeftJoin(
		"mst_division AS d",
		"up.division_id = d.id",
	).LeftJoin(
		"mst_role AS role",
		"up.role_id = role.id",
	).LeftJoin(
		"room_type",
		"room_type.id = room.room_type_id",
	).Where(
		"tr.reservation_date",
		">=",
		startYear,
	).AndWhere(
		"tr.reservation_date",
		"<=",
		endYear,
	).Run()

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}
