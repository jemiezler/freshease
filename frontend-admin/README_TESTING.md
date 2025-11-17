# Frontend-Admin Testing Guide

## Overview
This document provides a comprehensive guide for testing the frontend-admin application.

## Test Structure

```
__tests__/
├── components/          # Component tests
│   ├── login.test.tsx
│   └── products.test.tsx
├── integration/         # Integration tests
│   └── auth-flow.test.tsx
└── utils/              # Test utilities
    └── test-helpers.tsx
```

## Setup

### Installation
Tests use Jest and React Testing Library. Dependencies are already installed in `package.json`.

### Configuration
- `jest.config.js` - Jest configuration
- `jest.setup.js` - Test setup and mocks

## Running Tests

### Run all tests
```bash
npm test
```

### Run tests in watch mode
```bash
npm run test:watch
```

### Run tests with coverage
```bash
npm run test:coverage
```

### Run specific test file
```bash
npm test -- login.test.tsx
```

## Test Categories

### 1. Component Tests
Component tests verify that UI components render correctly and handle user interactions.

**Example: Login Page**
- Renders login form
- Displays Google OAuth button
- Shows email/password form
- Handles form submission
- Displays error messages

### 2. Integration Tests
Integration tests verify that multiple components work together correctly.

**Example: Authentication Flow**
- Google OAuth login
- Email/password login
- Error handling
- Navigation flow

## Writing Tests

### Component Test Example
```typescript
import { render, screen } from '@testing-library/react'
import LoginPage from '@/app/login/page'

describe('LoginPage', () => {
  it('renders login page with title', () => {
    render(<LoginPage />)
    expect(screen.getByText('Freshease Admin')).toBeInTheDocument()
  })
})
```

### Integration Test Example
```typescript
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

describe('Authentication Flow', () => {
  it('completes email/password login flow', async () => {
    const user = userEvent.setup()
    render(<LoginPage />)
    
    // Interact with the component
    await user.click(screen.getByText('Continue with Email'))
    
    // Assert expected behavior
    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument()
    })
  })
})
```

## Mocking

### Next.js Router
The router is automatically mocked in `jest.setup.js`.

### API Calls
Mock API calls using Jest mocks:
```typescript
jest.mock('@/lib/auth', () => ({
  startOAuth: jest.fn(),
  loginWithPassword: jest.fn(),
}))
```

### Context Providers
Use the custom render function from `test-helpers.tsx` to include providers:
```typescript
import { render } from '@/__tests__/utils/test-helpers'

render(<Component />) // Includes AuthProvider
```

## Best Practices

1. **Test User Behavior**: Test what users see and do, not implementation details
2. **Use Accessibility Queries**: Prefer `getByLabelText`, `getByRole`, etc.
3. **Wait for Async Operations**: Use `waitFor` for async updates
4. **Mock External Dependencies**: Mock API calls, router, etc.
5. **Clean Up**: Clear mocks between tests

## Coverage Goals

- **Components**: 80%+ coverage
- **Utilities**: 90%+ coverage
- **Integration**: Key user flows covered

## Troubleshooting

### Common Issues

1. **Router not mocked**: Ensure `jest.setup.js` is loaded
2. **Async operations**: Use `waitFor` for async updates
3. **Window methods**: Mock `window.confirm`, `window.alert`, etc.
4. **Next.js components**: Mock Next.js specific components if needed

### Debugging

- Use `screen.debug()` to see rendered output
- Use `jest.fn().mockImplementation(() => console.log(...))` to debug mocks
- Check test output for error messages

## Future Enhancements

1. **E2E Tests**: Add Playwright or Cypress tests
2. **Visual Regression**: Add snapshot tests
3. **Performance Tests**: Add performance testing
4. **Accessibility Tests**: Add accessibility testing

