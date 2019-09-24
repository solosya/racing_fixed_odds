FROM golang:latest

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 8080

# docker run -it -v $(pwd):/go/src/github.com/user/myProject/app golang /bin/bash
#./go build
#./app -api=LAD