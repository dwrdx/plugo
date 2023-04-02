
import time
import socket
import msgpack



# connect to tcp server and send msgpack data 
def send_tcp(in_data):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('localhost', 8080))
    # data = msgpack.packb(in_data, use_bin_type=True)
    # s.recv(1024)
    # receive data from socket server
    data = s.recv(1024)



    s.send(data)
    s.close()


if __name__ == '__main__':
    while True:
        in_data = {
            'Foo': 'test',
            'Params': 
                {'Name': 'param1', 'Value': 1},
            
        }
        send_tcp(in_data)
        time.sleep(1)
