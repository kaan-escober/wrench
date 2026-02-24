# Getting Started

## Prerequisites

- [Factory CLI](https://docs.factory.ai/cli/) installed
- Go 1.21 or higher (for building from source)

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/kartvya69/byok-cli.git
cd byok-cli/droid-cfg
```

### 2. Build the binary

```bash
go build -o droid-cfg .
```

### 3. (Optional) Install system-wide

```bash
# Linux / macOS
mv droid-cfg /usr/local/bin/droid-cfg

# Termux (Android)
mv droid-cfg $PREFIX/bin/droid-cfg
```

### 4. Run

```bash
droid-cfg
```

## First Run

On first launch you will land on the main menu. Each row shows a live summary of the current value for that category pulled from `~/.factory/settings.json`.

```
D R O I D  CONFIG
  ~/.factory/settings.json

>  BYOK   Custom Models        no models configured
   MOD    Model & Reasoning    opus  ·  auto
   AUTO   Autonomy             normal
   DISP   Display              github  ·  pinned
   SND    Sound                fx-ok01  ·  always
   SEC    Security             ● shield  ● co-author  ○ bg-proc
   BEHV   Agent Behavior       ● cloud  ○ hooks  ● droids
   CMD    Command Policies     factory defaults
```

Navigate with `↑↓` and press `Enter` to open any category.

## Basic Workflow

```
┌──────────────────────┐
│     Main Menu        │
│  Select a category   │
└──────────┬───────────┘
           │ Enter
           ▼
┌──────────────────────┐
│   Category Screen    │
│  Browse settings     │
└──────────┬───────────┘
           │ Enter
           ▼
┌──────────────────────┐
│  Picker / Input      │
│  Choose or type      │
└──────────┬───────────┘
           │ Enter
           ▼
┌──────────────────────┐
│  Auto-saved  ✓       │
│  Returns to list     │
└──────────────────────┘
```

Settings are written to `~/.factory/settings.json` immediately after you confirm a change.

## Next Steps

- [Add a custom model via the BYOK wizard](./byok.md)
- [Browse every setting](./settings.md)
- [See supported providers](./providers.md)
- [Understand config file formats](./configuration.md)
