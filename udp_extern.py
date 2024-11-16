import socket

sock = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)

ret = sock.sendto(bytes("0"*100,"utf-8"),("127.0.0.1",53))
print(ret)
