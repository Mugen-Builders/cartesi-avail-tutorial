from os import environ
import logging
import requests
import binascii
import re

logging.basicConfig(level="INFO")
logger = logging.getLogger(__name__)

rollup_server = environ["ROLLUP_HTTP_SERVER_URL"]
logger.info(f"HTTP rollup_server url is {rollup_server}")

def hex2str(hex_string):
    return bytes.fromhex(hex_string[2:]).decode('utf-8')

def str2hex(payload):
    return "0x" + binascii.hexlify(payload.encode('utf-8')).decode('utf-8')

def is_numeric(num):
    return bool(re.match(r"^-?\d+(\.\d+)?$", str(num)))

users = []
to_upper_total = 0

def handle_advance(data):
    global users
    global to_upper_total

    logger.info(f"Received advance request data {data}")

    metadata = data["metadata"]
    sender = metadata["msg_sender"]
    payload = data["payload"]

    # Convert hex payload to string and validate
    sentence = hex2str(payload)
    if is_numeric(sentence):
        report_req = requests.post(f"{rollup_server}/report", json={
            "payload": str2hex("sentence is not in string format")
        })
        if report_req.status_code < 200 or report_req.status_code > 299:
            logger.error("Failed to send report")
        return "reject"

    # Convert sentence to uppercase
    updated_sentence = sentence.upper()
    user_structure = {"sender": sender, "request input": sentence, "request output": updated_sentence}
    users.append(user_structure)
    global to_upper_total
    to_upper_total += 1

    # Send notice
    notice_req = requests.post(f"{rollup_server}/notice", json={
        "payload": str2hex(updated_sentence)
    })
    if notice_req.status_code < 200 or notice_req.status_code > 299:
        logger.error("Failed to send notice")
        return "reject"

    return "accept"


def handle_inspect(data):
    logger.info(f"Received inspect request data {data}")

    payload = data["payload"]
    route = hex2str(payload)

    if route == "list":
        response_object = {"users": users}
    elif route == "total":
        response_object = {"toUpperTotal": to_upper_total}
    else:
        response_object = "route not implemented"

    report_req = requests.post(f"{rollup_server}/report", json={
        "payload": str2hex(str(response_object))
    })
    if report_req.status_code < 200 or report_req.status_code > 299:
        logger.error("Failed to send report")

    return "accept"

handlers = {
    "advance_state": handle_advance,
    "inspect_state": handle_inspect,
}

finish = {"status": "accept"}

while True:
    logger.info("Sending finish")
    response = requests.post(rollup_server + "/finish", json=finish)
    logger.info(f"Received finish status {response.status_code}")
    if response.status_code == 202:
        logger.info("No pending rollup request, trying again")
    else:
        rollup_request = response.json()
        data = rollup_request["data"]
        handler = handlers[rollup_request["request_type"]]
        finish["status"] = handler(rollup_request["data"])
