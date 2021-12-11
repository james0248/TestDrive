extern crate reqwest;
extern crate scraper;

use scraper::{Html, Selector};

use crate::lib::TestCase;
use crate::request::cache::{cache_cases, check_cache};

pub mod cache;

#[tokio::main]
pub async fn parse_cases(problem_num: &str) -> Result<Vec<TestCase>, Box<dyn std::error::Error>> {
    // Check if cached data exists
    match check_cache(problem_num) {
        Ok(test_cases) => {
            if !test_cases.is_empty() {
                return Ok(test_cases);
            }
        }
        Err(err) => return Err(err),
    }

    let mut test_cases: Vec<TestCase> = Vec::new();
    let url = format!("https://www.acmicpc.net/problem/{}", problem_num);
    let resp = reqwest::get(url).await?.text().await?;

    let fragment = Html::parse_fragment(&resp);
    let num_data = get_num_data(&fragment);
    for i in num_data {
        let input_id = format!("#sample-input-{}", i);
        let output_id = format!("#sample-input-{}", i);

        // TODO: Better error handling
        let input_data = match parse_data(&input_id, &fragment) {
            Ok(data) => data,
            Err(error) => panic!(error),
        };
        let output_data = match parse_data(&output_id, &fragment) {
            Ok(data) => data,
            Err(error) => panic!(error),
        };

        let test_case = TestCase::new(input_data, output_data);
        test_cases.push(test_case);
    }
    cache_cases(problem_num, &test_cases);
    Ok(test_cases)
}

fn get_num_data(fragment: &Html) -> usize {
    let selector = Selector::parse(".sampledata").unwrap();
    fragment.select(&selector).count() / 2
}

fn parse_data(css_id: &str, fragment: &Html) -> Result<String, Box<dyn std::error::Error>> {
    let selector = Selector::parse(css_id).unwrap();
    return match fragment.select(&selector).next() {
        Some(tag) => {
            let data = tag.text().collect::<Vec<_>>();
            Ok(data.join(""))
        }
        None => Err("No such selector found!".into()),
    };
}
