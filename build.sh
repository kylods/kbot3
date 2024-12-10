echo "Building linux server binary..."
go build -o bin/kbot-server ./cmd/server

echo "Building linux client binary..."
go build -o bin/kbot-client ./cmd/client

echo "Building windows client binary..."
fyne-cross windows -arch=amd64 -app-id github.com/kylods/kbot3/client ./cmd/client