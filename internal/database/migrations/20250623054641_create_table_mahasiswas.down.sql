-- Migrasi 'down' ini akan menghapus semua objek yang dibuat oleh file 'up'
-- dalam urutan terbalik untuk memastikan tidak ada error dependensi.

-- STEP 1: Hapus trigger dari tabel 'mahasiswas'.
-- Trigger bergantung pada tabel dan fungsi, jadi ini harus dihapus terlebih dahulu.
DROP TRIGGER IF EXISTS set_timestamp ON mahasiswas;

-- STEP 2: Hapus tabel 'mahasiswas'.
-- Tabel harus dihapus sebelum tipe data (ENUM) dan fungsi yang mungkin bergantung padanya.
DROP TABLE IF EXISTS mahasiswas;

-- STEP 3: Hapus fungsi yang digunakan oleh trigger.
-- Fungsi ini sudah tidak lagi digunakan setelah trigger dan tabelnya dihapus.
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- STEP 4: Hapus tipe ENUM kustom.
-- Tipe data ini sudah tidak lagi digunakan setelah tabelnya dihapus.
DROP TYPE IF EXISTS status_pembayaran_enum;
DROP TYPE IF EXISTS status_mahasiswa_enum;