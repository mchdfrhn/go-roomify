// Package model meredefinisikan entitas data aplikasi.
package model

// RoomType merepresentasikan kategori atau tipe dari suatu ruangan.
type RoomType struct {
	Id   string `json:"id" binding:"required"`   // ID unik tipe ruangan (UUID)
	Name string `json:"name" binding:"required"` // Nama tipe ruangan (contoh: Meeting Room, Auditorium, Classroom)
}