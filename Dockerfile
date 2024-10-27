FROM golang:1.22-alpine AS development

RUN apk --no-cache add openjdk11 bash libreoffice util-linux \
  font-droid-nonlatin font-droid ttf-dejavu ttf-freefont ttf-liberation \
	msttcorefonts-installer fontconfig && update-ms-fonts && fc-cache -f && \
  rm -rf /var/cache/apk/*

WORKDIR /app
RUN mkdir -p /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 8080

CMD ["./main"]
