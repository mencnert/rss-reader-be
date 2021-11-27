deploy:
	git push heroku master

clean:
	rm -f rss-reader

update-example-env:
	cp app.env app.example.env

db-reset:
	go run main.go migrate down
	go run main.go migrate up