# mmv

_Move multiple files_

mmv lets you move or rename one or more files at once, by editing their paths in your favourite text editor

## API

```sh
$ mmv [directory]
```

## Example

```sh
$ mmv ./mydir
```

This will open the contents of the `mydir` in the text editor pointed to by 1) the `$VISUAL` environment variable, 2) the `$EDITOR` variable, or `vim` if neither are set.

You can then edit the filenames in your editor. When you save and exit, `mmv` will rename each file. 

`mmv` is useful for complex renaming tasks, such as:

- Using your text editor's powerful search/replace to edit filenames
- Replacing spaces in filenames with underscores or dashes
- Adding or removing a common substring from multiple filenames
