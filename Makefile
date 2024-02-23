.PHONY: deploy
deploy:
	go build -o btrenamer
	mv ./btrenamer ~/.storage/qbittorrent/tools/rename/
	cp ./config.yaml ~/.storage/qbittorrent/tools/rename/