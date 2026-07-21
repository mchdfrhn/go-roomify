FROM golang:alpine AS build
# Build Stage 
WORKDIR /app
# semua isi yang ada dalam project go-roomify akan disalin kedalam /app 
COPY . .

RUN go mod download
RUN go build -o go-roomify

# Final Stage 
FROM alpine 
WORKDIR /app
COPY --from=build /app/go-roomify /app/go-roomify

ENTRYPOINT ["/app/go-roomify"]