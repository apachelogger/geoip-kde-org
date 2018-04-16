clean:
	rm -rf doc

node_modules:
	npm install apidoc

doc: node_modules
	node_modules/.bin/apidoc \
		--debug \
		-e node_modules \
		-e vendor \
		-e doc \
		-o doc

install:
	systemctl --user stop geoip-kde-org.service || true
	go install
	cp -rv systemd/* ~/.config/systemd/user/
	systemctl --user daemon-reload
	systemctl --user restart geoip-kde-org.socket

test:
	go test -v ./...

run:
	/lib/systemd/systemd-activate -l 0.0.0.0:8080 geoip-kde-org

.PHONY: doc deploy
