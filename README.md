# mmv

_Move multiple files_

mmv lets you move or rename one or more files at once, by editing their paths in your favourite text editor

## API

```sh
$ mmv [directory]
```

## Example

```sh
$ ls ./mydir
a.txt
b.txt
d.txt
$ mmv ./mydir
```

This will open the contents of the `mydir` in the text editor pointed to by 1) the `$VISUAL` environment variable, 2) the `$EDITOR` variable, or `vim` if neither are set.

```
a.txt
b.txt
d.txt
```

Let's rename `d.txt` to `c.txt`, and save and close the file:

```
a.txt
b.txt
c.txt
```

`mmv` has renamed `d.txt` to `c.txt`:

```sh
$ ls ./mydir
a.txt
b.txt
c.txt
```

This is a trivial example. `mmv` becomes more useful with more complex renaming tasks, such as:

- Using your text editor's powerful search/replace to edit filenames
- Replacing spaces in filenames with underscores or dashes
- Adding or removing a common substring from multiple filenames
