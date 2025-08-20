binaryName=password_gen
versionPath=github.com/IIIcecream/version
version=v1.0.0
outputPath=./output
buildTime=`date "+%Y-%m-%d %H:%M:%S"`
goBuild=CGO_ENABLED=0 go build -ldflags "-X '${versionPath}.Version=$(version)' \
	-X '${versionPath}.BuildTime=$(buildTime)' \
	-X '${versionPath}.GoVersion=`go version`' \
	-X '${versionPath}.GitCommit=`git rev-parse --short HEAD`'" 

build: linux-amd64 darwin-amd64 darwin-arm64

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(goBuild) -o $(outputPath)/$(binaryName)-$@
	chmod +x $(outputPath)/$(binaryName)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(goBuild) -o $(outputPath)/$(binaryName)-$@
	chmod +x $(outputPath)/$(binaryName)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(goBuild) -o $(outputPath)/$(binaryName)-$@
	chmod +x $(outputPath)/$(binaryName)-$@

release: build
	cp README.md $(outputPath)
	tar zcvf $(binaryName)-$(version).tar.gz -C $(outputPath) .

clean:
	rm -rf output

.PHONY: all linux-amd64 darwin-amd64 darwin-arm64 release clean