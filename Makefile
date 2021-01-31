start:
	go run main.go

git:
	git add .
	git commit -m "$(msg)"
	git push origin master

.PHONY: start