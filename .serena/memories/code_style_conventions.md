# Code Style & Conventions

## Go Backend (CLI & API)

### File Organization
- `cmd/` - CLI commands with Cobra framework
- `crypto/verkle/` - Verkle tree implementation  
- `blockchain_integration/` - Ethereum integration
- `config/` - Configuration management

### Naming Conventions
- **Files**: snake_case (e.g., `term_aggregation.go`, `api_server.go`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Constants**: UPPER_SNAKE_CASE
- **Structs**: PascalCase

### Error Handling
```go
// Standard pattern
if err != nil {
    return fmt.Errorf("failed to do X: %w", err)
}

// CLI commands use cobra.Command with RunE
RunE: func(cmd *cobra.Command, args []string) error {
    // implementation
}
```

## React/TypeScript Frontend

### File Organization
- `src/app/` - Next.js app router pages
- `src/components/` - Reusable React components
- `src/lib/` - Utility libraries (API, blockchain)
- `src/providers/` - React context providers

### Naming Conventions
- **Files**: PascalCase for components, camelCase for utilities
- **Components**: PascalCase function names
- **Hooks**: camelCase starting with `use`
- **Types**: PascalCase interfaces

### Component Structure
```typescript
'use client'; // For client components

import { useState } from 'react';
import { type ComponentProps } from '@/lib/types';

export function ComponentName({ prop1, prop2 }: ComponentProps) {
  const [state, setState] = useState();
  
  return (
    <div className="tailwind-classes">
      {/* JSX */}
    </div>
  );
}
```

### Web3 Integration
- Uses wagmi v2 + viem v2 for Ethereum interaction
- ConnectKit for wallet connection UI
- Type-safe contract interactions
- Error boundary patterns for Web3 failures

## Configuration Management
- Environment variables in `.env` (never commit secrets)
- JSON config in `config/micert.json` 
- TypeScript config strict mode enabled
- ESLint with Next.js recommended rules