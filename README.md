# 🚖 TestDrive
## What does it do?
TestDrive automatically scrapes input/output data from [BOJ(Baekjoon Online Judge)](https://www.acmicpc.net/) and runs tests for your executable binary file!

## Discalimer
I have changed the language from Rust to Go for several reasons (development performance, etc)
So the **How to use** part is not currently working, and I haven't worked on the command line interface. 
I will be working on it ASAP ☺️

## How to use (Not available)
1. Clone this repo
2. Build with [Rust](https://www.rust-lang.org/) with `cargo build --release`
3. `test_drive` gets 2 arguments. First, the path to the binary, second the problem number. Following is the example running `test_drive` for [#1000 (A+B)](https://www.acmicpc.net/problem/1000) at BOJ.
### Example
```shell
/path-to-testdrive/target/release/test_drive ~/BOJ/bin/1000 1000
```
### Output
```
Running tests on 1 cases...
[Case #1 Passed!]
```
### Tip
If your editor that you use for problem solving supports any kind of user scripts, it is recommended to automate the process using it. Currently I am using VSCode as my main editor, and my user script for TestDrive is the following:
```json5
{
    "tasks": [
        {
            "label": "compile and run for C++ (with BOJ support)",
            "command": "g++-11",
            "args": [
                "${file}",
                "-o",
                "${fileDirname}/bin/${fileBasenameNoExtension}",
                "-std=c++11",
                "&&",
                // Path to TestDrive
                "~/Documents/test_drive/target/release/test_drive",
                // Path to executable binary 
                "${fileDirname}/bin/${fileBasenameNoExtension}",
                // Problem number (I save my source code name as the problem number)        
                "${fileBasenameNoExtension}"
            ],
            "group": "build",
            // Problem matcher
            "problemMatcher": {
                "fileLocation": [
                    "relative",
                    "${workspaceRoot}"
                ],
                "pattern": {
                    // The regular expression. 
                    "regexp": "^(.*):(\\d+):(\\d+):\\s+(warning error):\\s+(.*)$",
                    "file": 1,
                    "line": 2,
                    "column": 3,
                    "severity": 4,
                    "message": 5
                }
            }
        },
    ]
}
```

BTW I save my source code name as the problem number, so the process is simplified.

## ⚠️ Warning! ⚠️
This project contains lots of unhandled errors, bugs, etc. Also, it currently only supports macOS, since this program uses `~/Library/Caches/` as the directory to store caches. 