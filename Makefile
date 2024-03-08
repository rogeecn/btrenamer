.PHONY: deploy
deploy: build
	mv ./btrenamer ~/.storage/qbittorrent/tools/rename/
	cp ./config.yaml ~/.storage/qbittorrent/tools/rename/

.PHONY: build
build:
	CGO_ENABLED=0 go build -o btrenamer