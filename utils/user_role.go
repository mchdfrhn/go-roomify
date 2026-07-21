// Package utils menyediakan sekumpulan fungsi pembantu (utility) untuk aplikasi go-roomify.
package utils

// Deklarasi konstanta peran pengguna (User Role) dalam sistem.
const (
	USER_ROLE_ADMIN    = "admin"    // Peran administrator sistem dengan hak akses penuh
	USER_ROLE_GA       = "ga"       // Peran General Affair (GA) pengelola ruangan dan fasilitas
	USER_ROLE_EMPLOYEE = "employee" // Peran karyawan (Employee) sebagai peminjam/pemesan ruangan
)
