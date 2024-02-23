.PHONY: deploy
deploy:
	CGO_ENABLED=0 go build -o btrenamer
	mv ./btrenamer ~/.storage/qbittorrent/tools/rename/
	cp ./config.yaml ~/.storage/qbittorrent/tools/rename/