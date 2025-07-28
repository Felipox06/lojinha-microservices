@"
# Lojinha Microservices

Sistema de e-commerce usando arquitetura de microserviços em Go.

## Estrutura

- services/
  - user-service/
  - product-service/
  - auth-service/
  - notification-service/
  - api-gateway/
  - config-service/
- shared/
- scripts/
"@ | Out-File -FilePath README.md -Encoding UTF8