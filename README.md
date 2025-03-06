# 🚀 Custom Version Control System

This is a lightweight version control system built in Go, inspired by Git. It provides essential version control commands to track changes, manage branches, and commit updates efficiently.

## 🛠 Features

- Add files to the staging area
- Commit changes
- View commit logs
- Check the status of your working directory
- Compare file differences
- Manage branches
- Switch between branches

## 📌 Commands

| Command   | Description |
|-----------|-------------|
| `add`    | Add files to the staging area |
| `branch` | Create and list branches |
| `commit` | Commit staged changes |
| `diff`   | Show differences between working directory, index, and commits |
| `init`   | Initialize a new repository |
| `log`    | View commit history |
| `status` | Show the working directory and staging area status |
| `switch` | Switch between branches, with `-c` flag to create a branch if it does not exist |

## 🚀 Getting Started

1. Clone the repository:
   ```sh
   git clone https://github.com/aryandutt/gvc.git
   cd gvc
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Build the project:
   ```sh
   go build -o gvc cmd/main.go
   ```
4. Initialize a repository:
   ```sh
   ./gvc init
   ```
5. Start tracking files:
   ```sh
   ./gvc add file.txt
   ./gvc commit -m "Initial commit"
   ```

## 🤝 Contributing

This project is still evolving, and contributions are welcome! If you'd like to help make the project more organized and provide a great learning opportunity for others, feel free to open issues or submit pull requests.

## 📜 License

This project is open-source and available under the MIT License.

---
Happy coding! 😊

