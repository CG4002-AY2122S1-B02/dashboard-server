# Test Script to send data (one-way communication) from Ultra96 FPGA to Database Server
# Implementation using Python Socket API

#python3 scripts/ultra96position.py
# 8880

import sys
import socket
import time
import packet_pb2

position_stream_test = [
    packet_pb2.Position(
        position="123"
    ),
    packet_pb2.Position(
        position="213"
    ),
    packet_pb2.Position(
        position="312"
    ),
    packet_pb2.Position(
        position="321"
    ),
    packet_pb2.Position(
        position="231"
    ),
    packet_pb2.Position(
        position="132"
    ),
]

class Client():
    def __init__(self, ip_addr, port_num, group_id):
        super(Client, self).__init__()

        # Create a TCP/IP socket and connect to database server
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_address = (ip_addr, port_num)
        self.group_id = group_id

        print('trying to connect to %s port %s' % server_address)
        self.socket.connect(server_address)
        print("Successfully connected to the dashboard server")

    def send_data(self, position):
        position.end = "\x7F"
#         print(f"Sending data to dashboard comm client", position)
        self.socket.sendall(position.SerializeToString())

    def stop(self):
        self.connection.close()
        self.shutdown.set()
        self.timer.cancel()


def main():
    ip_addr = '127.0.0.1'
    port_num = 8880
    group_id = 2
    if len(sys.argv) == 2:
        port_num = int(sys.argv[1])
    elif len(sys.argv) != 4:
        print('Invalid number of arguments')
        print('python eval_client.py [IP address] [Port] [groupID]')
        print('using default:localhost:8080, 2')
        # sys.exit()
    else:
        ip_addr = sys.argv[1]
        port_num = int(sys.argv[2])
        group_id = sys.argv[3]

    user = port_num % 10
    my_client = Client(ip_addr, port_num, group_id)
    action = ""

    count = 0
    while action != "logout":
        # Send the Evaluation Server the received data from the 3 laptops
        position = position_stream_test[count]
        position.epoch_ms = int(time.time() * 1000)
        my_client.send_data(position)
        # my_client.send_data("# 1 2 3 | dab | 1.00")
        # Receive the new dance move instructions from the Evaluation Server
        time.sleep(3)
        count = (count + 1) % len(position_stream_test)
#         if (count == len(packet_stream_test)):
#             my_client.stop()


if __name__ == '__main__':
    main()
