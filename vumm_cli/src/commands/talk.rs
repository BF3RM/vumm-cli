use clap::Args;

#[derive(Args)]
pub struct SayHello {
    name: Option<String>,
}

impl SayHello {
    pub async fn run(&self) {
        if let Some(input) = &self.name {
            println!("Hello {}!", input);
        } else {
            println!("Running command say hello without input.");
        }
    }
}

#[derive(Args)]
pub struct SayHelloNicely {
    name: Option<String>,
}

impl SayHelloNicely {
    pub async fn run(&self) {
        if let Some(input) = &self.name {
            println!("Love You {}!", input);
        } else {
            println!("Love You!");
        }
    }
}
