#!/bin/bash
# Startup script para Azure App Service

# Script de inicio para ejecutar la aplicación
echo "🚀 Iniciando WhatsApp API Server..."
echo "📊 Variables de entorno:"
echo "   - MONGODB_URL: ${MONGODB_URL:0:30}..."
echo "   - API_PORT: $API_PORT"
echo "   - WABA_PHONE_ID: $WABA_PHONE_ID"

# Ejecutar la aplicación
./whatsapp-api-server

