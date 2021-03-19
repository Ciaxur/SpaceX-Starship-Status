## SpaceX SN Status
Go CLI App that checks **SpaceX's Starship** Status through given Twitter User.
Given a Twitter User ID, the app parses through their latest tweets to find tweets related to the SpaceX Startship launch/status.

---

### 💻 Supported Operating Systems
`Linux`
- All Features

`MacOS/Win10`:
- All CLI Features
- Notification Feature not Supported (*notify-send*)
  - Can probably forward the Command Required through the Script for it to work 😊

### ✨ Scripts
- `run_check.sh`: Helper Script that sets the right args to the app for the `notify-send` command
  - A script for **dynamically** issuing a Notification from the Active Session, check out [Dbus-Session-Notify](https://gist.github.com/Ciaxur/53f88d82721141461bc4f8e556f40860)
  - Be sure to set the script variables as needed. Full Path is required

### 📦 Dependancies
- `go`: Used to Build the App
- `notify-send` (*Optional*): Be sure to have this package on your system. As this will be used to issue a notification
  - Can be substituted based on your OS

### 🔧 Setup
1) Copy the [.env.sample](.env.sample) file to **.env** and setup variables
2) Build using Golang (*instructions below*)

### 🚀 Build & Run
```sh
# Create a Build Directory
mkdir build

# Use Golang to Build Application
go build -o build/app ./src
```

### 📄 License
Licensed under the [MIT License](LICENSE).