use std::process::Command;
use std::{env, fs, process};

mod lib;
mod request;

fn main() {
    let args: Vec<String> = env::args().collect();
    let config = lib::Config::new(&args).unwrap_or_else(|err| {
        println!("Problem parsing arguments: {}", err);
        process::exit(1);
    });

    let test_cases = request::parse_cases(&config.problem_num).unwrap();
    lib::run_tests(&config.binary, &test_cases);
}
