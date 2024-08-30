# go-stock-api

Untuk table yang harus di crete yaitu :

CREATE TABLE stocks (
    id SERIAL PRIMARY KEY,
    nama_barang VARCHAR(100) NOT NULL,
    jumlah INT NOT NULL,
	nomor_seri INT NOT NULL,
	additional_info JSONB,
	gambar_barang VARCHAR(250),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Untuk testing postman nya https://www.postman.com/poolapack-teknologi/workspace/test-dki-public/collection/29465757-97a3d020-eea3-455b-8732-d97d10611ccf?action=share&creator=29465757