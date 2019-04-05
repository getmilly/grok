define get_version
$(shell cat current_version)
endef

VERSION=$(call get_version,)

install-autotag:
	wget -O autotag https://github.com/pantheon-systems/autotag/releases/download/1.1.1/Linux && sudo mv autotag /usr/local/bin/

set_version:
	autotag > current_version

run-rebase:
	git rebase -p master 2>/dev/null | grep "Your branch is up-to-date with 'origin/master'." || echo "\nPlease rebase your branch with master!"

run-tests:
	go test -failfast -coverprofile=coverage.out ./...

build-package:
	dep ensure

install-dep:
	go get -u github.com/golang/dep/cmd/dep

tag-version: set_version
	git tag $(VERSION) && git push origin $(VERSION)
