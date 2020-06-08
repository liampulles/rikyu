<div align="center"><img src="bamboo-vase.webp" alt="Photograph of a bamboo vase"></div>
<div align="center"><small><sup>Yuji, Sawaki. Bamboo Vase. 1940. Photograph. Sekisui Museum, Tu-city, Japan </sup></small></div>
<h1 align="center">
  Rikyu
</h1>

<h4 align="center">A DVD-ripping tool for power users.</a></h4>

<p align="center">
  <a href="#status">Status</a> •
  <a href="#key-objectives">Key Objectives</a> •
  <a href="#system-requirements">System Requirements</a> •
  <a href="#install">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#how-it-works">How it works</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

<p align="center">
  <a href="https://travis-ci.com/liampulles/rikyu">
    <img src="https://travis-ci.com/liampulles/rikyu.svg?branch=master" alt="[Build Status]">
  </a>
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/liampulles/rikyu">
  <a href="https://goreportcard.com/report/github.com/liampulles/rikyu">
    <img src="https://goreportcard.com/badge/github.com/liampulles/rikyu" alt="[Go Report Card]">
  </a>
  <a href="https://codecov.io/gh/liampulles/rikyu">
    <img src="https://codecov.io/gh/liampulles/rikyu/branch/master/graph/badge.svg" />
  </a>
</p>

## Status

Rikyu is currently in heavy development, and as such does not have a release candidate available.

## Key Objectives

* Use Docker for cross-platform jobs.
* Actual command line arguments to x264, mkvmerge, etc. are left entirely up to the user. This provides high flexibility for advanced use cases (e.g. PAL slowdown).
* Specify "human steps" in the pipeline, for e.g. verification stages, subtitle tweaking, author comments, etc.
* Specify pipelines as DAGs, to enable parallel processing of steps.
* Provide a "project" concept to houses multiple templates, which can select titlesets in a DVD structure based on previous selections.
* Use mustache templates on pipelines for reusability.
* Utilize IMDB metadata

## System Requirements

* Docker

Yep, that is all you need!

## Install

  1. Clone this repository for the bleeding edge version, or download from the [Releases](https://github.com/liampulles/cabiria/releases) page for a stable version.
  1. Run `make install`

## Usage

  1. Start `rikyud`
  1. Head to http://localhost:9119 and start running and observing jobs!

## Workflow

### Ripping the DVDs

1. Start the interface
1. Select a Project / Create new
1. Insert the DVD
1. Interface is populated with title sets
1. Largest title is automatically selected as 'Main' if previous categories don't exist on project - otherwise those are assigned where predicted.
1. Assign categories to titles (can create new category)
    * Categories can be unique or sequential, if sequential the no. is guessed, but can be set explicitly.
    * You can guess assign categories.
    * You can launch the title in a video player, which starts from the middle, just to help you see what it is.
1. Select 'Rip' to rip the titles to appropriate category folders.

## Contributing

Please submit an issue with your proposal.

## License

See [LICENSE](LICENSE)
