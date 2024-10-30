from os import environ
import json
import logging
import requests
import binascii

logging.basicConfig(level="INFO")
logger = logging.getLogger(__name__)

rollup_server = environ["ROLLUP_HTTP_SERVER_URL"]
logger.info(f"HTTP rollup_server url is {rollup_server}")

def hex2str(hex_string):
    return bytes.fromhex(hex_string[2:]).decode('utf-8')

def str2hex(payload):
    return "0x" + binascii.hexlify(payload.encode('utf-8')).decode('utf-8')

class Student:
    def __init__(self, name, age, wallet_address):
        self.name = name
        self.age = age
        self.wallet_address = wallet_address
        self.attendance_count = 0

    def take_attendance(self):
        self.attendance_count += 1
        logger.info(f"{self.name} has taken attendance. Attendance count: {self.attendance_count}")

all_students = []

def handle_advance(data):
    logger.info(f"Received advance request data {data}")
    logger.info(f"Received advance request data {data}")

    metadata = data["metadata"]
    sender = metadata["msg_sender"]
    payload = data["payload"]

    json_payload = json.loads(hex2str(payload))
    method = json_payload.get("method")

    if method == "create":
        create_student(json_payload)
    elif method == "sign_attendance":
        sign_attendance(sender)
    elif method == "delete":
        delete_student(json_payload.get("wallet_address"))
    else:
        logger.info("Unsupported method in payload")

    return "accept"

def create_student(json_payload):
    new_student = Student(json_payload["name"], json_payload["age"], json_payload["wallet_address"])
    all_students.append(new_student)
    logger.info(f"New student added: {new_student.name}, {new_student.age}, {new_student.wallet_address}")
    emit_notice(f"New student created: {json.dumps(new_student.__dict__)}")

def sign_attendance(wallet_address):
    found_student = next((student for student in all_students if student.wallet_address.lower() == wallet_address.lower()), None)
    if found_student:
        found_student.take_attendance()
        logger.info(f"{found_student.name} has signed attendance.")
        emit_notice(f"Attendance signed by: {wallet_address}")
    else:
        logger.info(f"Student not found: {wallet_address}")
        emit_report(f"Student not found: {wallet_address}")

def delete_student(wallet_address):
    student_index = next((index for index, student in enumerate(all_students) if student.wallet_address.lower() == wallet_address.lower()), -1)
    if student_index >= 0:
        deleted_student = all_students.pop(student_index)
        logger.info(f"Student deleted: {deleted_student.wallet_address}")
        emit_notice(f"Student deleted successfully: {wallet_address}")
    else:
        logger.info(f"Student not found: {wallet_address}")
        emit_report(f"Student not found: {wallet_address}")

def handle_inspect(data):
    logger.info(f"Received inspect request data {data}")
    logger.info(f"Received inspect request data {data}")

    payload = data["payload"]
    route = hex2str(payload)
    routes = route.split("/")

    if routes[0] == "all_students":
        students_list = [{"name": student.name, "age": student.age, "wallet_address": student.wallet_address} for student in all_students]
        emit_report(json.dumps({"students_list": students_list}))
        return json.dumps({"students_list": students_list})
    elif routes[0] == "student":
        student_address = routes[1]
        student = next((student for student in all_students if student.wallet_address.lower() == student_address.lower()), None)
        emit_report(json.dumps({"student": student.__dict__ if student else None}))
        return json.dumps({"student": student.__dict__ if student else None})
    else:
        logger.info("Unsupported route in payload")
        emit_report("Function not implemented")

    return "accept"

def emit_notice(payload):
    notice_req = requests.post(f"{rollup_server}/notice", json={
        "payload": str2hex(payload)
    })
    logger.info(f"Notice sent successfully to rollup server, with response: {notice_req.status_code}")

def emit_report(payload):
    report_req = requests.post(f"{rollup_server}/report", json={
        "payload": str2hex(payload)
    })
    logger.info(f"Report sent successfully to rollup server, with response: {report_req.status_code}")


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
