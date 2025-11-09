# 🏢 API de Gestión de Empresas

## Base URL
```
http://localhost:8080/api/companies
```

## Endpoints Disponibles

### 1. Listar Todas las Empresas
```http
GET /api/companies
GET /api/companies?status=active
GET /api/companies?status=inactive
```

**Query Parameters:**
- `status` (opcional): `active` | `inactive` - Filtrar por estado

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "companies": [
    {
      "id": "uuid-123",
      "code": "10",
      "name": "nexti",
      "business_type": "Organización empresarial",
      "whatsapp_number": "+593992686734",
      "phone_number_id": "123456789",
      "is_active": true,
      "created_at": "2025-10-23T10:00:00Z",
      "updated_at": "2025-10-23T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

### 2. Obtener Empresa por ID
```http
GET /api/companies/{id}
```

**Respuesta exitosa (200):**
```json
{
  "id": "uuid-123",
  "code": "10",
  "name": "nexti",
  "business_type": "Organización empresarial",
  "whatsapp_number": "+593992686734",
  "phone_number_id": "123456789",
  "is_active": true,
  "created_at": "2025-10-23T10:00:00Z",
  "updated_at": "2025-10-23T10:00:00Z"
}
```

---

### 3. Crear Nueva Empresa
```http
POST /api/companies
Content-Type: application/json
```

**Body:**
```json
{
  "code": "ACTIVA_TEST",
  "name": "Empresa Activa Test",
  "business_type": "Organización empresarial",
  "whatsapp_number": "+593999999999",
  "phone_number_id": "123456789012345",
  "access_token": "EAAxxxxxxxxxx",
  "webhook_token": "mi_token_seguro_123"
}
```

**Campos:**
- `code` (requerido): Código único de la empresa
- `name` (requerido): Nombre de la empresa
- `business_type` (opcional): Tipo de organización
- `whatsapp_number` (requerido): Número de WhatsApp
- `phone_number_id` (opcional): ID del número en Meta
- `access_token` (opcional): Token de acceso de Meta
- `webhook_token` (opcional): Token de verificación de webhook

> **Nota:** Si se proporcionan `phone_number_id` y `access_token`, la empresa se activa automáticamente.

**Respuesta exitosa (201):**
```json
{
  "success": true,
  "company": {
    "id": "uuid-456",
    "code": "ACTIVA_TEST",
    "name": "Empresa Activa Test",
    "business_type": "Organización empresarial",
    "whatsapp_number": "+593999999999",
    "phone_number_id": "123456789012345",
    "is_active": true,
    "created_at": "2025-10-23T12:00:00Z",
    "updated_at": "2025-10-23T12:00:00Z"
  },
  "message": "Empresa creada exitosamente"
}
```

---

### 4. Actualizar Empresa
```http
PUT /api/companies/{id}
Content-Type: application/json
```

**Body (todos los campos son opcionales):**
```json
{
  "name": "Nuevo Nombre",
  "business_type": "Otro tipo",
  "whatsapp_number": "+593988888888",
  "phone_number_id": "987654321",
  "access_token": "EAAyyyyyyyyyy",
  "webhook_token": "nuevo_token"
}
```

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "company": { /* empresa actualizada */ },
  "message": "Empresa actualizada exitosamente"
}
```

---

### 5. Activar Empresa
```http
POST /api/companies/{id}/activate
```

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "company": {
    "id": "uuid-123",
    "is_active": true,
    "updated_at": "2025-10-23T12:05:00Z"
  },
  "message": "Empresa activada exitosamente"
}
```

---

### 6. Desactivar Empresa
```http
POST /api/companies/{id}/deactivate
```

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "company": {
    "id": "uuid-123",
    "is_active": false,
    "invalid_date": "2025-10-23T12:10:00Z",
    "updated_at": "2025-10-23T12:10:00Z"
  },
  "message": "Empresa desactivada exitosamente"
}
```

---

### 7. Eliminar Empresa
```http
DELETE /api/companies/{id}
```

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Empresa eliminada exitosamente"
}
```

---

## Códigos de Estado HTTP

