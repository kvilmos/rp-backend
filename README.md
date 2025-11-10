# Room Planner Backend

Ez a repository a "Room Planner" alkalmazás Go nyelven írt backendjét tartalmazza. A szolgáltatás felelős az üzleti logika kezeléséért, az adatbázis-műveletekért, valamint a fájlok tárolásáért.

## Architektúra és Technológiák

- **Nyelv:** Go
- **Adatbázis:** MariaDB (MySQL-kompatibilis) a főbb adatok tárolására. (A beállításához lásd az `rp-database` repositoryt).
- **Objektumtároló:** MinIO a felhasználó által feltöltött fájlok (pl. bútormodellek, képek) tárolására.
- **Gyorsítótár/In-memory tároló:** Redis a gyorsítótárazott adatok kezelésére.
- **Konténerizáció:** Docker és Docker Compose a fejlesztői környezet függőségeinek (MinIO, Redis) egyszerűsített kezelésére.

## Előfeltételek

- [Go](https://go.dev/doc/install) (1.20+ verzió javasolt)
- [Docker](https://www.docker.com/products/docker-desktop/) és Docker Compose
- [A projekt adatbázis repositoryja (`rp-database`) letöltve és beállítva.](https://github.com/kvilmos/rp-database)

## Telepítés és Futtatás

A teljes alkalmazás elindításához kövesd az alábbi 3 fő lépést.

### 1. Adatbázis elindítása

Ez a szolgáltatás az `rp-database` repositoryban definiált adatbázisra támaszkodik. Mielőtt elindítanád a backendet, győződj meg róla, hogy az adatbázis fut.
[rp-database](https://github.com/kvilmos/rp-database)

### 2. Függőségek indítása (MinIO és Redis)

A backendnek szüksége van egy MinIO és egy Redis szerverre is. Ezeket a projektben található Docker Compose fájlokkal indíthatod el.

1.  **MinIO indítása:** Nyiss egy terminált, navigálj a `dev-minio` mappába, majd futtasd a következő parancsot:

    ```bash
    cd dev-minio
    docker-compose up -d
    ```

    Ez elindít egy MinIO szervert, és egy segéd szolgáltatást, ami automatikusan létrehozza a szükséges `furniture-models` bucket-et, és beállítja a jogosultságokat és a webhookokat.

    - **MinIO konzol:** `http://localhost:9001`
    - **Felhasználó:** `minioadmin`
    - **Jelszó:** `minioadmin`

2.  **Redis indítása:** Nyiss egy másik terminált, navigálj a `dev-redis` mappába, és futtasd a következő parancsot:
    ```bash
    cd dev-redis
    docker-compose up -d
    ```
    Ez elindítja a Redis szervert, ami a `localhost:6379` címen lesz elérhető.

### 3. Backend alkalmazás indítása

Ha mindhárom külső szolgáltatás (MariaDB, MinIO, Redis) fut, elindíthatod magát a Go alkalmazást.

1.  Navigálj az `rp-backend` projekt gyökérkönyvtárába.
2.  Töltsd le a Go függőségeket:
    ```bash
    go mod tidy
    ```
3.  Indítsd el az alkalmazást:
    ```bash
    go run .
    ```

A backend sikeresen elindult és csatlakozott a többi szolgáltatáshoz. A `storage` csomagban lévő kódrészletek alapján az alkalmazás a következő alapértelmezett címekkel próbál csatlakozni:

- **MySQL:** `develop:develop@tcp(127.0.0.1:3306)/room-planner`
- **MinIO:** `127.0.0.1:9000`
- **Redis:** `localhost:6379`

## A Fejlesztői Környezet Leállítása

Ha befejezted a munkát, a következő parancsokkal állíthatod le a háttérben futó konténereket:

```bash
# A MinIO leállítása
cd dev-minio
docker-compose down

# A Redis leállítása
cd dev-redis
docker-compose down

# Az adatbázis leállítása (az rp-database mappában)
cd ../../rp-database
docker-compose down
```
