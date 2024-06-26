# Markdown to HTML

A simple markdown to HTML converter written in Go. It kind of works. It only supports part of the markdown.

## Usage

1. Build the app

```
git clone https://github.com/Lisviks/markdown-to-html.git
cd markdown-to-html
go build
```

2. Convert markdown to HTML

```
./md-to-html sample.md          # On Unix systems
md-to-html.exe sample.md        # On Windows
```

Replace `sample.md` with the path to your markdown file. Output file will be named the same as the input. In this case it will be named `sample.html`. To output a file with a different name add a second argument. HTML files will be created inside `out` directory.

```
./md-to-html sample.md index          # On Unix systems
md-to-html.exe sample.md index        # On Windows
```

To output to a different directory use add a third argument with a name of a directory. Default is `out`

```
./md-to-html sample.md index output          # On Unix systems
md-to-html.exe sample.md index output        # On Windows
```
