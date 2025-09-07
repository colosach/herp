# Hotel ERP Monorepo

This repository contains the **Hotel ERP** system, structured as a monorepo with two main components:

- **backend/** → Go (Gin) project
- **frontend/** → Vue.js project

We use **Git submodules** to manage the `frontend` directory as a separate Git repository, while the main project is tracked in this parent repo.

---

## 📦 Repository Structure

hotel-erp/
├── backend/ # Go backend application
├── frontend/ # Vue.js frontend (Git submodule)
├── Makefile # Build automation for backend + frontend
├── .gitignore # Root gitignore for monorepo
└── README.md


---

## 🔧 Cloning the Repository

Since the frontend is a **submodule**, you need to initialize and pull its contents:

```bash
# Clone with submodules
git clone --recurse-submodules <repo-url>

# OR if already cloned without submodules
git submodule update --init --remote --recursive

## One-liner (backend + frontend at once)
```bash
git pull origin main && git submodule update --init --remote --recursive

## Auto-update submodules on pull
```bash
git config --global submodule.recurse true
git pull --rebase

