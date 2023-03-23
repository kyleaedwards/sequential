# Sequential

An interactive task queue for single-cored organisms

## Installation

```sh
git clone https://github.com/kyleaedwards/sequential
go install .
```

## Usage

To open in interactive mode, simply run `sequential`.

To skip interactive mode and print the current task directly to `stdout`, use `sequential -c`. This may be useful to show your most pressing todo in your command line prompt, tmux status bar, etc.
