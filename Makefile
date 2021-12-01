# Binary name
BINARY=ptg
GOBUILD=go build -ldflags "-s -w" -o ${BINARY}
GOCLEAN=go clean
RMTARGZ=rm -rf *.gz
VERSION=1.0.1

# Build
build:
	$(GOCLEAN)
	$(GOBUILD)


clean:
	$(GOCLEAN)
	$(RMTARGZ)


git-tag:
	git tag $(VERSION); \
	git push origin $(VERSION); \

release:
	# Clean
	$(GOCLEAN)
	$(RMTARGZ)
	# Build for mac
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD)
	tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
	# Build for arm
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD)
	tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
	# Build for linux
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)
	tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
	# Build for win
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD).exe
	tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
	$(GOCLEAN)


gen-go-proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    reflect/gen/user.proto

gen-log-proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        gui/io/log.proto

pkg-win:
	fyne package -os windows -src gui/ -icon pic/gopher.png -name ${BINARY} -appVersion $(VERSION)

pkg-macos:
	fyne package -os darwin -src gui/ -icon pic/gopher.png -name ${BINARY} -appVersion $(VERSION)
	tar czvf ${BINARY}-mac-gui.tar ./${BINARY}.app