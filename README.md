# qo

A tool for creating customizable sandboxed Linux environments for educational testing and evaluation.

`qo` enables instructors to create secure and reproducible testing environments where students can complete coding challenges in isolated Linux sandboxes. The tool provides complete control over available commands and binaries while automatically generating detailed evaluation reports.

## Features

- **Secure Sandboxing**: Creates isolated Linux environments using namespaces for safe student testing
- **Time-Locked Challenges**: Encrypts challenge archives with unlock times to prevent early access
- **Customizable Environments**: Control exactly which binaries and commands are available to students
- **Reproducible**: Ensures consistent testing environments across different machines

## Prerequisites

- Linux operating system (required for sandboxing features).
- [Go](https://go.dev/doc/install) installed on your system.


## Installation

Run this command to install `qo`.

```bash
curl -fsSL  https://raw.githubusercontent.com/ahmedYasserM/qo/main/scripts/install.sh | bash
```

## Quick Start

### For Instructors

1. Prepare your [challenge folder](#challenge-folder-structure) with levels and check scripts
2. Build and encrypt the challenge archive:
   ```bash
   qo build -f <challenge folder> -p <password> -k <starterkey> -u <unlock date and time>
   ```
   For example:
   ```bash
   qo build -f ./my-challenges -p mypassword -k starterkey -u "2025-12-01 14:30"
   ```

### For Students

1. Start the test session with the encrypted archive:
   ```bash
   sudo qo start -i <student id> -a <challenge archive> -p <password> -k <starter key> -d <duration>
   ```
   For example:
   ```bash
   sudo qo start -i 2021170034 -a test.enc -p mypassword -k starterkey -d 90m
   ```
   _**Note:** Setting duration is not implemented yet. The option is required but has no effect._


## Usage

### Instructor Command: `build`

Prepares and encrypts challenge folders for secure distribution to students.

**Workflow:**
1. Validates [challenge folder](#challenge-folder-structure) structure and scripts
2. Compresses folder into archive format
3. Encrypts with time-lock and starter key
4. Outputs ready-to-distribute encrypted file

**Required Flags:**
- `-f, --folder` — Path to challenge folder
- `-p, --password` — Archive encryption password  
- `-k, --key` — Starter key for students
- `-u, --unlock-time` — Unlock time (`YYYY-MM-DD HH:MM` format)

**Optional Flags:**
- `-o, --output` — Output path (default: `eval-archive.enc`)

**Example:**
```bash
qo build -f ./challenges -p securepass -k abc123 -u "2025-07-10 09:30" -o midterm-exam.enc
```

### Student Command: `start`

Launches secure testing environment for students to complete challenges.

**Workflow:**
1. Prompts for Student ID (used in reports and logs)
2. Verifies starter key and enforces unlock time
3. Creates isolated sandbox environment
4. Extracts challenges and starts interactive shell
5. Monitors all commands and activities
6. Generates evaluation report upon completion

**Required Flags:**
- `-i, --id` — Student ID
- `-a, --archive` — Path to encrypted challenge archive
- `-p, --password` — Archive decryption password
- `-k, --key` — Starter key provided by instructor
- `-d, --duration` — Test duration (e.g., `90m`, `2h`, `1h30m`) _(required but not implemented yet)_

**Optional Flags:**
- `-o, --output` — Results directory (default: `eval-results`) _(not implemented yet)_

**Example:**
```bash
sudo qo start -i 2021170034 -a midterm-exam.enc -p securepass -k abc123 -d 2h 
```

## Challenge Folder Structure

Your challenge folder should follow this structure:

```
challenges/
├── level1/
│   ├── description.md
│   ├── check.sh
│   └── files/
├── level2/
│   ├── description.md
│   ├── check.sh
│   └── files/
└── README.md
```

Each level should contain:
- **description.md**: Challenge instructions for students
- **check.sh**: Automated validation script
- **files/**: Any supporting files needed

## Customizing the Sandbox

### Adding System Binaries
First, extract the `rootfs`.

```bash
# in qo/pkg/sandbox/
sudo tar -xzvf rootfs.tar.gz
```
Then, check if the binary you would like to add is available in `busybox`.

```bash
# in qo/pkg/sandbox/bin/

./busybox --list | grep [command]
```
You will encounter one of two cases:
#### 1. The binary is available in busybox
Create a symbolic link to busybox with the name of the binary.
```bash
ln -s busybox [command]
```
#### 2. The binary is not available in busybox

In this case, you can copy the binary (along with its dependencies) from your system to the environment. To do so, use the provided script:

```bash
# In qo/scripts/
./inject.sh /usr/bin/gcc /path/to/your/rootfs
./inject.sh /bin/nano /path/to/your/rootfs
```

Im both cases, make sure to recompress the `rootfs` and recompile after modification.

```bash
# in qo/pkg/sandbox/
sudo tar -czvf rootfs.tar.gz rootfs
cd ../..
# In qo/
go install
```

### Environment Configuration

You can customize the sandbox environment by modifying:
- Available commands and utilities
- File system permissions
- Available users and groups, etc

## Coming soon

- **Command Monitoring**: Log all student commands and activities during testing sessions
- **Automated Reporting**: Generate comprehensive PDF reports of student performance 
- **Set Challenge Duaration**: Automatically end the session after a specified duration
