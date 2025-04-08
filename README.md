# Go Script to Count and Find Following List

This Go script is designed to count and find the following list for a list of people from a raw user data string.So you can find common streamers or check the percentage of any channel's current viewers to overlapping audiences, or even find bot accounts.

## Overview

The script reads a raw list of usernames and processes it to count the occurrences same streamer of each user, identifying the following relationships among them. This is useful for applications that analyze user connections or interactions.

## Features

- Parse a raw list of usernames.
- Count how many times each user appears in the list of following.
- Determine the percentage of tracked users and overlapping audiences.

## Installation

1. Ensure you have [Go](https://golang.org/dl/) installed on your system.

2. Clone this repository:

   ```bash
   git clone https://github.com/fl1ghterr/TwitchCommonFollowersFinder.git
   cd TwitchCommonFollowersFinder
  '''
3. Run the project
  ```bash
   go run main.go
  ```

4. All done
