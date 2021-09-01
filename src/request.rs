extern crate reqwest;
extern crate scraper;

use scraper::Html;
use scraper::Selector;

use crate::lib::TestCase;
use crate::request::cache::{cache_cases, check_cache};

pub mod cache;

#[tokio::main]
pub async fn parse_cases(problem_num: &str) -> Result<Vec<TestCase>, Box<dyn std::error::Error>> {
    match check_cache(problem_num) {
        Ok(test_cases) => {
            if !test_cases.is_empty() {
                return Ok(test_cases);
            }
        }
        Err(err) => return Err(err),
    }

    let input = "#sample-input-";
    let output = "#sample-output-";
    let mut test_cases: Vec<TestCase> = Vec::new();
    let url = format!("{}{}", "https://www.acmicpc.net/problem/", problem_num);
    let resp = reqwest::get(url).await?.text().await?;

    let fragment = Html::parse_fragment(&resp);
    let mut i = 1;
    loop {
        let input_id = format!("{}{}", input, i);
        let output_id = format!("{}{}", output, i);
        let input_data = parse_data(&input_id, &fragment);
        let output_data = parse_data(&output_id, &fragment);

        if input_data.is_ok() && output_data.is_ok() {
            test_cases.push(TestCase::new(input_data.unwrap(), output_data.unwrap()));
        } else {
            break;
        }
        i += 1;
    }
    cache_cases(problem_num, &test_cases);

    Ok(test_cases)
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
