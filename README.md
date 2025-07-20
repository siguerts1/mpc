# MPC — Multipass Cluster Controller

**MPC** is the central CLI tool that federates multiple Multipass-capable hosts into a unified VM cluster. It connects to `mpcd` agents running on each host and lets you view and manage distributed Multipass instances from a single command-line interface.

---

## ✅ Features

- Discover Multipass VMs across many hosts
- View instance details per host
- YAML-based host inventory (`hosts.yml`)
- Foundation for launching, executing, and destroying VMs remotely

---

## 🚀 Requirements

- Ubuntu or macOS machine for development
- Golang 1.21 or newer
- Network access to remote hosts where `mpcd` is running
- Multipass installed on the remote hosts (not needed on controller)

---

## 📦 Installation

```bash
git clone https://github.com/yourusername/mpc.git
cd mpc
go mod tidy
go run ./cmd/mpc status --hosts hosts.yml
go run ./cmd/mpc status --instances hosts.yml
