#!/home/exu/Workarea/plugo/venv/bin/python

import time
import socket
import msgpack
import sys



# connect to tcp server and send msgpack data 
def send_tcp(s, in_data):
    data = msgpack.packb(in_data, use_bin_type=True)
    s.send(data)

# receive data from socket server
def recv_tcp(s):
    data = s.recv(1024)
    print("echo.py recv: ", data)


if __name__ == '__main__':
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('localhost', 8080))
    while True:
        in_data = {
            "Type": 0,
            'Method': 'Register',
            'Params': 
                {'Name': 'echo.py', 'Address': s.getsockname()[1]},
            
        }
        send_tcp(s, in_data)
        in_data = {
            "Type": 1,
            'Method': 'Register',
            'Result': 
                {'Name': 'echo.py', 'Address': s.getsockname()[1]},
            
        }
        send_tcp(s, in_data)
        recv_tcp(s)
        time.sleep(1)
    s.close()
