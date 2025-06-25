start:
	sudo docker build --tag 'main' .; sudo docker run -p 8080:8080 -d --name main main
clear:
	sudo docker stop main; sudo docker container rm main;