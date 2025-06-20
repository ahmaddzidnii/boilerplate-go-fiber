package config

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// =================================================================
// Impelementasi AsyncWriter untuk menulis log secara asinkron
// =================================================================

// AsyncWriter adalah writer kustom yang akan menulis log secara asinkron.
// Ia mengimplementasikan interface io.Writer.
type AsyncWriter struct {
	logChan chan []byte
	writer  io.Writer
}

// NewAsyncWriter membuat instance baru dari AsyncWriter.
// Ia juga memulai sebuah goroutine di latar belakang untuk memproses log.
func NewAsyncWriter(writer io.Writer, bufferSize int) *AsyncWriter {
	aw := &AsyncWriter{
		// Membuat buffered channel untuk menampung log sementara.
		// Ukuran buffer bisa disesuaikan. Jika penuh, pemanggilan Write akan memblokir.
		logChan: make(chan []byte, bufferSize),
		writer:  writer,
	}

	// Memulai satu goroutine di latar belakang.
	// Goroutine ini akan bertanggung jawab untuk benar-benar menulis log.
	go func() {
		// Loop ini akan terus berjalan selama channel terbuka.
		for data := range aw.logChan {
			// Melakukan operasi I/O yang lambat di sini.
			aw.writer.Write(data)
		}
	}()

	return aw
}

// Write adalah method yang dipanggil oleh logrus.
// Method ini tidak langsung menulis ke file, tapi mengirim data log ke channel.
func (aw *AsyncWriter) Write(p []byte) (n int, err error) {
	// Kita perlu membuat salinan dari slice 'p' karena logrus dapat
	// menggunakan kembali buffer ini. Jika tidak disalin, data di channel bisa
	// berubah sebelum sempat ditulis oleh goroutine writer.
	data := make([]byte, len(p))
	copy(data, p)

	aw.logChan <- data
	return len(p), nil
}

func InitLogger() *logrus.Logger {
	logger := logrus.New()

	level := GetEnv("LOG_LEVEL", "info")

	switch level {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	}

	logger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		logger.Fatalf("Gagal membuka file log: %v", err)
	}

	multi := io.MultiWriter(file, os.Stdout)
	asyncWriter := NewAsyncWriter(multi, 100)
	logger.SetOutput(asyncWriter)
	return logger
}
