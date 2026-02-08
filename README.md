# 🐝 NearBeee — *buzzzz* the map

[![Kotlin Multiplatform](https://img.shields.io/badge/Kotlin_Multiplatform-FFD700?style=flat&logo=kotlin&logoColor=000000)](https://kotlinlang.org/docs/multiplatform.html)
[![Compose Multiplatform](https://img.shields.io/badge/Compose_UI-000000?style=flat&logo=jetpackcompose&logoColor=FFD700)](https://www.jetbrains.com/lp/compose-multiplatform/)
[![Go](https://img.shields.io/badge/Go_Backend-FFD700?style=flat&logo=go&logoColor=000000)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-000000?style=flat&logo=postgresql&logoColor=FFD700)](https://www.postgresql.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-FFD700?style=flat&logo=mongodb&logoColor=000000)](https://www.mongodb.com/)
[![Docker](https://img.shields.io/badge/Docker-000000?style=flat&logo=docker&logoColor=FFD700)](https://www.docker.com/)

**NearBeee** 
is a geo-social ecosystem that bridges the gap between digital talk and physical locations.
Discover any spots, join urban discussions, and see what’s *buzzzzing* at any point on the map.

---

## ✨ Key Features

* **Spatial Discovery:** Your feed dynamically adapts to your coordinates—see only what's happening around you.
* **Location-Tied Discussions:** Launch threads at any spot, from urban landmarks to local meetups, and engage with those nearby.
* **Cross-Platform:** A seamless, unified experience across Android and iOS devices.
* **Live Interaction:** Real-time commentary and instant feedback from the people sharing your space.

## 👥 Authors

* **Client:** [@skrpld](https://github.com/skrpld)
* **Server:** [@kr0uch](https://github.com/kr0uch)
* **UI/UX:** [@zaydiF](https://github.com/zaydiF)

## 🧰 Tech Stack

The project is built on a robust foundation, divided into client and server modules:

### Client
* [**Kotlin Multiplatform (KMP)**](https://kotlinlang.org/docs/multiplatform.html) — shared business logic for Android and iOS.
* [**Compose Multiplatform**](https://www.jetbrains.com/lp/compose-multiplatform/) — declarative shared UI.
* [**Ktor**](https://ktor.io/) — lightweight asynchronous multiplatform client.
* [**KOIN**](https://insert-koin.io/) — pragmatic dependency injection.

### Server
* [**Go**](https://go.dev/) — high performance and reliability.
* [**PostgreSQL**](https://www.postgresql.org/) — primary relational database for structured data.
* [**MongoDB**](https://www.mongodb.com/) — flexible NoSQL storage for posts and activity logs.
* [**Docker**](https://www.docker.com/) — containerization for rapid deployment and scaling.

## 🚀 Installation

### [Download the latest stable version from the releases page](https://github.com/skrpld/NearBeee/releases).

### Or for Developers:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/skrpld/NearBeee.git
    ```
2.  **Start Backend (Docker):**
    ```bash
    cd NearBeee/server # navigate to server directory
    docker-compose up -d
    ```
3.  **Run Mobile App:**
    Open the project in **Android Studio** and select the desired target (Android or iOS).

---
Built to create a *buzzzz* wherever you land. 🐝