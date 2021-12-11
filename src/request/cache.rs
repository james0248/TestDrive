extern crate directories;

use crate::lib::TestCase;
use directories::BaseDirs;
use path::{Path, PathBuf};
use regex::Regex;
use std::{fs, path, process};

pub fn cache_cases(
    problem_num: &str,
    test_cases: &Vec<TestCase>,
) -> Result<(), Box<dyn std::error::Error>> {
    let cache_dir = get_cache_dir(problem_num);
    fs::create_dir_all(&cache_dir)?;

    for (i, case) in test_cases.iter().enumerate() {
        let mut input_file = Path::new(&cache_dir).join(format!("{}-input", i + 1));
        input_file.set_extension("txt");
        fs::write(&input_file, &case.input());
        let mut output_file = Path::new(&cache_dir).join(format!("{}-output", i + 1));
        output_file.set_extension("txt");
        fs::write(&output_file, &case.output());
    }
    Ok(())
}

fn get_cache_dir(problem_num: &str) -> PathBuf {
    let mut cache_dir = PathBuf::new();
    if let Some(base_dirs) = BaseDirs::new() {
        cache_dir = [base_dirs.cache_dir(), "TestDrive", problem_num]
            .iter()
            .collect();
    } else {
        panic!("Cannot find cache directory")
    }

    cache_dir
}

pub fn check_cache(problem_num: &str) -> Result<Vec<TestCase>, Box<dyn std::error::Error>> {
    let cache_dir = get_cache_dir(problem_num);
    fs::create_dir_all(&path)?;

    let mut test_cases: Vec<TestCase> = Vec::new();
    let re = Regex::new(r"\d+-(input|output)\.txt").unwrap();
    let mut test_files = fs::read_dir(&cache_dir)?
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
