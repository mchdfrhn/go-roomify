// Package model meredefinisikan struktur data entitas bisnis aplikasi.
package model

// Division merepresentasikan entitas divisi kerja dalam sistem.
type Division struct {
	Id   string `json:"id"`   // ID unik divisi (UUID)
	Name string `json:"name"` // Nama divisi (contoh: IT, HR, Finance)
}
