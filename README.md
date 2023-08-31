# `gnols`

The Gno Language Server (`gnols`) is an implementation of the 
[Language Server Protocol][1] (LSP) for the [Gno programming][2] language.

It provides an interface to locally-installed Gno tools (`gno`, `gnokey`, ...)
with the goal of making it easy to add IDE-like features to any application 
that supports the Language Server Protocol.

> **Note** 
>
> This project is community-maintained (not official) and in early
> development.

## Installation

TODO.

## Features

<table>
    <tr>
        <th>Hover Tooltips</th>
        <th>Autocomplete</th>
    </tr>
    <tr>
        <td width="50%">
            <a href="https://github.com/errata-ai/vale/assets/8785025/e3ff4e27-6ba1-456c-b451-856b3e1c9f41">
                <img src="https://github.com/errata-ai/vale/assets/8785025/e3ff4e27-6ba1-456c-b451-856b3e1c9f41" width="100%">
            </a>
        </td>
        <td width="50%">
            <a href="https://github.com/errata-ai/vale/assets/8785025/cf3b36e7-4cf8-4a02-9578-f532a0cb9af6">
                <img src="https://github.com/errata-ai/vale/assets/8785025/cf3b36e7-4cf8-4a02-9578-f532a0cb9af6" width="100%">
            </a>
        </td>
    </tr>
    <tr>
        <td width="50%">
          In-editor documentation for all exported symbols in the Gno standard library.
        </td>
        <td width="50%">Autocomplete for all exported symbols in the Gno standard library.
    </tr>
    <tr>
        <th>Diagnostics</th>
        <th>Formatting</th>
    </tr>
    <tr>
        <td width="50%">
            <a href="https://github.com/errata-ai/vale/assets/8785025/09cb33ea-1bbf-4f89-aec0-e12066622a42">
                <img src="https://github.com/errata-ai/vale/assets/8785025/09cb33ea-1bbf-4f89-aec0-e12066622a42" width="100%">
            </a>
        </td>
        <td width="50%">
            <a href="https://github.com/errata-ai/vale/assets/8785025/2426c64a-1dfe-47f2-a1e0-b7512bca0df4">
                <img src="https://github.com/errata-ai/vale/assets/8785025/2426c64a-1dfe-47f2-a1e0-b7512bca0df4" width="100%">
            </a>
        </td>
    </tr>
    <tr>
        <td width="50%">
            Real-time precompile & build errors in your Gno files.
        </td>
        <td width="50%">
            Format Gno files with the tool of your choice.
        </td>
    </tr>
</table>

## Clients

TODO: tutorials for each client coming soon!

- [ ] Sublime Text
- [ ] Neovim
- [ ] JetBrains IDEs

[1]: https://microsoft.github.io/language-server-protocol/
[2]: https://gno.land/
