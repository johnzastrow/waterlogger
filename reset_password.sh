#!/bin/bash

# Simple password reset script
echo "Resetting password for user 'jcz'..."
echo "New password will be: Password123!"
echo

# Use printf to provide both password entries
printf "Password123!\nPassword123!\n" | ./waterlogger -reset-password jcz

echo "Password reset complete!"
echo "You can now login with:"
echo "Username: jcz"
echo "Password: Password123!"