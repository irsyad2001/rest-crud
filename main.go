package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Definisikan struktur untuk menyimpan data hasil query
type BarangDetail struct {
	IdBarang     uint 
	NamaBarang   string `json:"nama_barang"`
	Stok   		 string `json:"stok"`
	ExpDate   	 string `json:"expired"`
	FotoBarang   string `json:"foto_barang"`
	Harga        string `json:"harga"`
	NamaKategori string `json:"nama_kategori,omitempty"`
	NamaJenis    string `json:"nama_jenis,omitempty"`
	NamaMaterial string `json:"nama_material,omitempty"`
	NoBatch      string `json:"no_batch,omitempty"`
}

type CreateBulkRequest struct {
	Data []BarangDetail `json:"data"`
}


func createBulkData(data []BarangDetail) error {
	tx := db.Begin()

	for _, item := range data {
		print(item.ExpDate)
		// Sesuaikan dengan struktur dan kolom pada tabel ref_barang
		query := "INSERT INTO `ref_barang` (`id_ref_barang`, `id_barang`, `stok`, `expired`, `no_batch`, `created_date`) VALUES (NULL, ?, ?, ?, ?, current_timestamp());"
		if result := tx.Exec(query, item.IdBarang, item.Stok, item.ExpDate, item.NoBatch); result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func main() {
	// Konfigurasi database dengan GORM
	dsn := "root:@tcp(127.0.0.1:3306)/db_skripsi_2?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	// Auto Migrate, atau buat tabel sesuai model BarangDetail
	db.AutoMigrate(&BarangDetail{})

	// Definisikan handler untuk endpoint /barang
	http.HandleFunc("/barang", func(w http.ResponseWriter, r *http.Request) {
		// Query menggunakan GORM
		var result []BarangDetail
		db.Table("ref_barang").
			Select("barang.nama_barang, barang.foto_barang, barang.harga, kategori.nama_kategori, jenis.nama_jenis, ref_barang.no_batch").
			Joins("INNER JOIN barang ON ref_barang.id_barang = barang.id_barang").
			Joins("INNER JOIN kategori ON barang.id_kategori = kategori.id_kategori").
			Joins("INNER JOIN jenis ON barang.id_jenis = jenis.id_jenis").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_kategori", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("barang").
			Select("nama_barang, foto_barang, kategori.nama_kategori, harga").
			Joins("INNER JOIN kategori ON barang.id_kategori = kategori.id_kategori").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_jenis", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("barang").
			Select("nama_barang, foto_barang, jenis.nama_jenis, harga").
			Joins("INNER JOIN jenis ON barang.id_jenis = jenis.id_jenis").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_material", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("barang").
			Select("nama_barang, foto_barang, material.nama_material, harga").
			Joins("INNER JOIN material ON barang.id_material = material.id_material").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_no_batch", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("barang").
			Select("nama_barang, foto_barang, ref_barang.no_batch, harga").
			Joins("INNER JOIN ref_barang ON barang.id_barang = ref_barang.id_barang").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_expired", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("ref_barang").
			Select("barang.nama_barang, barang.foto_barang, barang.harga, ref_barang.no_batch, ref_barang.expired").
			Joins("INNER JOIN barang ON ref_barang.id_barang = barang.id_barang").
			Where("ref_barang.expired <= CURRENT_DATE").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/barang_not_expired", func(w http.ResponseWriter, r *http.Request) {
		var result []BarangDetail
		db.Table("ref_barang").
			Select("barang.nama_barang,barang.foto_barang, barang.harga, ref_barang.no_batch, ref_barang.expired").
			Joins("INNER JOIN barang ON ref_barang.id_barang = barang.id_barang").
			Where("ref_barang.expired >= CURRENT_DATE").
			Scan(&result)

		// Konversi hasil ke format JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Atur header dan kirim respons JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	http.HandleFunc("/update_stok", func(w http.ResponseWriter, r *http.Request) {
		// Misalnya, Anda menerima nomor batch dan stok baru melalui parameter URL
		noBatch := r.URL.Query().Get("no_batch")
		newStok := r.URL.Query().Get("new_stok")

		// Lakukan update stok berdasarkan nomor batch
		result := db.Table("ref_barang").
			Where("no_batch = ?", noBatch).
			Update("stok", newStok)

		// Periksa kesalahan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim respons sukses
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Stok berhasil diupdate"))
	})

	http.HandleFunc("/update_harga", func(w http.ResponseWriter, r *http.Request) {
		// Misalnya, Anda menerima nomor batch dan harga baru melalui parameter URL
		noBatch := r.URL.Query().Get("no_batch")
		newHarga := r.URL.Query().Get("new_harga")

		err := db.Exec("UPDATE barang b INNER JOIN ref_barang rb ON rb.id_barang = b.id_barang SET b.harga =? WHERE rb.no_batch =?", newHarga,noBatch)
	
		// Periksa kesalahan
		if err.Error != nil {
			http.Error(w, err.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim respons sukses
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Harga berhasil diupdate"))
	})

	http.HandleFunc("/create_bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Metode HTTP tidak didukung", http.StatusMethodNotAllowed)
			return
		}
	
		var requestData CreateBulkRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
	
		// Panggil fungsi untuk menyimpan data Create Bulk ke dalam database
		if err := createBulkData(requestData.Data); err != nil {
			http.Error(w, "Error creating bulk data", http.StatusInternalServerError)
			return
		}
	
		// Kirim respons sukses
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create bulk data successful"))
	})

	
	

	// Mulai server HTTP
	port := 8080
	fmt.Printf("Server started on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
