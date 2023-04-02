#!/home/exu/Workarea/plugo/venv/bin/python

import time
import socket
import msgpack



# connect to tcp server and send msgpack data 
def send_tcp(s, in_data):
    data = msgpack.packb(in_data, use_bin_type=True)
    s.send(data)

# receive data from socket server
def recv_tcp(s):
    data = s.recv(1024)
    print("test.py recv: ", data)

if __name__ == '__main__':
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('localhost', 8080))
    while True:
        in_data = {
            'Method': 'Register',
            'Params': 
                {'Name': 'test.py', 'Address': s.getsockname()[1]},

        }
        send_tcp(s, in_data)
        time.sleep(1)

    s.close()
