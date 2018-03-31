all: kuehlschrank

kuehlschrank: main.go $(wildcard */*.go)
	GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags="-s -w"
	upx kuehlschrank

clean:
	rm kuehlschrank

testdeploy: kuehlschrank
	scp kuehlschrank root@10.25.11.32:/tmp/kuehlschrank

deploy: kuehlschrank
	scp kuehlschrank root@10.25.11.32:/root/kuehlschrank.bin

