# 📌 API REST pour une To-Do List en Golang

## 📖 Description
Cette API REST permet aux utilisateurs de gérer une liste de tâches avec des fonctionnalités CRUD (Create, Read, Update, Delete). Elle est construite en Go avec le framework **Gin** et utilise **PostgreSQL** ou **SQLite** comme base de données. L'authentification est gérée avec **JWT**.

## 🚀 Fonctionnalités
- 🔑 Inscription et connexion des utilisateurs (JWT)
- ✅ Ajout, modification, suppression et récupération de tâches
- 📌 Statuts des tâches : `à faire`, `en cours`, `terminé`
- 🛠️ Documentation API avec Swagger
- 🔒 Sécurisation des endpoints
- 📦 Stockage des données avec PostgreSQL ou SQLite

---

## 🛠️ Technologies utilisées
- **Langage :** Go
- **Framework Web :** Gin (`gin-gonic/gin`)
- **ORM :** GORM (`gorm.io/gorm`)
- **Base de données :** PostgreSQL / SQLite
- **Authentification :** JWT (`golang-jwt/jwt`)
- **Gestion de configuration :** Godotenv (`joho/godotenv`)
- **Migration DB :** Golang Migrate (`golang-migrate/migrate`)
- **Documentation :** Swaggo (`swaggo/swag`)

---

## 📁 Structure du projet
```
todo-api/
│── cmd/                     # Point d'entrée principal
│   ├── main.go               # Lancement du serveur
│
├── config/                  # Gestion de la configuration
│   ├── config.go             # Chargement des variables d'env
│
├── internal/                 # Logique métier
│   ├── models/               # Définition des modèles
│   ├── repository/           # Gestion des interactions DB
│   ├── services/             # Logique métier
│   ├── handlers/             # Gestion des routes et controllers
│
│
├── routes/                   # Définition des routes
│
├── tests/                    # Tests unitaires et d'intégration
│
├── .env                      # Variables d’environnement
├── go.mod                    # Dépendances du projet
├── README.md                 # Documentation du projet
```

---

## 📥 Installation
### 1️⃣ Cloner le projet
```sh
git clone https://github.com/user/todo-api.git
cd todo-api
```
### 2️⃣ Installer les dépendances
```sh
go mod tidy
```
### 3️⃣ Configurer les variables d'environnement
Créer un fichier `.env` et y ajouter :
```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=todo_db
JWT_SECRET=my_secret_key
```
### 4️⃣ Lancer les migrations
```sh
go run cmd/migrate.go
```
### 5️⃣ Démarrer le serveur
```sh
go run cmd/main.go
```

---

## 🔥 Endpoints de l'API
### 🔑 Authentification
- **POST** `/register` → Inscription d'un utilisateur
- **POST** `/login` → Connexion et récupération du JWT

### ✅ Gestion des tâches (nécessite un JWT)
- **GET** `/tasks` → Récupérer toutes les tâches
- **POST** `/tasks` → Ajouter une tâche
- **GET** `/tasks/{id}` → Récupérer une tâche spécifique
- **PUT** `/tasks/{id}` → Modifier une tâche
- **DELETE** `/tasks/{id}` → Supprimer une tâche

---

## 🛠️ Documentation API
Générer et afficher la documentation Swagger :
```sh
swag init
```
L'interface Swagger sera accessible à : `http://localhost:8080/swagger/index.html`

---

## 🧪 Tests
Lancer les tests unitaires et d'intégration :
```sh
go test ./...
```

---

## 📝 Licence
Ce projet est sous licence MIT.

---

## 🤝 Contribution
Les contributions sont les bienvenues !
1. Forkez le projet
2. Créez une branche (`feature/ma-feature`)
3. Commitez vos modifications (`git commit -m 'Ajout d'une nouvelle fonctionnalité'`)
4. Poussez votre branche (`git push origin feature/ma-feature`)
5. Ouvrez une Pull Request 🚀


