CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Membuat tipe data custom (ENUM) untuk status mahasiswa dan pembayaran
CREATE TYPE status_mahasiswa_enum AS ENUM ('Aktif', 'Cuti', 'Non-Aktif');
CREATE TYPE status_pembayaran_enum AS ENUM ('Lunas', 'Belum Lunas');

-- Fungsi untuk trigger yang akan memperbarui kolom 'updated_at' secara otomatis
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Membuat tabel utama untuk data mahasiswa
CREATE TABLE mahasiswas (
    id_mahasiswa UUID NOT NULL DEFAULT gen_random_uuid(),
    nim VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    ipk NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    ips_lalu NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    tahun_akademik VARCHAR(10) NOT NULL,
    semester_berjalan INT NOT NULL DEFAULT 1,
    status_mahasiswa status_mahasiswa_enum NOT NULL DEFAULT 'Aktif',
    status_pembayaran status_pembayaran_enum NOT NULL DEFAULT 'Belum Lunas',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Koma yang hilang sudah ditambahkan

    PRIMARY KEY (id_mahasiswa)
);

-- Membuat trigger yang akan menjalankan fungsi trigger_set_timestamp
-- setiap kali ada pembaruan (UPDATE) pada baris di tabel mahasiswas
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON mahasiswas
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();