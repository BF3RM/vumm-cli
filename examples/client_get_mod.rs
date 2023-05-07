use vumm_api::ClientError;

#[tokio::main]
async fn main() {
    let mut client = vumm_api::Client::new();

    client.set_bearer_token("007e9635-7a92-4796-adb9-26d468de6b74".to_string());

    let mod_name = String::from("realitymod");
    let mod_version = String::from("0.2.0");

    let mod_response = client.mods().get_version(mod_name, mod_version).await;

    match mod_response {
        Ok(mod_) => {
            println!("Mod: {:?}", mod_);
        }
        Err(e) => match e {
            ClientError::StatusCode(response) => {
                println!("Error: {}", response.status());
            }
            ClientError::Internal(e) => {
                println!("Error: {}", e);
            }
        },
    }
}
