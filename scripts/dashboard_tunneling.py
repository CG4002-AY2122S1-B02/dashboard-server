# Test script to try out the sshtunnel package

import paramiko
import sshtunnel
from paramiko import SSHClient

TUNNEL_ONE_SSH_ADDR = "sunfire.comp.nus.edu.sg"
TUNNEL_ONE_SSH_USERNAME = "e0325893"
TUNNEL_ONE_SSH_PASSWORD = "iLoveCapstoneB02"

TUNNEL_TWO_SSH_ADDR = "137.132.86.225"
TUNNEL_TWO_SSH_USERNAME = "xilinx"
TUNNEL_TWO_SSH_PASSWORD = "cg4002b02"

port_nums = [8880, 8881, 8882, 8883]

tunnel_one =  sshtunnel.open_tunnel(
    # Port 22 open for SSH
    (TUNNEL_ONE_SSH_ADDR,22), # Remote Server IP
    ssh_username=TUNNEL_ONE_SSH_USERNAME,
    ssh_password=TUNNEL_ONE_SSH_PASSWORD,
    remote_bind_address=(TUNNEL_TWO_SSH_ADDR,22), # Private Server IP
)

tunnel_one.start()
print("Connection to tunnel_one (sunfire:22) OK...")

for i in range(len(port_nums)):

    tunnel_two = sshtunnel.open_tunnel(
        ssh_address_or_host=('127.0.0.1',tunnel_one.local_bind_port),
        remote_bind_address=('127.0.0.1',port_nums[i]),
        ssh_username=TUNNEL_TWO_SSH_USERNAME,
        ssh_password=TUNNEL_TWO_SSH_PASSWORD,
        local_bind_address=('127.0.0.1',port_nums[i])
    )

    tunnel_two.start()
    print(f"Connection to tunnel_two (137.132.86.225:{port_nums[i]}) OK...")

while True:
    continue