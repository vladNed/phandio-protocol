use serde::Deserialize;
use serde_json;

use std::{
    fs::File,
    io::Read,
    process::{Command, Stdio},
};

#[derive(Debug, Deserialize)]
struct ContractConfig {
    proxy: String,
    chain: String,
    gas_limit: u64,
    pem_path: String,
    deploy_args: Vec<String>,
    contract_address: String,
}

impl ContractConfig {
    fn load() -> Self {
        let mut file = File::open("interact/config.json")
            .expect("-- Failed to open deploy configuration !!");

        let mut config_content = String::new();
        file.read_to_string(&mut config_content)
            .expect("-- Failed to read deploy configuration !!");

        serde_json::from_str(&config_content).expect("-- Failed to parse deploy configuration !!")
    }

    fn deploy_contract(&self) {
        let mut cmd = Command::new("mxpy");
        cmd.args(&[
            "--verbose",
            "contract",
            "deploy",
            "--recall-nonce",
            "--bytecode=output/swap-router.wasm",
            format!("--pem={}", self.pem_path).as_str(),
            format!("--proxy={}", self.proxy).as_str(),
            format!("--chain={}", self.chain).as_str(),
            format!("--gas-limit={}", self.gas_limit).as_str(),
        ]);
        cmd.arg("--arguments");
        for arg in self.deploy_args.clone() {
            cmd.arg(arg);
        }
        cmd.arg("--send");
        cmd.stdout(Stdio::inherit());
        cmd.stderr(Stdio::inherit());

        let status = cmd.status().expect("-- Failed to execute deploy command !!");
        if status.success() {
            println!("Contract deployed successfully !!");
        } else {
            println!("Failed to deploy contract !!");
        }
    }

    fn _upgrade_contract(&self) {
        let mut cmd = Command::new("mxpy");
        cmd.args(&[
            "--verbose",
            "contract",
            "upgrade",
            self.contract_address.as_str(),
            "--recall-nonce",
            "--bytecode=output/swap.wasm",
            format!("--pem={}", self.pem_path).as_str(),
            format!("--proxy={}", self.proxy).as_str(),
            format!("--chain={}", self.chain).as_str(),
            format!("--gas-limit={}", self.gas_limit).as_str(),
            "--send"
        ]);
        cmd.stdout(Stdio::inherit());
        cmd.stderr(Stdio::inherit());

        let status = cmd.status().expect("-- Failed to execute upgrade command !!");
        if status.success() {
            println!("Contract upgraded successfully !!");
        } else {
            println!("Failed to upgrade contract !!");
        }

    }
}

fn main() {
    let config = ContractConfig::load();
    config.deploy_contract();
}
