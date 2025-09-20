# 🔗 IU-MiCert API Integration Guide

## 📖 Overview

Connect the IU-MiCert frontend with the issuer backend API for credential verification using modern React libraries.

## 📦 Install Dependencies

```bash
cd packages/client/iumicert-client
npm install @tanstack/react-query axios react-hot-toast react-error-boundary react-dropzone react-hook-form @hookform/resolvers yup
```

## ⚙️ Environment Setup

Create `.env.local`:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
NEXT_PUBLIC_API_TIMEOUT=30000
```

## 🚀 Quick Setup

### 1. Add Providers to Layout

```typescript
// src/app/layout.tsx
'use client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { ErrorBoundary } from 'react-error-boundary';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 5 * 60 * 1000, retry: 3 },
    mutations: { retry: 1 },
  },
});

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <ErrorBoundary fallback={<div>Something went wrong</div>}>
          <QueryClientProvider client={queryClient}>
            {children}
            <Toaster position="top-right" />
          </QueryClientProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
```

### 2. Create API Service

```typescript
// src/services/api.ts
import axios from 'axios';

const client = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
  timeout: 30000,
});

export const apiService = {
  async verifyReceipt(receiptData: any) {
    const response = await client.post('/verifier/receipt', receiptData);
    return response.data;
  },
  
  async getSystemStatus() {
    const response = await client.get('/status');
    return response.data;
  },
};
```

### 3. Create React Query Hook

```typescript
// src/hooks/useVerification.ts
import { useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { apiService } from '../services/api';

export const useVerifyReceipt = () => {
  return useMutation({
    mutationFn: apiService.verifyReceipt,
    onSuccess: () => toast.success('Credential verified!'),
    onError: () => toast.error('Verification failed'),
  });
};
```

### 4. Update Verifier Component

```typescript
// src/app/verifier/page.tsx
import { useVerifyReceipt } from '../hooks/useVerification';

export default function VerifierDashboard() {
  const verifyMutation = useVerifyReceipt();
  
  const handleVerify = async (credentialData: string) => {
    try {
      const parsedData = JSON.parse(credentialData);
      const result = await verifyMutation.mutateAsync(parsedData);
      
      // Handle success - show results
      setVerificationResult(result);
    } catch (error) {
      // Handle error - already shown via toast
    }
  };

  return (
    <VerificationUpload
      isLoading={verifyMutation.isPending}
      onVerify={handleVerify}
    />
  );
}
```

## 📡 Available API Endpoints

### ✅ Currently Working
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/status` | System status |
| `GET` | `/api/health` | Health check |
| `GET` | `/api/terms` | List terms |
| `POST` | `/api/receipts` | Generate receipt |
| `GET` | `/api/students` | List students |
| `POST` | `/api/blockchain/publish` | Publish to blockchain |

### ❌ Missing (Need Implementation)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/verifier/receipt` | **Verify complete receipt** |
| `GET` | `/api/verifier/journey/{student_id}` | **Get student journey** |

## 🧪 Testing

```bash
# Test backend
curl -X POST http://localhost:8080/api/verifier/receipt \
  -H "Content-Type: application/json" \
  -d @test-receipt.json

# Start backend
cd packages/issuer && ./dev.sh

# Start frontend
cd packages/client/iumicert-client && npm run dev
```

## 🔒 Error Handling

- Network errors: Automatic retry with exponential backoff
- Validation errors: Toast notifications with specific messages
- File upload: Size limits (10MB) and type validation (.json)

---

**Next Steps**: Replace mock verification logic in your components with the real API calls above.