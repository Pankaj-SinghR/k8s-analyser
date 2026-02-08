build:
	go build -o ./bin/k8s-analyser

run: build
	./bin/k8s-analyser