# Summary — Session 2026-06-22

## Estado del proyecto

- Rama actual: `feat/agents-branching-strategy` (basada en `dev`)
- Commit más reciente: `5506cf5 docs: add branching strategy to AGENTS.md`
- PR #1 abierta a `dev` (sin mergear)

## Decisión de arquitectura

- **AWS** como proveedor cloud
- **Lambda + API Gateway HTTP** (API serverless)
- **DynamoDB** como base de datos (pendiente de migrar, sin tablas aún)
- **Terraform** como IaC, dentro del repo en `infra/`
- **JWT auth** en la misma Lambda (middleware Gin)
- **Despliegue**: GitHub Action manual + local via Docker
- **Dominio**: API Gateway endpoint directamente (sin CloudFront)

## Infraestructura desplegada (dev)

- **Bucket S3**: `terraform-infra-211125667058` (estado remoto)
- **Lambda**: `TRAMA-LAMBDA-DEV` (runtime `provided.al2023`, placeholder hello-world)
  - **Problema resuelto**: error GLIBC_2.34 (cambiado de provided.al2 a provided.al2023)
- **API Gateway**: `TRAMA-GATEWAY-DEV`, ruta `ANY /api/{proxy+}`, CORS `*`
- **IAM Role Lambda**: `TRAMA-LAMBDA-DEV` (logs + dynamodb cuando se añada)
- **IAM Role GitHub**: `TRAMA-GITHUB-DEV` (AdministratorAccess, OIDC)
- **OIDC Provider**: GitHub Actions

**Endpoint**: https://o432bm06q8.execute-api.eu-west-1.amazonaws.com

**Lambda falla** aún con Internal Server Error — está pendiente de depurar (probablemente el binario no arranca bien en provided.al2023).

## Terraform

- Estado remoto en S3: `s3://terraform-infra-211125667058/trama/{env}/terraform.tfstate`
- Backend init: `make tf-init-dev` (migrar estado local → S3 pendiente)
- Despliegue via Docker (no necesita terraform local)

## Makefile targets nuevos

- `build-lambda` — compila binario Go para Lambda
- `tf-init-dev`, `tf-init-prod` — init backend S3
- `tf-plan-dev`, `tf-plan-prod` — plan via Docker
- `tf-apply-dev`, `tf-apply-prod` — apply via Docker

## Pendiente para mañana

1. **Depurar Lambda** — el placeholder responde Internal Server Error. Revisar logs de CloudWatch.
2. **Configurar retention** CloudWatch logs a 30 días (ya en terraform, sin deploy)
3. **Sistema de audit** con S3 + Athena (sin Glue crawler para ahorrar)
4. **Migrar estado local a S3** — hacer `make tf-init-dev` y responder yes
5. **Mergear PR #1** a `dev`
6. **Fase 2**: migrar SQLite → DynamoDB, Gin → Lambda wrapper, JWT auth
