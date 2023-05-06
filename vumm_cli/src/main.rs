use vumm_api::ClientError;

#[tokio::main]
async fn main() {
    println!("Hello, world!");

    let client = vumm_api::Client::new();

    let mod_name = String::from("mapeditor");
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
