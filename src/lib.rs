extern crate colored;

use colored::*;

use std::io::{Read, Write};
use std::process::{Command, Stdio};

pub fn run_tests(binary: &str, test_cases: &Vec<TestCase>) {
    println!(
        "{}",
        format!("Running tests on {} cases...", test_cases.len())
            .blue()
            .bold()
    );
    for (i, case) in test_cases.iter().enumerate() {
        let process = match Command::new(binary)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .spawn()
        {
            Err(why) => panic!("Couldn't spawn C++ executable! {}", why),
            Ok(process) => process,
        };
        match process.stdin.unwrap().write_all(&case.input.as_bytes()) {
            Err(why) => panic!("couldn't write to stdin: {}", why),
            Ok(_) => (),
        };
        let mut user_stdout = String::new();
        match process.stdout.unwrap().read_to_string(&mut user_stdout) {
            Err(why) => panic!("couldn't read stdout: {}", why),
            Ok(_) => {
                if &user_stdout == &case.output {
                    println!("{}", format!("[Case #{} Passed!]", i + 1).green().bold());
                } else {
                    let warning = format!("[Case #{} Failed!]\n", i + 1).red().bold();
                    let input = format!("{}{}", "Input data:\n".bold(), &case.input);
                    let user_output = format!("{}{}", "Your output:\n".bold(), user_stdout);
                    let case_output = format!("{}{}", "Ground truth:\n".bold(), &case.output);
                    println!("{}{}{}{}", warning, input, user_output, case_output);
                };
            }
        }
    }
}

pub struct TestCase {
    input: String,
    output: String,
}

impl TestCase {
    pub fn new(input: String, output: String) -> TestCase {
        TestCase { input, output }
    }
    pub fn input(&self) -> &str {
        &self.input
    }
    pub fn output(&self) -> &str {
        &self.output
    }
}

pub struct Config {
    pub binary: String,
    pub problem_num: String,
}

impl Config {
    pub fn new(args: &[String]) -> Result<Config, &str> {
        if args.len() < 3 {
            return Err("not enough arguments");
        }
        let binary = args[1].clone();
        let problem_num = args[2].clone();

        Ok(Config {
            binary,
            problem_num,
        })
    }
}
