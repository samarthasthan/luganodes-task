# Ethereum Deposit Tracker

## Overview

The Ethereum Deposit Tracker is designed to monitor and record ETH deposits on the Beacon Deposit Contract. This project efficiently tracks incoming deposits, handles real-time data fetching, and visualization features.

**Note:** Due to time constraints, several aspects of the project are still in need of improvement. These include variable names, additional comments, cache implementation, telegram bot alert and Zipkin tracing.

**Development Note:** To demonstrate functionality, a sample deposit has been added manually as the Beacon Deposit Contract updates infrequently. The tracker will automatically add new transactions to the database as they occur. To check the most recent transactions on the Beacon Deposit Contract, visit [Etherscan](https://etherscan.io/address/0x00000000219ab540356cBB839Cbe05303d7705Fa).

## Technology Stack

- **Backend**: Go
- **Frontend**: Next.js
- **Containerization**: Docker
- **Monitoring & Logging**: Grafana, Loki, Zipkin
- **Caching**: Redis
- **Message Broker**: Kafka

## High-Level Design

![High-Level Design](./others/Luganodes%20Architecture%20Task.png)
_High-level design of the multi-vendor e-commerce platform._

## Links

**Note:** NGINX Reverse proxy can be use to attach the domain

- **Grafana Dashboard**: [http://3.7.73.40:15000/dashboards](http://3.7.73.40:15000/dashboards)
- **Frontend**: [http://3.7.73.40:3000/](http://3.7.73.40:3000/)
- **API Endpoint**: [http://3.7.73.40:8000/deposits?page=1&limit=10](http://3.7.73.40:8000/deposits?page=1&limit=10)

## Snapshots

![Grafana-dashboard](./others/Screenshot%202024-09-10%20at%206.04.05 PM.png)
_Grafana Dashboard with Loki and Prometheus._

![Frontend](./others/Screenshot%202024-09-10%20at%206.04.08 PM.png)
_Next.JS Frontend._

![API-Endpoint](./others/Screenshot%202024-09-10%20at%206.04.10 PM.png)
_API Endpoint._

## Configuration

Configure the following environment variables in your `.env` file:

**Note:** Do not push .env file in production stage

```
# Kafka
KAFKA_PORT=29092
KAFKA_EXTERNAL_PORT=9092
KAFKA_HOST=kafka

# MySQL
MYSQL_PORT=3306
MYSQL_ROOT_PASSWORD=password
MYSQL_HOST=mysql

# Redis
REDIS_PORT=6379
REDIS_HOST=redis

# Ethereum RPC
WS_ENDPOINT=wss://eth-mainnet.g.alchemy.com/v2/Zu9sVXs0A9UIlF-phnzLNo48ch8QMGm6
ADDRESS=0x00000000219ab540356cBB839Cbe05303d7705Fa

# REST API
REST_API_PORT=8000

# Logging, Metrics and Tracing
GRAFANA_PORT=15000
GRAFANA_LOKI_HOST=loki
GRAFANA_LOKI_PORT=3100
PROMETHEUS_PORT=15002
ZIPKIN_HOST=zipkin
ZIPKIN_PORT=9411

# Frontend
FRONTEND_PORT=3000
```

## Setup and Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/samarthasthan/luganodes-task.git
   cd luganodes-task
   ```

2. **Build Docker Containers**:

   ```bash
   docker-compose -f build/compose/compose.yaml up --build
   ```

3. **Run the Application**:

   - The backend service will be accessible at `http://localhost:8000`.
   - The frontend application will be accessible at `http://localhost:3000`.

4. **Grafana and Monitoring**:

   - Grafana dashboard can be accessed at [http://localhost:15000](http://localhost:15000).
   - Prometheus metrics at [http://localhost:15002](http://localhost:15002).
   - Zipkin tracing at [http://localhost:9411](http://localhost:9411).

## Usage

- **Backend**: The backend service tracks Ethereum deposits and stores them in MySQL. It uses Kafka for asynchronous processing and Redis for caching.
- **Frontend**: The frontend provides a user interface to view deposit data.
- **Monitoring**: Check Grafana for visualizations of system metrics and logs.
