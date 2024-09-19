run:
	go run cmd/main.go

build:
	make swag-init && go build -o ./tmp/main ./cmd/main.go

build-docker:
	docker compose up --build --detach

jet-init:
	${HOME}/go/bin/jet -dsn=postgresql://plnmarketstg:79ZtFZVjmYHJaHUS@10.14.204.222:2345/plnmarketstg -path=./.gen