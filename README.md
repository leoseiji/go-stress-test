# go-stress-test


Criar a imagem do docker: docker build -t stress_test .

Rodar o stress test:  docker run stress_test --url=http://google.com --requests=1000 --concurrency=10