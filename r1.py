import socket
import sys
import time
import numpy as np
from pynq import allocate
from pynq import Overlay
from pynq import MMIO
import pandas as pd


class BaseClient:
    def __init__(self, timeout: int = 10, buffer: int = 4096):
        self.__socket = None
        self.__address = None
        self.__timeout = timeout
        self.__buffer = buffer

    def connect(self, address, family: int, typ: int, proto: int):
        self.__address = address
        self.__socket = socket.socket(family, typ, proto)
        self.__socket.settimeout(self.__timeout)
        self.__socket.connect(self.__address)

    def send(self, message: str = "") -> None:
        contentLength = len(message)
        self.__socket.send((str(contentLength) + "\n").encode('utf-8'))
        self.__socket.send(message.encode('utf-8'))
        try:
            pass
        except:
            self.__socket.shutdown(socket.SHUT_RDWR)
            self.__socket.close()

    def read(self):
        recv_message = ""
        while True:
            buf = self.__socket.recv(1).decode('utf-8')
            if buf == "\n":
                break
            recv_message += buf
        contentLength = int(recv_message)

        recv_message = ""
        recv_size = 0
        while True:
            buf = self.__socket.recv(self.__buffer).decode('utf-8')
            recv_message += buf
            recv_size += len(buf)
            if recv_size >= contentLength:
                break
        try:
            pass
        except:
            self.__socket.shutdown(socket.SHUT_RDWR)
            self.__socket.close
        return recv_message

    def close(self):
        self.__socket.shutdown(socket.SHUT_RDWR)
        self.__socket.close


class UnixClient(BaseClient):
    def __init__(self, path: str = "/tmp/server.sock"):
        self.server = path
        super().__init__(timeout=600, buffer=4096)
        super().connect(self.server, socket.AF_UNIX, socket.SOCK_STREAM, 0)


class FPGA():
    def __init__(self, sock_path: str = sys.argv[1], BIT_FILE: str = sys.argv[2]):
        self.BIT_FILE = BIT_FILE
        self.k = 100
        self.ovl = Overlay(BIT_FILE)
        self.dma = self.ovl.axi_dma_0
        self.done_flag = self.ovl.axi_gpio_2
        self.result_done = self.ovl.axi_gpio_3
        self.bit_pos_user = 32
        self.bit_len = 16
        self.index = 0
        self.datalist = []
        self.cli = UnixClient(sock_path)

    def exec(self):
        # self.start = time.time()
        self.cli.send("ready")

        msg = self.cli.read()
        print("data received")
        datalines = msg.splitlines(True)
        for i in range(len(datalines)):
            self.datalist.append(int(datalines[i], 2))
        print("Data size: ", len(self.datalist))

        tx_size = len(self.datalist)
        tx_buffer = allocate(shape=(tx_size,), dtype=np.uint64)
        for i in range(tx_size):
            tx_buffer[i] = self.datalist[i]
        # self.data_gen_time = time.time() - self.start

        print("data transportation for fpga")
        self.dma.sendchannel.transfer(tx_buffer)
        self.dma.sendchannel.wait()
        # self.tx_time = time.time() - self.start

        while True:
            done = self.done_flag.read()
            if done == 1:
                # self.elapsed_time = time.time() - self.start
                break
        print("data transportation has been completed")
        print("starting data processing")

        mode = self.ovl.axi_gpio_0
        mode.write(0, 1)
        while True:
            done = self.result_done.read()
            if done == 1:
                break

        print("data processing has been completed")
        print("data transportation (receive) start")
        rx_size = self.ovl.axi_gpio_1
        size = rx_size.read()
        rx_buffer = allocate(shape=(size,), dtype=np.uint64)

        self.dma.recvchannel.transfer(rx_buffer)
        self.dma.recvchannel.wait()
        # self.total_time = time.time() - start

        print("data transportation has been completed")

        # Make any map available for later retrieval from RM
        nodes_tmp = pd.read_csv("nodes_tmp_demo.csv")
        nodes_tmp.head()

        full_path = []
        bit_pos_count = 16
        bit_pos_node1 = 32
        bit_pos_node2 = 48
        bit_len = 16
        total_seg_count = 0
        more_k_seg = 0
        for i in range(len(rx_buffer) - 1):
            value = int(rx_buffer[i])
            result_count = (value >> bit_pos_count) & (0x000000000000FFFF)
            result_node1 = (value >> bit_pos_node1) & (0x000000000000FFFF)
            result_node2 = (value >> bit_pos_node2) & (0x000000000000FFFF)
            osm_node1 = nodes_tmp['osmid'][result_node1-1]
            osm_node2 = nodes_tmp['osmid'][result_node2-1]
            full_path.append([osm_node1, osm_node2, result_count])
            total_seg_count += result_count
            if result_count >= self.k:
                print("NODE_ID: ", osm_node1, " NODE_ID: ",
                      osm_node2, " COUNT: ", result_count)
                more_k_seg += result_count

        # output
        msg = ""
        for i in range(len(full_path)):
            msg += str(full_path[i][0]) + "," + str(full_path[i][1]) + "," + str(full_path[i][2]) + "\n"

        self.cli.send(msg)
        print(len(msg))


if __name__ == "__main__":
    #sock_path = sys.argv[2]
    print(sys.argv[0])
    sock_path = "/tmp/server.sock"
    fpga = FPGA(sock_path, sys.argv[2])
    while True:
        fpga.exec()
