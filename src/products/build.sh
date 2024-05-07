#!/bin/bash
par=$1


echo " " 
echo " >> Building..."

if [ "$par" == "clean" ]; then
    echo "Cleaning bin directory"
    rm -rf bin
    exit 0
fi

if [ "$par" == "all" ]; then
    archs=( "amd64" )
    oses=( "linux" "windows" )

    for os in "${oses[@]}"
    do
        for arch in "${archs[@]}"
        do
            for d in cmd/* ; do
                echo "[$os $arch] Building ${d##*/} -> ./bin/$os-$arch/${d##*/}"
                GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/$os-$arch/${d##*/} $d/main.go
            done
        done
    done
    exit 0    
fi

os=`go env GOOS`
arch=`go env GOARCH`

for d in cmd/* ; do
    echo "[$os $arch] Building ${d##*/} -> ./bin/$os-$arch/${d##*/}"
    GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/${d##*/} $d/main.go
done
