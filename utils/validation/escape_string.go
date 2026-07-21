// Package validation menyediakan helper untuk sanitasi dan validasi data string.
package validation

// Import package fmt dan strings.
import (
	"fmt"     // Format penulisan string
	"strings" // Fungsi manipulasi string
)

// EscapeString melepaskan karakter petik tunggal (') untuk mencegah SQL Injection jika digunakan secara manual.
func EscapeString(value any) string {
	// Melakukan pengecekan tipe data variabel value
	switch value.(type) {
	case string:
		// Mengubah petik tunggal menjadi dua petik tunggal (' -> '') dan mengapitnya dengan petik
		return "'" + strings.ReplaceAll(value.(string), "'", "''") + "'"
	default:
		// Mengembalikan format representasi string bawaan jika bukan tipe string
		return fmt.Sprintf("%v", value)
	}
}
