// Package request menyimpan DTO untuk body/payload permintaan dari client.
package request

// DivisionRequest merepresentasikan data input HTTP request saat membuat atau memperbarui divisi.
type DivisionRequest struct {
	Id   string `json:"id"`                   // ID divisi (opsional pada pembuatan, wajib pada pembaharuan)
	Name string `json:"name" binding:"required"` // Nama divisi yang wajib diisi
}
