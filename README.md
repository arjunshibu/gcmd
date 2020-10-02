# gcmd
A command-line wrapper written in Go to avoid commonly used one-liner commands.

Inspired by @tomnomnom 's [gf](https://github.com/tomnomnom/gf)

## Why?
When doing recon and dealing with large amounts of data, I often end up with complex commands.
For example I use [httpx](https://github.com/projectdiscovery/httpx) to probe for hosts like this
```
$ echo hackerone.com | httpx -silent -json -response-in-json -o data.json
```
Full output of this command is [here](https://raw.githubusercontent.com/arjunshibu/gcmd/master/data.json)

To get only the HTTP response of the host from this file I use
```
$ cat data.json | jq .serverResponse | sed "s/\\\r\\\n/\\n/g;s/\"//g"
```
It's easy to mess up commands like this especially when having complex *Regex* patterns.

With `gcmd` you can give names to command combinations and reuse them anytime.

Saving this command is easy as
```
$ gcmd -save -i httpx-response "jq .serverResponse | sed 's/\\\r\\\n/\\n/g;s/\"//g'"
```
Now I can use it as
```
$ cat data.json | gcmd httpx-response
```

### Installation
`gcmd` requires Go
```
$ go get -v github.com/arjunshibu/gcmd
```

### Usage
```
$ gcmd -h
```
This will show help for the tool.
Available switches are
| Flag                    | Description                                             | Example                                            |
|-------------------------|---------------------------------------------------------|----------------------------------------------------|
| -ls                     | List available commands                                 | gcmd -ls                                           |
| -save                   | Save a command                                          | gcmd -save test-cmd 'cat /etc/passwd \| grep root' |
| -i                      | Take input from stdin (for -save only)                  | gcmd -save -i test-cmd 'sort -u \| wc -l'          |
| -echo                   | Prints the command rather than executing it             | gcmd -echo test-cmd                                |
| -rm                     | Remove a command                                        | gcmd -rm test-cmd                                  |

### Command Files
You can create config directory as `~/.config/gcmd` or `~/.gcmd`

Command files are stored as JSON files like this
```
$ cat ~/.gcmd/httpx-response.json
{
   "cmds": [
      {
         "name": "jq",
         "args": ".serverResponse"
      },
      {
         "name": "sed",
         "args": "'s/\\\\r\\\\n/\\n/g;s/\"//g'"
      }
   ],
   "stdin": true
}
```

Some example command files are available in `examples` directory. Copy them to your config directory.
```
$ cp -r ~/go/src/github.com/arjunshibu/gcmd/examples ~/.gcmd
```

### Auto Completion
Tab autocompletion scripts for bash and zsh are included. So you can hit Tab to show your commands.
```
$ gcmd <TAB>
httpx-response       test-cmd
```
#### Bash
Place this in `~/.bashrc`
```
source ~/go/src/github.com/arjunshibu/gcmd/autocompletion/gcmd.bash
```
#### Zsh
Enable autocomplete if haven't done (no need if you have oh-my-zsh) by placing in `~/.zshrc`
```
autoload -U compaudit && compinit
```
Place this in `~/.zshrc`
```
source ~/go/src/github.com/arjunshibu/gcmd/autocompletion/gcmd.zsh
```
## Issues

When saving commands, seperate pipes with whitespaces like `find . | wc -l` instead of `find .|wc -l` because if you do so, it will get saved as
```
{
   "cmds": [
      {
         "name": "find",
         "args": ".|cat|wc -l"
      }
   ]
}
```
This may not be a problem, but avoid it.
#
gcmd cannot run commands in background. Eg: `vim &`, because it run commands as `bash -c "command"`.

## Contribution

Pull requests are always welcome. Feel free to contribute your ideas or bug fixes :heart:.
