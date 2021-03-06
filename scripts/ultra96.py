# Test Script to send data (one-way communication) from Ultra96 FPGA to Database Server
# Implementation using Python Socket API

#python3 scripts/ultra96.py 8881
#python3 scripts/ultra96.py 8882
#python3 scripts/ultra96.py 8883

import sys
import socket
import time
import packet_pb2
from random import seed
from random import random
seed(1)

user = 1

packet_stream_test = [
    packet_pb2.Packet(
        user=1,
        dance_move = "Dab",
        accuracy=2,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Scarecrow",
        accuracy=1,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Window360",
        accuracy=3,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "James Bond",
        accuracy=3,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Mermaid",
        accuracy=2,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Push Back",
        accuracy=2,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Snake",
        accuracy=3,
    ),
    packet_pb2.Packet(
        user=1,
        dance_move = "Cowboy",
        accuracy=1,
    )
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

    def send_data(self, packet):
        packet.end = "\x7F"
#         print(f"Sending data to dashboard comm client", packet)
        self.socket.sendall(packet.SerializeToString())

    def stop(self):
        self.connection.close()
        self.shutdown.set()
        self.timer.cancel()


def main():
    ip_addr = '127.0.0.1'
    port_num = 8881
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
        packet = packet_stream_test[count]
        packet.epoch_ms = int(time.time() * 1000)
        packet.user = user
        my_client.send_data(packet)
        # my_client.send_data("# 1 2 3 | dab | 1.00")
        # Receive the new dance move instructions from the Evaluation Server
        time.sleep(2) # + random() * 0.5)
        count = (count + 1) % len(packet_stream_test)
#         if (count == len(packet_stream_test)):
#             my_client.stop()


if __name__ == '__main__':
    main()
