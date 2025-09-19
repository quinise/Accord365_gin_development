# Accord365 (Gin + Blockchain)

[![Go](https://img.shields.io/badge/go-1.20-blue?logo=go)](https://go.dev/)
[![Solidity](https://img.shields.io/badge/solidity-ERC20-363636?logo=ethereum)](https://docs.soliditylang.org/)
[![Postgres](https://img.shields.io/badge/postgres-15-blue?logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)

Accord365 is an Ethereum-based banking application built using the Go Gin framework and modern web technologies. The project explored blockchain integrations, user account management, and real-world transaction workflows.

Although development was not fully completed, this repo documents significant progress toward building a decentralized banking app that combines blockchain with a traditional backend.

## 🚀 Features & Accomplishments
### - Go/Gin Server
- Built server with HTML templates and dynamic routes.
- Implemented cookies and session management.

### - Authentication & Security
- Integrated OAuth2 for Google in Gin framework.
- Validated JavaScript forms using regular expressions.

### - Database Layer
- Designed and updated Postgres and MySQL schemas.
- Created database models and implemented Create & Update queries.

### - Blockchain & Web3 Integration
- Developed an ERC-20 token buying page using Solidity, JSON, and JavaScript.
- Connected with Web3.js API for contract and transaction calls.
- Implemented wallet functionality for users to load Ether.
- Used JavaScript Promises to complete web3 transactions.

### - Testing & Tools
- Conducted blockchain testing with Ganache, Metamask, Geth, Rinkeby, and Ropsten testnets.
- Used Truffle-Contracts to interact with smart contracts.

### - Front-End Work
- Built pages with JavaScript, jQuery, HTML, CSS, and Bootstrap.
- Prototyped wireframes and collaborated with the client in focus groups.
- Experimented with the Buffalo framework (Go front end).

## 🛠️ Tech Stack
### Frontend
- HTML, CSS, Bootstrap
- JavaScript, jQuery
- Web3.js

### Backend
- Go (Gin framework)
- MySQL, PostgreSQL
- Solidity (ERC-20 smart contracts)

### Tools & Testing
- Truffle, Ganache, Metamask, Geth
- OAuth2 (Google)
- Regular expression validation

## ⚡ Getting Started
This project was developed using **Go Gin**, **Postgres/MySQL**, and **Web3.js**.  
To explore the code:  
1. Clone the repository  
2. Review the `/models`, `/routes`, and `/blockchain` directories  
3. See `main.go` for the application entry point

⚠️ Note: This project is not currently maintained as a runnable app,  
but the codebase provides a solid reference for blockchain integration in Gin.  

## 📌 Project Notes
This repo represents my work during the development phase of Accord365. While I was not able to complete every feature, the following milestones were achieved:

- Successful integration of blockchain interactions into a Go/Gin web app.
- End-to-end testing of token transactions in Ethereum testnets.
- Client meetings, documentation, and planning to support the next developer.

## 🤝 Collaboration
- Partnered with teammates on JavaScript and Go development challenges.
- Reached out for peer support on routing problems.
- Led debriefs with the client to ensure transparency and alignment.

## 📚 Lessons Learned
Working on Accord365 gave me hands-on experience with:
- Designing for both traditional databases and blockchain systems in the same application  
- Implementing authentication flows in a new backend framework (Gin)  
- Balancing client requirements with technical feasibility through prototyping and iteration  
- The importance of clear project organization, which I improved by documenting plans and creating a mobile to-do tracker  

## 🔮 Future Development
If continued, the next developer could expand this foundation by:

- Finalizing the Gin routing logic and improving scalability.
- Enhancing the wallet system for production readiness.
- Migrating front-end components to a modern framework for a smoother UX.
- Strengthening error handling and transaction security in blockchain interactions.
