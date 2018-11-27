
all:
	docker-compose build
	docker-compose run start_consul
	docker-compose up -d uservices

install:
	go get github.com/securego/gosec/cmd/gosec/...
	go get -u github.com/kardianos/govendor
	go get golang.org/x/tools/cmd/goimports
	./scripts/install_git_hooks.sh
	docker-compose pull

uninstall:
	docker-compose stop
	docker-compose down --rmi all
	./scripts/uninstall_git_hooks.sh

clean:
	docker-compose down --rmi local
