# Langkah pertama, pilih base image yang sesuai untuk Go
FROM golang:1.20.5-alpine

# Set working directory di dalam wadahnya
WORKDIR /app

# Copy file go.mod dan go.sum terlebih dahulu agar dependensi dapat di-cache
COPY go.mod go.sum ./

# Download dependensi menggunakan go mod
RUN go mod download

# Copy seluruh file kode Go dari direktori lokal ke direktori kerja di wadah
COPY . .

# Build aplikasi Go dan simpan binary dengan nama "app" di root direktori ("/app")
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Tandai port 8080 yang akan digunakan oleh aplikasi
EXPOSE 8080

# Run the Go Fiber application
CMD ["./app"]
