deploy:
	git push heroku master

clean:
	rm -f rss-reader

update-example-env:
	cp app.env app.example.env