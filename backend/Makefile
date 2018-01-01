build:
	GOOS=linux go build -o app
	GOOS=linux go build github.com/ktinubu/text2speech/client
	docker build -t gcr.io/text2speech-190611/say .
	rm -f app
	rm -f client

push:
	gcloud docker -- push gcr.io/text2speech-190611/say
