LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s "

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gomake.exe -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomake-linux -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gomake-linux-arm64 -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gomake-darwin -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o gomake-darwin-arm64 -trimpath -ldflags "${LDFLAGS}"

# sha256
sha256sum gomake* > gomake-sha256
cat gomake-sha256

# chmod 
chmod +x gomake-*

# gzip
gzip --best gomake*