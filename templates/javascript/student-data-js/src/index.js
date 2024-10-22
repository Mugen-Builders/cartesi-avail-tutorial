// XXX even though ethers is not used in the code below, it's very likely
// it will be used by any DApp, so we are already including it here
const { ethers } = require("ethers");

const rollup_server = process.env.ROLLUP_HTTP_SERVER_URL;
console.log("HTTP rollup_server url is " + rollup_server);


function Circle(name, age, wallet_address) {
  this.name = name;
  this.age = age;
  this.wallet_address = wallet_address;
  this.attendance_count = 0;

  this.take_attendance = function () {
    this.attendance_count++;
    console.log(`${this.name} has taken attendance. Attendance count: ${this.attendance_count}`);
  }
}

let all_students = [];

function hex2str(hex) {
  return ethers.toUtf8String(hex)
}

function str2hex(payload) {
  return ethers.hexlify(ethers.toUtf8Bytes(payload))
}

async function handle_advance(data) {
  console.log("Received advance request data " + JSON.stringify(data));

  const metadata = data["metadata"]
  const sender = metadata["msg_sender"]
  const payload = data["payload"]

  let json_payload = JSON.parse(hex2str(payload));

  let method = json_payload["method"];

  switch (method) {
    case "create":
      create_student(json_payload);
      break;

    case "sign_attendance":
      sign_attendance(sender);
      break;

    case "delete":
      delete_student(json_payload["wallet_address"]);
      break;

    default:
      console.log("Unsupported method in payload");
      break;
  }

  return "accept";
}

async function create_student(json_payload) {
  let new_student = new Circle(json_payload["name"], json_payload["age"], json_payload["wallet_address"]);
  all_students.push(new_student);
  console.log(`New student added: ${new_student.name}, ${new_student.age}, ${new_student.wallet_address}`);
  emit_notice(`New student created: ${JSON.stringify(new_student)}`);
}


async function sign_attendance(wallet_address) {
  let found_student = all_students.find(student => student.wallet_address.toLowerCase() === wallet_address.toLowerCase());
  if (found_student) {
    found_student.take_attendance();
    console.log(`${found_student.name} has signed attendance.`);
    emit_notice(`Attendance signed by: ${wallet_address}`);
  } else {
    console.log(`Student not found: ${student_name}`);
    await emit_report(`Student not found: ${wallet_address}`);
  }
}

async function delete_student(wallet_address) {
  let student_index = all_students.findIndex(student => student.wallet_address.toLowerCase() === wallet_address.toLowerCase());
  if (student_index >= 0) {
    all_students.splice(student_index, 1);
    console.log(`Student deleted: ${wallet_address}`);
    emit_notice(`Student deleted sucessfully: ${wallet_address}`);

  } else {
    console.log(`Student not found: ${wallet_address}`);
    await emit_report(`Student not found: ${wallet_address}`);

  }
}

async function handle_inspect(data) {
  console.log("Received inspect request data " + JSON.stringify(data));

  const payload = data["payload"]
  const route = hex2str(payload)

  let routes = route.split("/");

  switch (routes[0]) {
    case "all_students":
      let students_list = all_students.map(student => ({ name: student.name, age: student.age, wallet_address: student.wallet_address }));
      emit_report(JSON.stringify({ students_list }));
      return JSON.stringify({ students_list });

    case "student":
      let student_address = routes[1];
      let student = all_students.find(student => student.wallet_address.toLowerCase() === student_address.toLowerCase());
      await emit_report(JSON.stringify({ student }));
      return JSON.stringify({ student });

    default:
      console.log("Unsupported route in payload");
      await emit_report("Function not implemented");
      break;
  }
  return "accept";
}


async function emit_notice(payload) {
  const notice_req = await fetch(rollup_server + "/notice", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ payload: str2hex(payload) }),
  });

  console.log(`notice sent succesfully to rollup server, with response: ${notice_req.status}`);

}

async function emit_report(payload) {
  const notice_req = await fetch(rollup_server + "/report", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ payload: str2hex(payload) }),
  });

  console.log(`report sent succesfully to rollup server, with response: ${notice_req.status}`);
}

var handlers = {
  advance_state: handle_advance,
  inspect_state: handle_inspect,
};

var finish = { status: "accept" };

(async () => {
  while (true) {
    const finish_req = await fetch(rollup_server + "/finish", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ status: "accept" }),
    });

    console.log("Received finish status " + finish_req.status);

    if (finish_req.status == 202) {
      console.log("No pending rollup request, trying again");
    } else {
      const rollup_req = await finish_req.json();
      var handler = handlers[rollup_req["request_type"]];
      finish["status"] = await handler(rollup_req["data"]);
    }
  }
})();
