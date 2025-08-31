# fitbyte
Go app for second project of ProjectSprint

# How to Run
1. Clone the repo

   ```bash
    git clone git@github.com:fatjan/fitbyte.git
    cd fitbyte
   ```

2. Create `.env` file

  Can copy from `.env-example` but adjust the value
   ```bash
    cp .env-example .env
   ```

3. Create database `fitbyte`

4. Run the migration

  ```bash
    make migrate-up
  ```

5. Run the app

  ```bash
  make run
  ```

6. Health check

```bash
curl http://localhost:3002/health-check
```

