use vumm_cli;

fn main() {
    if let Err(e) = vumm_cli::run_cli() {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}
