build:
	sudo docker build -t arbitrage_watcher .

run:
	sudo docker rm -f arbitrage_watcher || false
	sudo docker run --name arbitrage_watcher --env-file .env arbitrage_watcher

all: build run