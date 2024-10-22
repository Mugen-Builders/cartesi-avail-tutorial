use json::{object, JsonValue};
use std::env;


pub struct Storage {
    pub toUpperTotal: u128,
    pub userStruct: Vec<UserStructure>,
}

pub struct UserStructure {
    pub sender: String,
    pub request_input: String,
    pub request_output: String,
}

pub fn hex2str(_payload: &str) -> String {
    let bytes = hex::decode(&_payload[2..]).expect("Decoding failed");
    let string_data = String::from_utf8(bytes).expect("Invalid UTF-8 sequence");
    return string_data;
}

pub fn is_numeric(input: &str) -> bool {
    input.parse::<f64>().is_ok()
}

pub async fn emit_notice(
    payload: &str,
    _server_addr: &str,
    _client: &hyper::Client<hyper::client::HttpConnector>,
) {
    // convert float to hex representation
    let hex_value = format!("0x{}", hex::encode(payload.to_string()));
    println!("hex_value is {:?}", hex_value);

    let response = object! {
        "payload" => format!("{}", hex_value)
    };

    // Send out a notice with the result
    let request = hyper::Request::builder()
        .method(hyper::Method::POST)
        .header(hyper::header::CONTENT_TYPE, "application/json")
        .uri(format!("{}/notice", &_server_addr))
        .body(hyper::Body::from(response.dump()))
        .expect("Error creating request");
    let response = _client
        .request(request)
        .await
        .expect("Error sending request");
    println!("Notice sending status {}", response.status());
}

pub async fn emit_report(
    payload: &str,
    _server_addr: &str,
    _client: &hyper::Client<hyper::client::HttpConnector>,
) {
    // convert float to hex representation
    let hex_value = format!("0x{}", hex::encode(payload.to_string()));
    println!("hex_value is {:?}", hex_value);

    let response = object! {
        "payload" => format!("{}", hex_value)
    };

    // Send out a notice with the result
    let request = hyper::Request::builder()
        .method(hyper::Method::POST)
        .header(hyper::header::CONTENT_TYPE, "application/json")
        .uri(format!("{}/report", &_server_addr))
        .body(hyper::Body::from(response.dump()))
        .expect("Error creating request");
    let response = _client
        .request(request)
        .await
        .expect("Error sending request");
    println!("Notice sending status {}", response.status());
}


pub async fn handle_advance(
    _client: &hyper::Client<hyper::client::HttpConnector>,
    _server_addr: &str,
    request: JsonValue,
    storage: &mut Storage,
) -> Result<&'static str, Box<dyn std::error::Error>> {
    println!("Received advance request data {}", &request);
    let _payload = request["data"]["payload"]
        .as_str()
        .ok_or("Missing payload")?;
    // TODO: add application logic here

    let sender = request["data"]["metadata"]["msg_sender"].as_str().ok_or("Missing sender")?;
    let sentence = hex2str(_payload);

    if is_numeric(&sentence) {
        let message = String::from("Sentence is not in string format");
        emit_report(&message, _server_addr, _client).await;
        return Ok("reject");

    } else {
        let updated_sentence = sentence.to_uppercase();
        storage.toUpperTotal += 1;
        storage.userStruct.push(UserStructure {
            sender: sender.to_string(),
            request_input: sentence.to_string(),
            request_output: updated_sentence.to_string(),
        });
        emit_notice(&updated_sentence, _server_addr, _client).await;
        println!("Processed advance request. Total upper: {}", storage.toUpperTotal);
    }

    Ok("accept")
}

pub async fn handle_inspect(
    _client: &hyper::Client<hyper::client::HttpConnector>,
    _server_addr: &str,
    request: JsonValue,
    storage: &Storage,
) -> Result<&'static str, Box<dyn std::error::Error>> {
    println!("Received inspect request data {}", &request);
    let _payload = request["data"]["payload"]
        .as_str()
        .ok_or("Missing payload")?;
    // TODO: add application logic here

    let query = hex2str(_payload);
    if query.contains("list") {
        let mut result = String::new();
        for user in &storage.userStruct {
            result.push_str(&format!("{{Sender: '{}', Input: '{}', Output: '{}'}}\n", user.sender, user.request_input, user.request_output));
        }
        emit_report(&result, _server_addr, _client).await;
        println!("Processed inspect request. all upper: {}", storage.toUpperTotal);
    } else if query.contains("total") {
        emit_report(&format!("Total upper: {}", storage.toUpperTotal), _server_addr, _client).await;
        println!("Processed inspect request. Total upper: {}", storage.toUpperTotal);
    }
    Ok("accept")
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = hyper::Client::new();
    let server_addr = env::var("ROLLUP_HTTP_SERVER_URL")?;
    let mut storage: Storage = Storage {
        toUpperTotal: 0,
        userStruct: Vec::new(),
    };

    let mut status = "accept";
    loop {
        println!("Sending finish");
        let response = object! {"status" => status.clone()};
        let request = hyper::Request::builder()
            .method(hyper::Method::POST)
            .header(hyper::header::CONTENT_TYPE, "application/json")
            .uri(format!("{}/finish", &server_addr))
            .body(hyper::Body::from(response.dump()))?;
        let response = client.request(request).await?;
        println!("Received finish status {}", response.status());

        if response.status() == hyper::StatusCode::ACCEPTED {
            println!("No pending rollup request, trying again");
        } else {
            let body = hyper::body::to_bytes(response).await?;
            let utf = std::str::from_utf8(&body)?;
            let req = json::parse(utf)?;

            let request_type = req["request_type"]
                .as_str()
                .ok_or("request_type is not a string")?;
            status = match request_type {
                "advance_state" => handle_advance(&client, &server_addr[..], req, &mut storage).await?,
                "inspect_state" => handle_inspect(&client, &server_addr[..], req, &storage).await?,
                &_ => {
                    eprintln!("Unknown request type");
                    "reject"
                }
            };
        }
    }
}
