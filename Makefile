BINARY = chatbot
GOPATH = $(HOME)/go


all: run

run:
	chmod +x install.sh
	./install.sh
	go get github.com/wcamiller/chatbot
	export GOPATH=$(GOPATH)
	cd $(GOPATH)/src/github.com/wcamiller/chatbot && go run goServer.go
	echo $(GOPATH)
clean:
	rm -rf $(GOPATH)/src/github.com/codegangsta
	rm -rf $(GOPATH)/src/github.com/martini-contrib
	rm -rf $(GOPATH)/src/github.com/wcamiller
	rm -rf $(GOPATH)/src/github.com/go-martini
	rm -rf $(GOPATH)/src/github.com/oxtoacart
