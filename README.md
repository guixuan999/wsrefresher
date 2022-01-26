build:
go build -o wsrefresherservice.exe wisersoft.com.cn/wsrefresher

install service
wsrefresherservice.exe install

start service
sc start wsrefresher