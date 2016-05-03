![License](http://img.shields.io/:license-mit-blue.svg)

# gotoc
:page_facing_up: _Markdown table of contents generator_

![logo_godoc](img/gotoc.png)

`gotoc` is a tool made to generate markdown table of contents inside local git repo.  
Links generated refers to Github anchors.

## Install

```sh
go install github.com/Belekkk/gotoc
```


## Usage

```sh
gotoc -file=README.md
```

```md
<!-- Table of Contents generated by [gotoc](https://github.com/Belekkk/gotoc) -->
**Table of Contents** (by [`gotoc`](https://github.com/Belekkk/gotoc))
- [gotoc](#gotoc)
  - [Install](#install)
  - [Usage](#usage)
  - [Custom TOC title](#custom-toc-title)
  - [Max heading level for TOC entries](#max-heading-level-for-toc-entries)
  - [Features](#features)
    - [Done](#done)
    - [Still in development](#still-in-development)
```

## Custom TOC title

To specify custom TOC title like `**Repo : Table of Contents**` you can pass the argument : `-title='<yourtitle>'`.  
To remove title from TOC, just use the option `-notitle`.

## Max heading level for TOC entries

To limit TOC entries to a specified level of headings, use : `-depth=3`.

## Features

### Done

- [X] Generate TOC a top of file
- [X] Enable custom TOC title editing
- [X] Limit TOC entries

### Still in development

- [ ] Update an existing TOC in Markdown file
- [ ] Handle multiple files input
- [ ] Adding TOC to all files in a directory/sub directories
