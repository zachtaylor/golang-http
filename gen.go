package http

//go:generate go-jenny -f=session/cache.go -p=session -t=Cache -k=string -v=*T

//go:generate go-jenny -f=user/cache.go -p=user -t=Cache -k=string -v=*T

//go:generate go-jenny -f=user/sockets.go -i=taylz.io/http/websocket -p=user -t=Sockets -k=string -v=*websocket.T

//go:generate go-jenny -f=websocket/cache.go -p=websocket -t=Cache -k=string -v=*T
