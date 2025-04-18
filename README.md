Replacements for some coreutils on Linux, rewrite of [tutils](https://github.com/kivattt/tutils)

<img src="screenshot.png" alt="hex, ls and xxd utilities">

# Commands
`ls` List files\
`pwd` Print working directory\
`cat` Print files or STDIN\
`hex` Encode/decode hexadecimal\
`xxd` Visualize as hex\
`ascii` Print ASCII-range bytes\
`countchars` Show distribution of characters sorted\
`bytes` Show bytes in human-readable size

# Building
`./compile.bash`\
The built binaries will be located in the `bin` directory

# Installing
```console
cd
git clone https://github.com/kivattt/tutils2
cd tutils2
./compile.bash
```
Then add this to your `.bashrc` file, and re-open a terminal
```bash
tutils2path=~/tutils2/bin
if test -d $tutils2path; then
    alias ls="$tutils2path/ls"
    alias pwd="$tutils2path/pwd"
    alias cat="$tutils2path/cat"
    alias hex="$tutils2path/hex"
    alias xxd="$tutils2path/xxd"
    alias ascii="$tutils2path/ascii"
    alias countchars="$tutils2path/countchars"
    alias bytes="$tutils2path/bytes"
else
    echo "Could not find tutils2 programs in $tutils2path"
fi
```

Since adding `tutils2` to your path environment variable could break existing scripts that rely on system utilities specific behaviour, we use shell aliases so that shell scripts will continue to use the existing utilities, rather than `tutils2`.
