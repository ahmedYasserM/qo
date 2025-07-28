#!/bin/bash

# Set base folder name
CHALLENGE_DIR="test"

# Create base folder
mkdir -p "$CHALLENGE_DIR"

echo "Creating challenge folder: $CHALLENGE_DIR"

# Level 1
mkdir -p "$CHALLENGE_DIR/level1"
# Add question.txt
cat <<'EOF' > "$CHALLENGE_DIR/level1/question.txt"
Level 1 Challenge:
------------------
Create a directory called "testdir" in your current working directory.

Hint: Use the `mkdir` command.
EOF

# Add check.sh
cat <<'EOF' > "$CHALLENGE_DIR/level1/check.sh"
#!/bin/bash
# Level 1 check: Did the student create a folder called "testdir"?

if [ -d "$PWD/testdir" ]; then
    echo "Level 1 passed!"
    exit 0
else
    echo "Level 1 failed: 'testdir' does not exist."
    exit 1
fi
EOF
chmod +x "$CHALLENGE_DIR/level1/check.sh"

# Level 2
mkdir -p "$CHALLENGE_DIR/level2"
# Add question.txt
cat <<'EOF' > "$CHALLENGE_DIR/level2/question.txt"
Level 2 Challenge:
------------------
Create a new user named "studentuser" on the system.

Hint: Use the `useradd` command.
EOF

# Add check.sh
cat <<'EOF' > "$CHALLENGE_DIR/level2/check.sh"
#!/bin/bash
# Level 2 check: Did the student create a user named "studentuser"?

if id "studentuser" &>/dev/null; then
    echo "Level 2 passed!"
    exit 0
else
    echo "Level 2 failed: user 'studentuser' does not exist."
    exit 1
fi
EOF
chmod +x "$CHALLENGE_DIR/level2/check.sh"

# Level 3
mkdir -p "$CHALLENGE_DIR/level3"
# Add question.txt
cat <<'EOF' > "$CHALLENGE_DIR/level3/question.txt"
Level 3 Challenge:
------------------
Copy the file "secret.txt" from this folder to your home directory and
set its permissions so that only you (the owner) can read and write it.

Hint: Use `cp` and `chmod 600`.
EOF

# Add a supporting file for the challenge
echo "This is a secret file for Level 3." > "$CHALLENGE_DIR/level3/secret.txt"

# Add check.sh
cat <<'EOF' > "$CHALLENGE_DIR/level3/check.sh"
#!/bin/bash
# Level 3 check: Was secret.txt copied and permissions set?

if [ -f "$HOME/secret.txt" ] && [ "$(stat -c "%a" "$HOME/secret.txt")" == "600" ]; then
    echo "Level 3 passed!"
    exit 0
else
    echo "Level 3 failed: Check file copy and permissions."
    exit 1
fi
EOF
chmod +x "$CHALLENGE_DIR/level3/check.sh"

echo "Challenge folder created successfully!"
