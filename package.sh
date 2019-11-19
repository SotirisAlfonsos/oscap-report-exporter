#!/usr/bin/env bash

VERSION=$1

go test ./... -v -cover -count=1

package="oscap-report-exporter"
package_name="oscap-exporter"

platforms=("linux/arm64" "linux/amd64" "linux/386" "linux/mips" "linux/mips64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    tar_output_name=$package_name'-'$VERSION'.'$GOOS'-'$GOARCH

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $package_name $package
    tar -czvf $tar_output_name.tar.gz $package_name
    rm $package_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

# go build scheduler.go
# tar -czvf oscap-exporter-${VERSION}.tar.gz scheduler