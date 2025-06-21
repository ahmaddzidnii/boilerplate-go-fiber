# =========================================================================
# TAHAP 1: Build Stage (Membangun binary aplikasi)
# =========================================================================
# Menggunakan Go versi 1.24 (sesuai permintaan)
FROM golang:1.24-alpine AS builder

# Menetapkan direktori kerja di dalam container
WORKDIR /app

# Meng-copy file dependensi terlebih dahulu untuk optimasi cache Docker.
COPY go.mod go.sum ./
RUN go mod download

# Meng-copy seluruh source code aplikasi ke dalam container
COPY . .

# Membangun (compile) aplikasi Go.
# Kita menargetkan './cmd/api' sebagai sumber package 'main'.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o /app/main ./cmd/api

# =========================================================================
# TAHAP 2: Final Stage (Menjalankan aplikasi)
# =========================================================================
# Menggunakan image distroless yang sangat kecil dan aman.
FROM gcr.io/distroless/static-debian11
# Menetapkan direktori kerja
WORKDIR /

# Meng-copy HANYA binary yang sudah di-build dari 'builder' stage
COPY --from=builder /app/main /main

# PERUBAHAN: Meng-copy file .env langsung ke dalam image.
# PERINGATAN: Tidak direkomendasikan untuk produksi!
COPY .env .

# Memberi tahu Docker bahwa container akan mendengarkan di port 8080.
# Port ini harus sesuai dengan nilai APP_PORT di file .env Anda.
EXPOSE 1001

# Perintah default untuk menjalankan aplikasi ketika container dimulai.
CMD ["/main"]