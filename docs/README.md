# NowYouSeeMe Documentation

Complete documentation for the NowYouSeeMe platform.

## 📚 Documentation Structure

```
docs/
├── README.md              # This file
├── .context/              # Detailed context for development
│   ├── ARCHITECTURE.md    # System architecture
│   ├── API.md            # API documentation
│   ├── SETUP.md          # Detailed setup guide
│   └── PROJECT_CONTEXT.md # Project context
```

## 🚀 Quick Links

### For Users
- [COMMANDS.md](../COMMANDS.md) - **Quick command reference** (START HERE)
- [QUICKSTART.md](../QUICKSTART.md) - Quick start guide
- [README.md](../README.md) - Project overview

### For Developers
- [Seed Scripts Guide](../sdk/scripts/README.md) - Database seeding and data generation

### For Deep Dive
- [Architecture](../docs/.context/ARCHITECTURE.md) - Technical architecture
- [API Documentation](../docs/.context/API.md) - API endpoints
- [Setup Guide](../docs/.context/SETUP.md) - Detailed setup
- [Project Context](../docs/.context/PROJECT_CONTEXT.md) - Design decisions

## 🎯 Quick Navigation

**I want to...**

- **Get started quickly** → [COMMANDS.md](../COMMANDS.md)
- **Run the platform** → `make backend` + `make frontend`
- **Test CRUD operations** → `make test-full`
- **Add sample data** → `make populate`
- **Seed database** → [Seed Scripts Guide](../sdk/scripts/README.md)
- **Understand the system** → [Architecture](../docs/.context/ARCHITECTURE.md)

## 📖 Documentation by Role

### 🔰 Beginner
Start here if you're new to the project:
1. [README.md](../README.md) - What is this?
2. [COMMANDS.md](../COMMANDS.md) - How do I use it?
3. [QUICKSTART.md](../QUICKSTART.md) - Step-by-step setup

### 👨‍💻 Developer
You want to understand and modify the code:
1. [Architecture](../docs/.context/ARCHITECTURE.md) - System design
2. [API Documentation](../docs/.context/API.md) - API details
3. [Project Context](../docs/.context/PROJECT_CONTEXT.md) - Why decisions were made

### 🤖 SDK & Database Seeding
You want to populate the database or integrate with the platform:
1. [Seed Scripts Guide](../sdk/scripts/README.md) - Database seeding with Event Sourcing
2. [API Documentation](../docs/.context/API.md) - Current Event Sourcing API reference

**Note:** Python SDK is being updated to match the new Event Sourcing API.
See API documentation for current endpoint specifications.

## 🎨 Visual Structure

```
NowYouSeeMe
├── 📖 User Documentation
│   ├── README.md              ⭐ Start here
│   ├── COMMANDS.md            ⭐ Most useful
│   └── QUICKSTART.md
│
├── 🔧 Developer Documentation
│   └── docs/.context/
│       ├── ARCHITECTURE.md
│       ├── API.md
│       ├── SETUP.md
│       └── PROJECT_CONTEXT.md
│
└── 🤖 Scripts & Utilities
    └── sdk/scripts/
        └── README.md   ⭐ Seed scripts (Event Sourcing)
```

## 💡 Tips

- **Just want to use it?** → Read [COMMANDS.md](../COMMANDS.md)
- **Want to understand it?** → Read [ARCHITECTURE.md](../docs/.context/ARCHITECTURE.md)
- **Want to extend it?** → Read [PROJECT_CONTEXT.md](../docs/.context/PROJECT_CONTEXT.md)
- **Seeding database?** → Read [Seed Scripts Guide](../sdk/scripts/README.md)

## 🔍 Finding What You Need

| I need... | Go to... |
|-----------|----------|
| Quick commands | [COMMANDS.md](../COMMANDS.md) |
| Setup steps | [QUICKSTART.md](../QUICKSTART.md) |
| API endpoints | [API.md](../docs/.context/API.md) |
| System design | [ARCHITECTURE.md](../docs/.context/ARCHITECTURE.md) |
| Seed scripts | [Seed Scripts Guide](../sdk/scripts/README.md) |
| Design decisions | [Project Context](../docs/.context/PROJECT_CONTEXT.md) |

---

**Stuck?** Start with [COMMANDS.md](../COMMANDS.md) - it has everything you need!
