#[tokio::main]
async fn main() {
    println!("Hello, world!");

    let mut client = vumm_api::Client::new();

    let mod_name = String::from("realitymod");
    let mod_version = String::from("0.2.0");

    // let mod_response = client.mods().get(mod_name).await;
    let mod_archive = client.mods().download_version(mod_name, mod_version).await;

    mod_archive.unwrap().unpack("test").unwrap();
}
