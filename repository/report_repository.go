// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql dan query builder.
import (
	"database/sql"           // Interface koneksi database SQL
	"go-roomify/utils/query" // Builder query SQL dinamis
)

// ReportRepository merupakan kontrak interface untuk generasi laporan data reservasi.
type ReportRepository interface {
	GetReportByYear(startYear string, endYear string) (*sql.Rows, error) // Mengambil baris data laporan reservasi berdasar tentang tahun/tanggal
}

// reportRepository merupakan struktur konkrit pengelola repository laporan.
type reportRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// GetReportByYear mengambil data gabungan reservasi, peminjam, ruangan, dan fasilitas berdasarkan rentang tanggal/tahun.
func (self *reportRepository) GetReportByYear(startYear string, endYear string) (*sql.Rows, error) {

	// Inisialisasi builder SELECT
	query := query.QSelect{DB: self.db}

	// Menyusun relasi JOIN dan subquery agregasi string untuk fasilitas ruangan & fasilitas tambahan
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
		"tr.reservation_date",
		"tr.start_time",
		"tr.end_time",
		"tr.request_message",
		"trs.name AS status",
		"tr.response_message",
	).
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

	// Memeriksa jika terjadi error saat eksekusi query
	if err != nil {
		return nil, err
	}

	// Mengembalikan cursor *sql.Rows hasil query
	return rows, nil
}

// NewReportRepository menginisialisasi provider ReportRepository baru.
func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}
