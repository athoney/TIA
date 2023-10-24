# TIA

*Project Lead:* Alicia Thoney (@athoney)

*TOC:*
- [Default Project Structure](#default-project-structure)
  - [About](#about)
  - [Folder Structure](#folder-structure)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [Technologies](#technologies)

## About

The Threat Information Assistant (TIA) is a web application that compiles and displays vulnerabilities associated with a specific software or hardware configuration. The user will provide a configuration and TIA will represent the cyber-threat information (CTI) using the Structured Threat Information eXpression (STIX) language.

<!-- [See WIKI](https://uwcedar.io/lab-tech/default-project/wikis/home)
 change this URL to your project's wiki after reviewing the default project wiki -->

## Folder Structure

(In progress)

- **.gitignore** => This file excludes files/folders from being committed to your repository. Use this file to exclude large files that exceed the gitlab upload limit of 10MB, build artifacts, etc.
- **code/** => This folder contains all the code being used for the project.
- **papers/** => This directory contains all the papers relevant to the current research project. These are the papers that make up the literature review.
  - **papers.md** => This file contains a summary/overview of all the papers for quick reference.
  - **TIA_Project_Proposal.pdf** => Official project proposal documentation.
- **results/** => This folder contains all the raw and processed results related to the experiments found in the code/ folder.
- **README.md** => This is the first file that is read by a newcomer to your project. This file needs to provide enough details and link to all the major topics in the repository.

## Installation

- Install Go
- Install Postgres (see `./code/db/Tia Database Instructions.txt` for more info)

## Usage

CD to `./code/server/server.go` and run `go run server.go`. Project runs locally on port 8080. Landing page: http://localhost:8080/

## Contributing

#### Graduate Mentors:
- Clay Carper
- Danny Radosevich
- Andey Robins

#### Project Lead:
- Alicia Thoney

#### Developers:
- Ally Hays
  - @ahays8
- Calvin VanWormer
  - @cvanworm
- Jenna Goodrich
  - @jgoodr10
- Marc Wodahl 
  - @mwodahl1
- William Frost
  - @wfrost2

## Technologies
- Bootstrap 4 (v4.6)
- Go (v1.18.3)
- Node.js (v16.15.1)
- Postgres (v14)
