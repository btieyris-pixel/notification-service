# Notification Service v1

Servicio backend encargado de enviar notificaciones push a drivers.

## Arquitectura

Order Service → SNS → SQS → Notification Service → FCM → Mobile

## Funcionalidad

- Consume eventos desde SQS
- Obtiene token desde PostgreSQL
- Envía push vía FCM
- Retry simple (3 intentos)
- Elimina mensaje solo si éxito

## Variables

- POSTGRES_DSN
- AWS_REGION
- SQS_QUEUE_URL
- FCM_ENABLED
- FCM_PROJECT_ID
- FCM_CREDENTIALS_JSON

## Eventos

- NEW_ORDER
- ORDER_STATUS_UPDATED