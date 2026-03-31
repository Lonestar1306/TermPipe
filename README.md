# ⚡ TermPipe Builder

<p align="center">
Interactive terminal tool to **visually build Unix pipelines** and generate the equivalent Bash command in real time.
</p>

<p align="center">
<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go"/>
<img src="https://img.shields.io/badge/TUI-BubbleTea-ff69b4?style=for-the-badge"/>
<img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge"/>
<img src="https://img.shields.io/badge/Platform-Linux%20%7C%20MacOS-lightgrey?style=for-the-badge"/>
</p>

---

# 🚀 What is TermPipe?

**TermPipe Builder** is an interactive **Terminal UI (TUI)** that helps you experiment with classic Unix text-processing tools like:

`grep | awk | sed | jq | sort | uniq | wc | head | tail`

Instead of writing complex shell pipelines manually, **you build them step-by-step**, preview the results instantly, and get the **final Bash command automatically generated**.

Perfect for:

* exploring logs
* inspecting CSV files
* extracting fields
* testing regex
* learning Unix pipelines
* debugging shell commands

---

# ✨ Why TermPipe?

Writing pipelines like this:

```bash
cat logs.txt | grep "ERROR" | awk -F';' '{print $3}' | sort | uniq -c
```

can be annoying when experimenting.

With **TermPipe** you:

1️⃣ Add filters interactively
2️⃣ See the output immediately
3️⃣ Copy the generated command when ready

---

# 🎬 Example

Input data:

```
save-ts;event-ts;op;D1;D2
2026-01-26 03:10:59;2026-01-26 02:10:31;13;45.89;13.49
2026-01-26 03:10:59;2026-01-26 02:09:08;12;3.54;9.25
```

Interactive pipeline:

```
[1] grep '26'
[2] awk -F';' col 3
[3] sort
[4] uniq -c
```

Generated Bash command:

```bash
cat input.txt | grep -E -i '26' | awk -F';' '{print $3}' | sort | uniq -c
```

All while seeing **live preview of the result**.

---

# 🧠 Supported Pipeline Tools

| Tool   | Description             |
| ------ | ----------------------- |
| `grep` | Regex filtering         |
| `awk`  | Column extraction       |
| `sed`  | Regex replacement       |
| `jq`   | JSON field extraction   |
| `sort` | Sort lines              |
| `uniq` | Remove duplicates       |
| `wc`   | Count lines/words/chars |
| `head` | First N lines           |
| `tail` | Last N lines            |
| `>`    | Redirect output to file |

---

# 🖥 Interface Highlights

### Live pipeline visualization

```
[1] grep 'error'   [2] awk col 3   [3] sort
```

### Output preview pane

Scrollable with:

* mouse wheel
* arrow keys

### Automatic Bash generation

```
$ cat input.txt | grep -E -i 'error' | awk '{print $3}' | sort
```

---

# 📦 Installation

### Requirements

* **Go 1.21+**

### Clone the repo

```bash
git clone https://github.com/YOUR_USERNAME/termpipe.git
cd termpipe
```

### Build

```bash
go build -o termpipe
```

### Run

```bash
./termpipe
```

---

# ⚡ Usage

### Interactive demo mode

```bash
./termpipe
```

Loads built-in sample data.

---

### Use with pipelines

```bash
cat logfile.txt | ./termpipe
```

or

```bash
journalctl | ./termpipe
```

---

# ⌨️ Keyboard Shortcuts

| Key      | Action         |
| -------- | -------------- |
| `g`      | grep filter    |
| `a`      | awk column     |
| `s`      | sed replace    |
| `j`      | jq filter      |
| `o`      | sort           |
| `u`      | uniq           |
| `w`      | wc             |
| `h`      | head           |
| `t`      | tail           |
| `x`      | export to file |
| `Delete` | undo last step |
| `Esc`    | cancel input   |
| `Enter`  | confirm step   |
| `q`      | quit           |

---

# 👀 Smart AWK Helper

When you use `awk`, TermPipe automatically:

✔ detects common delimiters
✔ shows a **column preview**
✔ helps you select the column index

Example:

```
1: timestamp | 2: event | 3: value
```

---

# 🛠 Built With

* **Go**
* **Bubble Tea** (TUI framework)
* **Bubbles** (UI components)
* **Lip Gloss** (terminal styling)

---

# 💡 Use Cases

TermPipe is great for:

* Log analysis
* CSV inspection
* Learning Unix pipelines
* Testing regex filters
* Exploring JSON logs
* Debugging shell scripts

---

# 🧭 Roadmap

Future ideas:

* copy command to clipboard
* pipeline presets
* interactive regex tester
* CSV mode
* export pipeline as script
* syntax highlighting
* multi-column awk support

---

# 🤝 Contributing

Contributions are welcome!

If you have ideas, improvements, or bug fixes:

1. Fork the repository
2. Create a branch
3. Submit a PR 🚀

---

# 📄 License

MIT License

---

# ⭐ If you like this project

Give it a **star on GitHub** — it really helps!

And feel free to share it with anyone who loves **terminal tools**.
