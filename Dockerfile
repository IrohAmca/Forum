# Üst imaj olarak resmi bir Golang imajı kullanıldı
FROM golang:1.22.3-alpine

# Konteyner içerisindeki geçerli çalışma dizini ayarlandı
# /app klasorunu docker icin olusturuyor
WORKDIR /app

# go mod ve sum dosyaları kopyalanır
COPY go.mod go.sum ./

# Tüm bağımlılıklar inidirilir
RUN go mod download

# Derleme bağımlılıkları yüklenir
RUN apk add --no-cache build-base


# Kaynak kodu konteynere kopyalanır
COPY . .

# Go uygulaması oluşturulur
# chat.go yu da dahil eder main adinda tek bir executable olusturuyor
ENV CGO_ENABLED=1 
RUN go build -o main .


# 8080 numaralı bağlantı noktasını dış dünyaya açar
EXPOSE 8080

# Yürütülebilir dosyayı çalıştırır
CMD ["./main"]
