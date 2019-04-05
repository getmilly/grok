TAG = $(shell autotag)

run-rebase:
	git rebase -p master 2>/dev/null | grep "Your branch is up-to-date with 'origin/master'." || echo "\nPlease rebase your branch with master!"

run-tests:
	go test -failfast -coverprofile=coverage.out ./...

build-package:
	dep ensure

install-dep:
	go get -u github.com/golang/dep/cmd/dep

tag-version:
	git tag $(TAG) && git push origin $(TAG)