| Código | Descripción |
|--------|-------------|
| `200` | Operación exitosa |
| `201` | Recurso creado |
| `400` | Solicitud inválida (JSON mal formado) |
| `404` | Empresa no encontrada |
| `409` | Conflicto (empresa ya existe) |
| `500` | Error interno del servidor |

---

## Ejemplos de Uso con Frontend

### JavaScript/TypeScript (Fetch API)

```typescript
// Listar empresas activas
const listActiveCompanies = async () => {
  const response = await fetch('http://localhost:8080/api/companies?status=active');
  const data = await response.json();
  return data.companies;
};

// Crear nueva empresa
const createCompany = async (companyData) => {
  const response = await fetch('http://localhost:8080/api/companies', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(companyData),
  });
  return await response.json();
};

// Activar empresa
const activateCompany = async (companyId) => {
  const response = await fetch(
    `http://localhost:8080/api/companies/${companyId}/activate`,
    { method: 'POST' }
  );
  return await response.json();
};

// Desactivar empresa
const deactivateCompany = async (companyId) => {
  const response = await fetch(
    `http://localhost:8080/api/companies/${companyId}/deactivate`,
    { method: 'POST' }
  );
  return await response.json();
};

// Actualizar empresa
const updateCompany = async (companyId, updates) => {
  const response = await fetch(
    `http://localhost:8080/api/companies/${companyId}`,
    {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updates),
    }
  );
  return await response.json();
};

// Eliminar empresa
const deleteCompany = async (companyId) => {
  const response = await fetch(
    `http://localhost:8080/api/companies/${companyId}`,
    { method: 'DELETE' }
  );
  return await response.json();
};
```

### React Hook Example

```typescript
// hooks/useCompanies.ts
import { useState, useEffect } from 'react';

interface Company {
  id: string;
  code: string;
  name: string;
  business_type: string;
  whatsapp_number: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export const useCompanies = (status?: 'active' | 'inactive') => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCompanies();
  }, [status]);

  const fetchCompanies = async () => {
    try {
      setLoading(true);
      const url = status 
        ? `http://localhost:8080/api/companies?status=${status}`
        : 'http://localhost:8080/api/companies';
      
      const response = await fetch(url);
      const data = await response.json();
      
      setCompanies(data.companies || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const activateCompany = async (id: string) => {
    const response = await fetch(
      `http://localhost:8080/api/companies/${id}/activate`,
      { method: 'POST' }
    );
    if (response.ok) {
      await fetchCompanies(); // Recargar lista
    }
  };

  const deactivateCompany = async (id: string) => {
    const response = await fetch(
      `http://localhost:8080/api/companies/${id}/deactivate`,
      { method: 'POST' }
    );
    if (response.ok) {
      await fetchCompanies(); // Recargar lista
    }
  };

  return {
    companies,
    loading,
    error,
    refresh: fetchCompanies,
    activateCompany,
    deactivateCompany,
  };
};
```

---

## CORS

El backend ya tiene configurado CORS para permitir peticiones desde cualquier origen (`*`). En producción, deberías configurarlo para dominios específicos.

---

## Notas Importantes

1. **Seguridad**: Los campos `access_token` y `webhook_token` **nunca** se retornan en las respuestas JSON (están marcados como `json:"-"`).

2. **Auto-activación**: Al crear una empresa con credenciales de Meta completas, se activa automáticamente.

3. **Filtros**: Usa `?status=active` o `?status=inactive` para filtrar empresas por estado.

4. **UUIDs**: Los IDs se generan automáticamente como UUIDs v4.

5. **Timestamps**: Todas las fechas están en formato ISO 8601 UTC.

---

## Testing con cURL

```bash
# Listar empresas activas
curl http://localhost:8080/api/companies?status=active

# Crear empresa
curl -X POST http://localhost:8080/api/companies \
  -H "Content-Type: application/json" \
  -d '{
    "code": "TEST_001",
    "name": "Mi Empresa",
    "business_type": "Retail",
    "whatsapp_number": "+1234567890"
  }'

# Activar empresa
curl -X POST http://localhost:8080/api/companies/{id}/activate

# Desactivar empresa
curl -X POST http://localhost:8080/api/companies/{id}/deactivate
```

---

**¡Tu API de empresas está lista!** 🎉

