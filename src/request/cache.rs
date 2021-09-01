extern crate home;

use crate::lib::TestCase;
use regex::Regex;
use std::path::PathBuf;
use std::{fs, process};

pub fn cache_cases(
    problem_num: &str,
    test_cases: &Vec<TestCase>,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut path = home::home_dir().unwrap();
    path.push("Library");
    path.push("Caches");
    path.push("TestDrive");
    path.push(problem_num);
    fs::create_dir_all(&path)?;

    let input = "-input";
    let output = "-output";
    for (i, case) in test_cases.iter().enumerate() {
        path.push(format!("{}{}", i + 1, input));
        path.set_extension("txt");
        fs::write(&path, &case.input());
        path.pop();
        path.push(format!("{}{}", i + 1, output));
        path.set_extension("txt");
        fs::write(&path, &case.output());
        path.pop();
    }
    Ok(())
}

pub fn check_cache(problem_num: &str) -> Result<Vec<TestCase>, Box<dyn std::error::Error>> {
    let mut path = home::home_dir().unwrap();
    let mut test_cases: Vec<TestCase> = Vec::new();
    path.push("Library");
    path.push("Caches");
    path.push("TestDrive");
    path.push(problem_num);

    fs::create_dir_all(&path)?;

    let re = Regex::new(r"\d+-(input|output)\.txt").unwrap();
    let mut test_files = fs::read_dir(&path)?
        .into_iter()
        .filter(|p| p.is_ok())
        .map(|p| p.unwrap().path())
        .filter(|p| re.is_match(p.to_str().unwrap()))
        .collect::<Vec<_>>();
    test_files.sort_by(|a, b| a.to_str().unwrap().cmp(b.to_str().unwrap()));

    for i in 0..test_files.len() / 2 {
        let input = fs::read_to_string(&test_files[2 * i]).unwrap();
        let output = fs::read_to_string(&test_files[2 * i + 1]).unwrap();
        test_cases.push(TestCase::new(input, output));
    }

    Ok(test_cases)
}
