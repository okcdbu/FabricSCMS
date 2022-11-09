from csv import writer as wt
from subprocess import Popen, PIPE
from sys import stderr, stdout
from timeit import default_timer
from requests import request
from os import path
import sys

script_dir = path.dirname(path.abspath(__file__))
cmd = "deployCC -ccn basic -ccp ../GurkhaContracts/asset-transfer-basic/chaincode-go/ -ccl go"


def shut_Fabric():
    url = "http://localhost:8080/fabric/network/down"
    response = request("POST", url)
    if response.status_code != 200:
        sys.exit(response.text)


def start_Fabric():
    url = "http://localhost:8080/fabric/network/up"
    response = request("POST", url)
    if response.status_code != 200:
        sys.exit(response.text)


def automated_cURL():
    shut_Fabric()
    start_Fabric()
    benchmark_start = default_timer()
    session = Popen(
        [script_dir + "/automation.sh"], shell=True, stdout=PIPE, stderr=PIPE)
    session.wait()
    benchmark_end = default_timer()
    return (benchmark_end - benchmark_start)


def deployCC():
    shut_Fabric()
    start_Fabric()
    benchmark_start = default_timer()
    session = Popen(
        [script_dir + "/network.sh deployCC -ccn basic -ccp ../GurkhaContracts/asset-transfer-basic/chaincode-go/ -ccl go" ], shell=True, stdout=PIPE, stderr=PIPE)  # TODO: This needs changing!!!
    session.wait()
    benchmark_end = default_timer()
    return (benchmark_end - benchmark_start)


def curl_test(repeats):
    with open('automation_cURL_r' + str(repeats) + '.csv', 'w', newline='')as file:
        writer = wt(file)
        writer.writerow(["Run", "automated cURL"])
        for i in range(repeats):
            benchmark = automated_cURL()
            writer.writerow([i+1, benchmark])
            print(i+1, benchmark)


def deployCC_test(repeats):
    with open('deployCC_r' + str(repeats) + '.csv', 'w', newline='')as file:
        writer = wt(file)
        writer.writerow(["Run", "deployCC"])
        for i in range(repeats):
            benchmark = deployCC()
            writer.writerow([i+1, benchmark])
            print(i+1, benchmark)


if __name__ == "__main__":
    if sys.argv[1] == "curl":
        curl_test(int(sys.argv[2]))
    elif sys.argv[1] == "deployCC":
        deployCC_test(int(sys.argv[2]))
    else:
        print(sys.argv)
