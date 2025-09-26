sudo docker rm ascii-art-web -f
sudo docker rmi ascii-art-web -f

sudo docker image build -f Dockerfile -t ascii-art-web .

sudo docker container run -p 8080:8080 --detach --name ascii-art-web ascii-art-web

sudo docker run ascii-art-web