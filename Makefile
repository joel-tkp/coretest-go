gobuild:
	@go build -o coretest cmd/real/main.go

gorun:
	make gobuild
	@./coretest